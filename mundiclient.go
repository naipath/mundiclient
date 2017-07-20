package mundiclient

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	m.conn.Write(message)
	result, err := bufio.NewReader(m.conn).ReadBytes(endOfTransmission)
	fmt.Printf("Got the following response:\n%08b\n", result)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}

func createConnection(ip string, port int) net.Conn {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))

	if err != nil {
		fmt.Print("error opening connection\n", err)
		os.Exit(1)
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
