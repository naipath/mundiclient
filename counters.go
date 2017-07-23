package mundiclient

import "encoding/binary"

const (
	getCounters                = 0x47
	resetCurrentCount          = 0x48
	getIncrementalCounterValue = 0x51
	setIncrementalCounterValue = 0x52
)

type Counters struct {
	Lifetime uint32
	Recent   uint32
}

func (m MundiClient) GetCounters() Counters {
	response := m.sendAndReceiveMessage([]byte{getCounters, emptyLength})

	lifetime := binary.BigEndian.Uint32(response[3:7])
	recent := binary.BigEndian.Uint32(response[7:11])

	return Counters{lifetime, recent}
}

func (m MundiClient) ResetCurrentCount() {
	response := m.sendAndReceiveMessage([]byte{resetCurrentCount, emptyLength})

	if response[0] != acknowledge {
		panic("Reset not acknowledged")
	}
}

type IncrementalCounterValue struct {
	FieldID byte
	Data    string
}

func (m MundiClient) GetIncrementalCounterValue(fieldId byte) IncrementalCounterValue {
	var length byte = 0x01
	response := m.sendAndReceiveMessage([]byte{getIncrementalCounterValue, length, fieldId})

	return IncrementalCounterValue{
		response[3],
		string(response[5 : response[4]*2+5]),
	}
}

func (m MundiClient) SetIncrementalCounterValue(input IncrementalCounterValue) {
	length := byte(len(input.Data)) + 0x2 // Double check for utf-8 to ansi conversion

	startOfMessage := []byte{setIncrementalCounterValue, length, input.FieldID, byte(len(input.Data))}
	message := append(startOfMessage, []byte(input.Data)...)

	response := m.sendAndReceiveMessage(message)

	if response[0] != acknowledge {
		panic("Could not set incremental counter")
	}
}
