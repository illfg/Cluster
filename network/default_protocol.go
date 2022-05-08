package network

import (
	"strconv"
	"strings"
)

type protocol interface {
	make(event Event) string
	parseEvent(data string) (int, *Event)
}

const protocolLogFlag = "protocol"

//defaultProtocol 封装解析事件
type defaultProtocol struct {
}
type Event struct {
	EType   State
	aUUID   string
	Content string
	From    string
}

type State int

const (
	PING State = 1 + iota
	LEASE
	VOTE
	HEARTBEAT
	JOIN
	ACK

	TASK
	RESPONSE
)

const (
	SFD       = "$+$+$\n"
	EPD       = "$-$-$\n"
	LEN_SPLIT = len("\n")
	LEN_SFD   = len(SFD)
	LEN_EPD   = len(EPD)
)

func (t *defaultProtocol) make(event Event) string {
	builder := strings.Builder{}
	builder.WriteString(SFD)
	builder.WriteString(strconv.Itoa(int(event.EType)) + "\n")
	builder.WriteString(event.aUUID + "\n")
	builder.WriteString(event.Content)
	builder.WriteString(EPD)
	//glog.Infof("[%s]:make event[%s]\n", protocolLogFlag, builder.String())
	return builder.String()
}

func (t *defaultProtocol) parseEvent(data string) (int, *Event) {
	//glog.Infof("[%s]:unparsed data[%s]\n", protocolLogFlag, data)
	end := strings.Index(data, EPD)
	start := strings.LastIndex(data[:end], SFD) + LEN_SFD
	if start == 4 || end == -1 {
		return -1, nil
	}
	data = data[start:end]
	idxState := strings.Index(data, "\n")
	state, _ := strconv.Atoi(data[:idxState])
	data = data[idxState+LEN_SPLIT:]
	idxUUID := strings.Index(data, "\n")
	UUID := data[:idxUUID]
	return end + LEN_EPD, &Event{
		EType:   State(state),
		aUUID:   UUID,
		Content: data[idxUUID+LEN_SPLIT:],
	}
}
