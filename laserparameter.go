package mundiclient

import "encoding/binary"

const (
	getLaserParameter = 0xA6
)

type LaserParameter struct {
	Group        uint16
	Frequency    uint16
	Duty         uint16
	MarkSpeed    uint16
	JumpSpeed    uint16
	OffDelay     uint16
	OnDelay      uint16
	JumpDelay    uint16
	MarkDelay    uint16
	PolygonDelay uint16
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
