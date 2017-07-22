package mundiclient

import "encoding/binary"

const (
	getLaserParameter        = 0xA6
	setLaserParameter        = 0xDC
	setLaserParameterGroupID = 0xE0
)

type LaserParameter struct {
	Group        byte   // Range: 0-255
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

func (m MundiClient) GetLaserParameter(group byte) LaserParameter {
	var length byte = 0x01
	lsb, msb := calculateChecksum(getLaserParameter, length, group)
	message := []byte{startOfText, getLaserParameter, length, group, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	return LaserParameter{
		response[4],
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

func (m MundiClient) SetLaserParameterDuty(group byte, duty byte) {

	var dutyID byte = 0xE2
	var length byte = 0x06

	lsb, msb := calculateChecksum(setLaserParameter, length, setLaserParameterGroupID, group, dutyID, duty)
	message := []byte{startOfText, setLaserParameter, length, setLaserParameterGroupID, group, 0, dutyID, duty, 0, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("Could not set Laser Parameter")
	}
}

func (m MundiClient) SetLaserParameterFrequency(group byte, frequency uint16) {

	var frequencyID byte = 0xE1
	var length byte = 0x06

	msbFrequency, lsbFrequency := byte(frequency>>8), byte(frequency&0xff)

	lsb, msb := calculateChecksum(setLaserParameter, length, setLaserParameterGroupID, group, frequencyID, msbFrequency, lsbFrequency)
	message := []byte{startOfText, setLaserParameter, length, setLaserParameterGroupID, group, 0, frequencyID, lsbFrequency, msbFrequency, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("Could not set Laser Parameter")
	}
}
