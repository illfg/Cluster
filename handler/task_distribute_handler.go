package handler

import (
	"Cluster/conf"
	"Cluster/dao"
	"Cluster/network"
	"container/list"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/jackdanger/collectlinks"
	"github.com/olivere/elastic/v7"
)

var (
	regex, _ = regexp.Compile("[\u4e00-\u9fa5]")
	fliter   = make(map[string]string)
)

type Task struct {
	Content   string
	Status    TaskStauts
	TimeStamp time.Time
}

type ScoreUnit struct {
	url   string
	score float32
}

type TaskStauts int32

const (
	Default TaskStauts = iota
	Distached
	Finished
)

type TaskDistributeHandler struct {
	queue         *list.List
	election      *ElectionHandler
	client        *elastic.Client
	resultUrl     *os.File
	resultContent *os.File
}

func NewTaskDistributeHandler(election *ElectionHandler) (*TaskDistributeHandler, error) {
	queue := list.New()
	b, err := os.ReadFile("task.txt")
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}
		queue.PushBack(&Task{
			Status:  Default,
			Content: line,
		})
	}
	os.Remove("result_url.txt")
	url, err := os.Create("result_url.txt")
	if err != nil {
		return nil, err
	}
	os.Remove("result_content.txt")
	content, err := os.Create("result_content.txt")
	if err != nil {
		return nil, err
	}
	return &TaskDistributeHandler{
		election:      election,
		queue:         queue,
		client:        dao.ConnectES(),
		resultUrl:     url,
		resultContent: content,
	}, nil
}

func (receiver *TaskDistributeHandler) PutEvent(event network.Event) {
	if !receiver.election.IsMaster || receiver.election.InElection {
		return
	}
	switch event.EType {
	case network.TASK:
		// glog.Infof("receiver task req: %v",event)
		task := receiver.nextTask()
		if task != nil {
			log.Printf("[dispatch task] task sended: %v\n", task)
			event.Content = fmt.Sprintf("url: %v\n", task.Content)
			event.EType = network.ACK
			network.Send(event.From, event, nil)
		}
	case network.RESPONSE:
		// glog.Infof("task finished resp: %v",event)
		receiver.finishTask(event.Content)
	}
}

func (receiver *TaskDistributeHandler) DoProcess() {
	for {
		if (!receiver.election.InElection) && (!receiver.election.IsMaster) && (receiver.election.CurrentMaster != "") {
			glog.Info("fetching task")
			time.Sleep(time.Second * 2)
			client := &network.SyncClient{}
			resp := client.SendSync(receiver.election.CurrentMaster, network.Event{
				EType: network.TASK,
			})
			if resp.EType != network.ACK {

				continue
			}
			content := conf.ParseYaml(resp.Content)
			glog.Infof("[handle url] url to fetch: %v\n", content["url"])
			if err := receiver.HandleTask(content["url"]); err == nil {
				client := &network.SyncClient{}
				client.SendSync(receiver.election.CurrentMaster, network.Event{
					EType:   network.RESPONSE,
					Content: content["url"],
				})
			}
		}
	}
}

func (receiver *TaskDistributeHandler) nextTask() *Task {
	head := receiver.queue.Front()
	for head != nil {
		task, ok := head.Value.(*Task)
		if !ok {
			return nil
		}
		switch task.Status {
		case Finished:
			// receiver.queue.Remove(head.Front())
		case Distached:
			if time.Now().After(task.TimeStamp) {
				task.TimeStamp = time.Now().Add(time.Minute)
				return task
			}
		case Default:
			task.Status = Distached
			task.TimeStamp = time.Now().Add(time.Minute)
			return task
		}
		head = head.Next()
	}
	return nil
}

func (receiver *TaskDistributeHandler) finishTask(content string) {
	head := receiver.queue.Front()
	for head != nil {
		task, ok := head.Value.(*Task)
		if !ok {
			return
		}
		if task.Content == content {
			task.Status = Finished
			return
		}
		head = head.Next()
	}
}

