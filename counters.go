package mundiclient

import (
	"encoding/binary"
	"errors"
)

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

func (m MundiClient) GetCounters() (Counters, error) {
	response, err := m.sendAndReceiveMessage([]byte{getCounters, emptyLength})
	if err != nil {
		return Counters{}, err
	}

	lifetime := binary.BigEndian.Uint32(response[3:7])
	recent := binary.BigEndian.Uint32(response[7:11])

	return Counters{lifetime, recent}, nil
}

func (m MundiClient) ResetCurrentCount() error {
	response, err := m.sendAndReceiveMessage([]byte{resetCurrentCount, emptyLength})

	if err != nil || response[0] != acknowledge {
		return errors.New("reset current count failed")
	}
	return nil
}

type IncrementalCounterValue struct {
	FieldID byte
	Data    string
}

func (m MundiClient) GetIncrementalCounterValue(fieldId byte) (IncrementalCounterValue, error) {
	var length byte = 0x01
	response, err := m.sendAndReceiveMessage([]byte{getIncrementalCounterValue, length, fieldId})

	if err != nil {
		return IncrementalCounterValue{}, err
	}

	return IncrementalCounterValue{
		response[3],
		string(response[5 : response[4]*2+5]),
	}, nil
}

func (m MundiClient) SetIncrementalCounterValue(input IncrementalCounterValue) error {
	length := byte(len(input.Data)) + 0x2 // Double check for utf-8 to ansi conversion

	startOfMessage := []byte{setIncrementalCounterValue, length, input.FieldID, byte(len(input.Data))}
	message := append(startOfMessage, []byte(input.Data)...)

	response, err := m.sendAndReceiveMessage(message)

	if err != nil || response[0] != acknowledge {
		return errors.New("Could not set incremental counter")
	}
	return nil
}
