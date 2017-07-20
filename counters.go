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

	lsb, msb := calculateChecksum(getCounters, emptyLength)
	message := []byte{startOfText, getCounters, emptyLength, lsb, msb, endOfTransmission}
	response := m.sendAndReceive(message)

	lifetime := binary.BigEndian.Uint32(response[3:7])
	recent := binary.BigEndian.Uint32(response[7:11])

	return Counters{lifetime, recent}
}

func (m MundiClient) ResetCurrentCount() {
	lsb, msb := calculateChecksum(resetCurrentCount, emptyLength)
	message := []byte{startOfText, resetCurrentCount, emptyLength, lsb, msb, endOfTransmission}
	response := m.sendAndReceive(message)

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
	lsb, msb := calculateChecksum(getIncrementalCounterValue, length, fieldId)
	message := []byte{startOfText, getIncrementalCounterValue, length, fieldId, lsb, msb, endOfTransmission}
	response := m.sendAndReceive(message)

	return IncrementalCounterValue{
		response[3],
		string(response[5 : response[4]*2+5]),
	}
}

func (m MundiClient) SetIncrementalCounterValue(input IncrementalCounterValue) {
	length := byte(len(input.Data)) + 0x2
	lsb, msb := calculateChecksum(0x00)
	startOfMessage := []byte{
		startOfText,
		setIncrementalCounterValue,
		length,
		input.FieldID,
		byte(len(input.Data)),
	}
	endOfMessage := []byte{lsb, msb, endOfTransmission}

	message := append(append(startOfMessage, []byte(input.Data)...), endOfMessage...)
	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("Could not set incremental counter")
	}
}
