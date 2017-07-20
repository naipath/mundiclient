package mundiclient

import (
	"fmt"
	"net"
	"strconv"
)

const (
	startOfText       = 0x02
	endOfTransmission = 0x04
	emptyLength       = 0x00
)

type MundiClient struct {
	conn net.Conn
}

func New(ip string, port int) *MundiClient {
	return &MundiClient{createConnection(ip, port)}
}

func (m MundiClient) Close() {
	m.conn.Close()
}

func (m MundiClient) sendAndReceive(message []byte) []byte {
	return m.sendAndReceiveWithCustomDelim(message, endOfTransmission)
}

func (m MundiClient) sendAndReceiveWithCustomDelim(message []byte, delim byte) []byte {
	m.conn.Write(message)

	reply := make([]byte, 1024)

	m.conn.Read(reply)

	for i := len(reply) - 1; i >= 0; i-- {
		if reply[i] != 0x00 {
			fmt.Printf("Got the following response:\n%08b\n", reply[:i+1])
			return reply[:i+1]
		}
	}
	panic("Got no response!")
}

func createConnection(ip string, port int) net.Conn {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))

	if err != nil {
		panic(err)
	}
	return conn
}

func calculateChecksum(input ...byte) (byte, byte) {
	var total uint16
	for _, value := range input {
		total += uint16(value)
	}
	lsb := byte(total & 0xFF)
	msb := byte((total >> 8) & 0xFF)
	return lsb, msb
}
