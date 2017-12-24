package mundiclient

import (
	"net"
	"testing"
)

type handleResponse func(net.Conn, []byte)

func equals(value1 interface{}, value2 interface{}, t *testing.T) {
	if value1 != value2 {
		t.Error("Value", value1, "is not equal to", value2)
	}
}

func notEquals(value1 interface{}, value2 interface{}, t *testing.T) {
	if value1 == value2 {
		t.Error("Value", value1, "is equal to", value2)
	}
}

func setupMockedListener(t *testing.T, handle handleResponse) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:50001")
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				t.Error(err)
			}
			reply := make([]byte, 1024)
			conn.Read(reply)
			handle(conn, reply)
		}
	}()

	return l
}

func TestCreatingNewMundiClient(t *testing.T) {
	l := setupMockedListener(t, func(conn net.Conn, reply []byte) {})
	defer l.Close()
	client, err := New("127.0.0.1", 50001)

	equals(err, nil, t)
	notEquals(client, nil, t)
}

func TestGetVersion(t *testing.T) {
	l := setupMockedListener(t, func(conn net.Conn, reply []byte) {
		equals(reply[1], byte(getVersionID), t)
		equals(reply[2], byte(emptyLength), t)
		response := append([]byte{0, 0, 0}, "AAAAAA"...)
		conn.Write(response)
	})
	defer l.Close()
	client, _ := New("127.0.0.1", 50001)

	version, _ := client.GetVersion()
	equals(version, "AAAAAA", t)
}
