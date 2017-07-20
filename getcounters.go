package mundiclient

import "encoding/binary"

const (
	getCounters = 0x47
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
