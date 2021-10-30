package tcpconn

import (
	"Cluster/handler"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
)

type defaultDecoder struct {
	protocol TCPProtocol
	IPPort   string
}

const (
	BUFFER_SIZE = 128
)

//读取数据，并将其封装成事件
func (d *defaultDecoder) createEventFromConn(conn net.Conn) {
	d.IPPort = conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	buffer := bytes.Buffer{}
	var block [BUFFER_SIZE]byte
	for {
		n, err := reader.Read(block[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err: ", err)
			break
		}
		d.splitStickPackage(&buffer, block, n)
	}
}

//解析粘包，将解析出的数据包封装成事件
func (d *defaultDecoder) splitStickPackage(buffer *bytes.Buffer, block [BUFFER_SIZE]byte, dataSize int) {
	start := 0
	for strings.Contains(string(block[start:dataSize]), EPD) {
		end := start + strings.Index(string(block[start:dataSize]), EPD) + LEN_EPD
		buffer.Write(block[start:end])
		_, event := d.protocol.ParseEvent(buffer.String())
		event.From = d.IPPort
		handler.DoDistribute(*event)
		buffer.Reset()
		start = end
	}
	if start == 0 {
		buffer.Write(block[:dataSize])
	} else {
		buffer.Write(block[start:dataSize])
	}
}
