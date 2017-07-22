package mundiclient

import "encoding/binary"

const (
	getLineSettings = 0xA7

	Single1       = detectorType(0)
	Single2       = detectorType(1)
	Dual1         = detectorType(2)
	Dual2         = detectorType(3)
	DualBothLong  = detectorType(5)
	DualBothShort = detectorType(6)
	DualEither    = detectorType(7)
)

type detectorType byte

type LineSettings struct {
	Direction byte
	Detector  detectorType
	Delay     uint16
}

func (m MundiClient) GetLineSettings() LineSettings {
	lsb, msb := calculateChecksum(getLineSettings, emptyLength)
	message := []byte{startOfText, getLineSettings, emptyLength, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	return LineSettings{
		response[4],
		detectorType(response[6]),
		binary.BigEndian.Uint16([]byte{response[9], response[8]}),
	}
}
