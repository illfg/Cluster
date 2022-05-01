package handler

import (
	"Cluster/conf"
	"Cluster/network"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

type ElectionHandler struct {
	InElection     bool
	isMaster       bool
	votedFollower  map[string]string
	follower       map[string]time.Time
	currentMaster  string
	reElectionTime time.Time
}

const electionLogFlag = "election handler"

func (receiver *ElectionHandler) PutEvent(event network.Event) {
	switch event.EType {
	case network.PING:
		receiver.handlePINGEvent(event)
	case network.LEASE:
		receiver.handleLeaseEvent(event)
	case network.VOTE:
		receiver.handleVoteEvent(event)
	case network.HEARTBEAT:
		//receiver.handleHeartbeat(event)
	case network.JOIN:
		receiver.handleJoinEvent(event)
	}
}

func (receiver *ElectionHandler) DoProcess() {
	go func() {
		for {
			//租约过期便发起选举
			if time.Now().After(receiver.reElectionTime) && !receiver.InElection {
				receiver.startElection()
			}
			time.Sleep(time.Second / 5)
		}
	}()
}

func (receiver *ElectionHandler) handlePINGEvent(event network.Event) {
	config := conf.GetConf()
	// glog.Infof("[%s]: receive PING event, from: %s\n", electionLogFlag, event.From)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("MasterCandidate: %s\n", (*config)["MasterCandidate"]))
	builder.WriteString(fmt.Sprintf("CurrentMaster: %s\n", receiver.currentMaster))
	event.Content = builder.String()
	event.EType = network.ACK
	network.Send(event.From, event, nil)
}

func (receiver *ElectionHandler) handleLeaseEvent(event network.Event) {
	glog.Infof("[%s]: receive lease event, form: %s\n", electionLogFlag, event.From)
	//收到租约时回复应答，主节点没收到应答则说明从节点下线，从节点没收到租约则说明主节点下线
	receiver.reElectionTime = time.Now().Add(time.Second * 10)
	receiver.sendACK(event)
	if receiver.currentMaster != event.From {
		receiver.InElection = true
		receiver.joinToMaster(event.From)
	}
}

func (receiver *ElectionHandler) handleVoteEvent(event network.Event) {
	// glog.Infof("[%s]: receive vote event, form: %s\n", electionLogFlag, event.From)
	receiver.votedFollower[event.From] = event.From
	receiver.sendACK(event)
}

//func (receiver *ElectionHandler) handleHeartbeat(event network.Event) {
//	glog.Infof("[%s]: receive heartbeat event, form: %s\n", electionLogFlag, event.From)
//	if receiver.isMaster {
//		receiver.follower[event.From] = time.Now()
//	}
//}

func (receiver *ElectionHandler) handleJoinEvent(event network.Event) {
	glog.Infof("[%s]: receive join event, form: %s\n", electionLogFlag, event)
	if receiver.isMaster {
		receiver.follower[event.From] = time.Now()
		receiver.sendACK(event)
	}
}

func (receiver *ElectionHandler) startElection() {
	glog.Infof("[%s]:start election\n", electionLogFlag)
	receiver.reset()

	//ping all nodes
	masterCandidate, currentMaster := receiver.PINGAllNodes()

	//vote or conn to master
	if currentMaster != "" {
		receiver.joinToMaster(currentMaster)
		return
	}
	masterCandidate = append(masterCandidate, (*conf.GetConf())["IPPort"])
	sort.Strings(masterCandidate)
	receiver.voteToMaster(masterCandidate[0])
}

func (receiver *ElectionHandler) PINGAllNodes() ([]string, string) {
	masterCandidate := make([]string, 0)
	currentMaster := ""

	nodes := strings.Split((*conf.GetConf())["nodes"], ",")
	for _, node := range nodes {
		query := network.Event{EType: network.PING}
		client := network.SyncClient{}
		response := client.SendSync(node, query)
		if response.Content == "" {
			continue
		}
		yaml := conf.ParseYaml(response.Content)
		if yaml["MasterCandidate"] == "true" {
			masterCandidate = append(masterCandidate, response.From)
		}
		if yaml["CurrentMaster"] != "" {
			currentMaster = yaml["CurrentMaster"]
		}
	}
	glog.Infof("[%s]: all nodes had been query, [MasterCandidate]:%v,[currentMaster]:%s\n",
		electionLogFlag,
		masterCandidate,
		currentMaster,
	)

	return masterCandidate, currentMaster
}
func (receiver *ElectionHandler) joinToMaster(IPPort string) {
	glog.Infof("[%s]: join to master, master addr is %s\n", electionLogFlag, IPPort)

	//send join event
	event := network.Event{EType: network.JOIN}
	client := network.SyncClient{}
	response := client.SendSync(IPPort, event)

	if response.EType != network.ACK {
		glog.Warningf("[%s]: join timeout, addr is %s\n", electionLogFlag, IPPort)
		receiver.InElection = false
		return
	}

	//update field
	receiver.isMaster = false
	receiver.InElection = false
	receiver.currentMaster = IPPort
	receiver.reElectionTime = time.Now().Add(time.Second * 10)

	//receiver.reportHeartbeat()
	//通过租约监控健康，由定时器发起下轮选举
}

//func (receiver *ElectionHandler) reportHeartbeat() {
//	receiver.InElection = false
//	for {
//		if !receiver.joined || receiver.isMaster {
//			break
//		}
//		builder := receiver.createCommonContent()
//		joinEvent := network.Event{
//			EType: network.HEARTBEAT,
//			Content: builder.String(),
//		}
//		network.Send(receiver.currentMaster, joinEvent, nil)
//		time.Sleep(time.Second)
//	}
//}

func (receiver *ElectionHandler) voteToMaster(IPPort string) {
	glog.Infof("[%s]: vote to master, master addr is %s\n", electionLogFlag, IPPort)
	if err := receiver.sendVoteEvent(IPPort); err != nil {
		receiver.InElection = false
		return
	}
	//wait for vote
	receiver.waitForVoteResult()
}

func (receiver *ElectionHandler) sendVoteEvent(IPPort string) error {
	config := *conf.GetConf()
	if config["IPPort"] == IPPort {
		receiver.votedFollower[IPPort] = IPPort
		return nil
	}
	voteEvent := network.Event{EType: network.VOTE}
	client := network.SyncClient{}
	response := client.SendSync(IPPort, voteEvent)
	if response.EType != network.ACK {
		glog.Warningf("[%s]: no reply after voted, addr is %s\n", electionLogFlag, IPPort)
		return errors.New("unable to connect master")
	}
	return nil
}

func (receiver *ElectionHandler) waitForVoteResult() {
	timer := time.Now().Add(time.Second * 10)
	for {
		total, _ := strconv.Atoi((*conf.GetConf())["total"])
		if len(receiver.votedFollower) > total/2 {
			receiver.sendLeaseToFollower()
			return
		}
		//如果在等待投票过程中收到租约，则已经加入集群
		if receiver.currentMaster != "" {
			break
		}
		if time.Now().After(timer) {
			glog.Info("wait for vote reuslt timeout")
			break
		}
		time.Sleep(time.Second / 10)
	}
	receiver.InElection = false
}

func (receiver *ElectionHandler) sendLeaseToFollower() {
	glog.Infof("[%s]: was elected master, send lease to follower\n", electionLogFlag)

	receiver.isMaster = true
	receiver.InElection = false
	receiver.currentMaster = (*conf.GetConf())["IPPort"]

	//sendLease
	for _, value := range receiver.votedFollower {
		if value == receiver.currentMaster {
			continue
		}
		event := network.Event{EType: network.LEASE}
		client := network.SyncClient{}
		response := client.SendSync(value, event)
		if response.EType == network.ACK {
			receiver.follower[response.From] = time.Now()
		}
	}
	time.Sleep(time.Second * 3)

	//maintain cluster
	total, _ := strconv.Atoi((*conf.GetConf())["total"])
	for {
		glog.Infof("[%s]: cluster info: %v\n", electionLogFlag, receiver.follower)

		if len(receiver.follower) < total/2 {
			glog.Warningf("[%s]: less than half of total number of followers\n", electionLogFlag)
			receiver.giveUpMaster()
			break
		}
		receiver.sendLease()
		time.Sleep(time.Second)
	}
}

func (receiver *ElectionHandler) reset() {
	receiver.InElection = true
	receiver.isMaster = false

	receiver.votedFollower = make(map[string]string, 8)
	receiver.follower = make(map[string]time.Time, 8)
	receiver.currentMaster = ""
}
func (receiver *ElectionHandler) giveUpMaster() {
	receiver.isMaster = false
	receiver.currentMaster = ""
}

func (receiver ElectionHandler) sendACK(event network.Event) {
	event.EType = network.ACK
	event.Content = ""
	network.Send(event.From, event, nil)
}

func (receiver *ElectionHandler) sendLease() {
	// glog.Infof("[%s]: sending lease to follower\n", electionLogFlag)
	receiver.reElectionTime = time.Now().Add(time.Second * 10)
	if !receiver.isMaster {
		return
	}
	followerToDel := make([]string, 0)
	for key, _ := range receiver.follower {
		event := network.Event{EType: network.LEASE}
		client := network.SyncClient{}
		response := client.SendSync(key, event)
		if response.EType != network.ACK {
			followerToDel = append(followerToDel, key)
		}
		if !receiver.isMaster {
			return
		}
	}
	// glog.Infof("[%s]: node with off-line: %v\n", electionLogFlag, followerToDel)
	for _, key := range followerToDel {
		delete(receiver.follower, key)
	}
}
