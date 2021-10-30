package tcpconn

import "strings"

type TCPProtocol struct {
}
type Event struct {
	Name string
	Content string
}
const (
	SFD = "♩♪♫♬\n"
	EPD = "$￡￥\n"
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

func (t *TCPProtocol) Make(event Event) string {
	builder := strings.Builder{}
	builder.WriteString(SFD)
	builder.WriteString(event.Name+"\n")
	builder.WriteString(event.Content)
	builder.WriteString(EPD)
	return builder.String()
}

func (t *TCPProtocol) ParseEvent(data string) (int, *Event) {
	end := strings.Index(data, EPD)
	start := strings.LastIndex(data[:end], SFD) + LEN_SFD
	if start == 4 || end == -1 {
		return -1,nil
	}
	result := data[start:end]
	idxName := strings.Index(result, "\n")

	return end+LEN_EPD,&Event{
		Name: result[:idxName],
		Content: result[idxName+1:],
	}
}