func (receiver *TaskDistributeHandler) queueOffset() int {
	offset := 0
	head := receiver.queue.Front()
	for head != nil {
		task, ok := head.Value.(*Task)
		if !ok {
			break
		}
		if task.Status == Finished {
			offset++
			head = head.Next()
			continue
		}
		break
	}
	return offset
}

func (receiver *TaskDistributeHandler) HandleTask(url string) error {
	limit := 10
	queue := list.New()
	queue.PushBack(ScoreUnit{
		url:   url,
		score: 0,
	})
	for i := 0; i < limit && queue.Len() > 0; i++ {
		receiver.sort(queue)
		if err := receiver.bfs(queue); err != nil {
			glog.Errorf("fail to search by bfs: %v", queue)
			return nil
		}
	}
	return nil
}

func (receiver *TaskDistributeHandler) bfs(queue *list.List) error {
	if queue.Len() == 0 {
		return fmt.Errorf("no url remaining")
	}
	unit := queue.Front().Value.(ScoreUnit)
	queue.Remove(queue.Front())
	resp, err := receiver.fetchContent(unit.url)
	if err != nil {
		glog.Errorf("fail to fetch content: %v, skipped", err)
		return nil
	}
	links := collectlinks.All(resp.Body)
	for _, link := range links {
		if !strings.HasPrefix(link, "http") {
			continue
		}
		score := receiver.calcScore(unit.url)
		if score > 10.0 {
			queue.PushBack(ScoreUnit{
				url:   link,
				score: score,
			})
		}
	}
	if _, ok := fliter[unit.url]; ok {
		return nil
	}
	fliter[unit.url] = unit.url
	receiver.writeUrl(unit.url)
	return nil
}

func (receiver *TaskDistributeHandler) calcScore(url string) float32 {
	if receiver.client == nil {
		return 11
	}
	resp, err := receiver.fetchContent(url)
	if err != nil {
		glog.Errorf("fail to fetchcontent: %v", err)
		return 11
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("fail to fetchcontent: %v", err)
		return 11
	}
	words := regex.FindAllString(string(b), -1)
	querys := make([]elastic.Query, 0)
	for _, word := range words {
		querys = append(querys, elastic.NewTermQuery("text", word))
	}
	docs, err := receiver.client.Search("cluster").
		Query(
			elastic.NewBoolQuery().
				Should(querys...)).
		Do(context.Background())
	if err != nil {
		glog.Errorf("fail to query es: %v", err)
		return 11
	}
	type score struct {
		Score float32 `json:"_score"`
	}
	results := docs.Each(reflect.TypeOf(score{}))
	for _, result := range results {
		return result.(score).Score
	}
	return 11
}

func (receiver *TaskDistributeHandler) sort(queue *list.List) {
	head := queue.Front()
	max := float32(0.0)
	maxUnitElem := queue.Front()
	for head != nil {
		unit, ok := head.Value.(ScoreUnit)
		if !ok {
			head = head.Next()
			continue
		}
		if max > unit.score {
			max = unit.score
			maxUnitElem = head
		}
		head = head.Next()
	}
	queue.MoveToFront(maxUnitElem)
}

func (receiver *TaskDistributeHandler) writeUrl(url string) error {
	_, err := receiver.resultUrl.WriteString(url + "\n")
	if err != nil {
		glog.Errorf("fail to fetch content: %v, skipped", err)
		return nil
	}
	resp, err := receiver.fetchContent(url)
	if err != nil {
		glog.Errorf("fail to fetch content: %v, skipped", err)
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("fail to fetch content: %v, skipped", err)
		return nil
	}
	b64 := base64.StdEncoding.EncodeToString(b)
	_, err = receiver.resultContent.WriteString(b64 + "\n")
	if err != nil {
		glog.Errorf("fail to fetch content: %v, skipped", err)
	}
	return nil
}

func (receiver *TaskDistributeHandler) fetchContent(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to make req: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get error : %v", err)
	}
	return resp, nil
}
