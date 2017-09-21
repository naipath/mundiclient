package mundiclient

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	startOfText         = 0x02
	endOfTransmission   = 0x04
	emptyLength         = 0x00
	acknowledge         = 0x06
	negativeAcknowledge = 0x15
)

type MundiClient struct {
	conn        net.Conn
	debug       bool
	lock        sync.Mutex
	readTimeout time.Duration
}

func New(ip string, port int) (*MundiClient, error) {
	conn, err := createConnection(ip, port)
	if err != nil {
		return nil, err
	}
	return &MundiClient{conn, false, sync.Mutex{}, 10 * time.Second}, nil
}

func (m *MundiClient) Close() error {
	return m.conn.Close()
}

func (m *MundiClient) SetDebug(debug bool) {
	m.debug = debug
}

func (m *MundiClient) SetReadTimeout(readTimeout time.Duration) {
	m.readTimeout = readTimeout
}

func (m *MundiClient) sendAndReceiveMessage(message []byte) ([]byte, error) {
	return m.sendAndReceive(constructMessage(message))
}

func (m *MundiClient) sendAndReceive(bytes []byte) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.debug {
		fmt.Printf("Sending the following request:\n%08b\n", bytes)
	}
	m.conn.Write(bytes)

	m.conn.SetReadDeadline(time.Now().Add(m.readTimeout))

	reply := make([]byte, 1024)
	m.conn.Read(reply)

	for i := len(reply) - 1; i >= 0; i-- {
		if reply[i] != 0x00 {
			if m.debug {
				fmt.Printf("Got the following response:\n%08b\n", reply[:i+1])
			}
			return reply[:i+1], nil
		}
	}
	return nil, errors.New("sendAndReceive: Did not receive response")
}

func createConnection(ip string, port int) (net.Conn, error) {
	var lastErr error
	for i := 0; i < 10; i++ {
		conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
		if err == nil {
			return conn, nil
		}
		lastErr = err
	}
	return nil, lastErr
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

func convertUInt32ToBytes(input uint32) []byte {
	return []byte{byte(input >> 24), byte(input >> 16), byte(input >> 8), byte(input & 0xff)}
}

func constructMessage(input []byte) []byte {
	lsb, msb := calculateChecksum(input...)
	return append(append([]byte{startOfText}, input...), lsb, msb, endOfTransmission)
}
