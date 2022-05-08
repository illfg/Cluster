package handler

import (
	"Cluster/network"
	"container/list"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackdanger/collectlinks"
)

type Task struct {
	Content   string
	Status    TaskStauts
	TimeStamp time.Time
}

type TaskStauts int32

const (
	Default TaskStauts = iota
	Distached
	Finished
)

type TaskDistributeHandler struct {
	queue    *list.List
	election *ElectionHandler
	offset   int
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
		queue.PushBack(Task{
			Status:  Default,
			Content: line,
		})
	}
	return &TaskDistributeHandler{
		election: election,
		queue:    queue,
	}, nil
}

func (receiver *TaskDistributeHandler) PutEvent(event network.Event) {
	if !receiver.election.IsMaster || receiver.election.InElection {
		return
	}
	switch event.EType {
	case network.TASK:
		task := receiver.nextTask()
		if task != nil {
			log.Printf("[dispatch task] task sended: %v\n",task)
			event.Content = task.Content
			event.EType = network.ACK
			network.Send(event.From, event, nil)
		}
	case network.RESPONSE:
		receiver.finishTask(event.Content)
	}
}

func (receiver *TaskDistributeHandler) DoProcess() {
	for {
		if !receiver.election.InElection && !receiver.election.IsMaster && receiver.election.CurrentMaster != "" {
			client := &network.SyncClient{}
			resp := client.SendSync(receiver.election.CurrentMaster, network.Event{
				EType: network.TASK,
			})
			if resp.EType != network.ACK {
				time.Sleep(time.Second)
				continue
			}
			log.Printf("[handle url] url to fetch: %v\n",resp.Content)
			if err := receiver.handleTask(resp.Content); err != nil {
				network.Send(receiver.election.CurrentMaster, network.Event{
					EType:   network.RESPONSE,
					Content: resp.Content,
				}, nil)
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

func (receiver *TaskDistributeHandler) handleTask(url string) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http get error : %v", err)
	}
	defer resp.Body.Close()
	links := collectlinks.All(resp.Body)
	for _, link := range links {
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			continue
		}
		file, err := os.Create(url)
		if err != nil {
			return err
		}
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		file.WriteString(string(b))
	}
	return nil
}
