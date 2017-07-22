package mundiclient

import "encoding/binary"

const (
	getLaserParameter = 0xA6
)

type LaserParameter struct {
	Group        uint16 // Range: 0-255
	Frequency    uint16 // Range: 500-50000 Hz
	Duty         uint16 // Range: 1-85%
	MarkSpeed    uint16 // Range: 10-20,000 mm/s
	JumpSpeed    uint16 // Range: 10-20,000 mm/s
	OffDelay     uint16 // Range: 2-8000 μSeconds
	OnDelay      uint16 // Range: -8000 to 8000 μSeconds
	JumpDelay    uint16 // Range: 0-32000 μSeconds
	MarkDelay    uint16 // Range: 0-32000 μSeconds
	PolygonDelay uint16 // Range: 0-32000 μSeconds
}

func (m MundiClient) GetLaserParameter(parameterGroup byte) LaserParameter {
	lsb, msb := calculateChecksum(getLaserParameter, 0x01, parameterGroup)
	message := []byte{startOfText, getLaserParameter, 0x01, parameterGroup, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	return LaserParameter{
		binary.BigEndian.Uint16([]byte{response[5], response[4]}),
		binary.BigEndian.Uint16([]byte{response[8], response[7]}),
		binary.BigEndian.Uint16([]byte{response[11], response[10]}),
		binary.BigEndian.Uint16([]byte{response[14], response[13]}),
		binary.BigEndian.Uint16([]byte{response[17], response[16]}),
		binary.BigEndian.Uint16([]byte{response[20], response[19]}),
		binary.BigEndian.Uint16([]byte{response[23], response[22]}),
		binary.BigEndian.Uint16([]byte{response[26], response[25]}),
		binary.BigEndian.Uint16([]byte{response[29], response[28]}),
		binary.BigEndian.Uint16([]byte{response[32], response[31]}),
	}
}
