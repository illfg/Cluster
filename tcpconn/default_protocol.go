package tcpconn

import (
	"Cluster/handler"
	"strconv"
	"strings"
)

type TCPProtocol struct {
}

const (
	SFD     = "$+$+$\n"
	EPD     = "$-$-$\n"
	LEN_SFD = len(SFD)
	LEN_EPD = len(EPD)
)

/**
“♩♪♫♬\n”
“name\n”
“value\n”
“value\n”
"content" (json)
“$￡￥\n”
*/

func (t *TCPProtocol) Make(event handler.Event) string {
	builder := strings.Builder{}
	builder.WriteString(SFD)
	builder.WriteString(string(event.EType) + "\n")
	builder.WriteString(event.UUID + "\n")
	builder.WriteString(event.Content)
	builder.WriteString(EPD)
	return builder.String()
}

func (t *TCPProtocol) ParseEvent(data string) (int, *handler.Event) {
	end := strings.Index(data, EPD)
	start := strings.LastIndex(data[:end], SFD) + LEN_SFD
	if start == 4 || end == -1 {
		return -1, nil
	}
	result := data[start:end]
	idxName := strings.Index(result, "\n")
	integer, _ := strconv.Atoi(result[:idxName])
	return end + LEN_EPD, &handler.Event{
		EType:   handler.State(integer),
		Content: result[idxName+1:],
	}
}
