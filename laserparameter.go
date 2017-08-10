package mundiclient

import (
	"encoding/binary"
	"errors"
)

const (
	getLaserParameter        = 0xA6
	setLaserParameter        = 0xDC
	setLaserParameterGroupID = 0xE0
)

type LaserParameter struct {
	Group        byte   // Group: 0-255
	Frequency    uint16 // Frequency: 500-50000 Hz
	Duty         uint16 // Duty: 1-85%
	MarkSpeed    uint16 // MarkSpeed: 10-20,000 mm/s
	JumpSpeed    uint16 // JumpSpeed: 10-20,000 mm/s
	OffDelay     uint16 // OffDelay: 2-8000 μSeconds
	OnDelay      uint16 // OnDelay: -8000 to 8000 μSeconds
	JumpDelay    uint16 // JumpDelay: 0-32000 μSeconds
	MarkDelay    uint16 // MarkDelay: 0-32000 μSeconds
	PolygonDelay uint16 // PolygonDelay: 0-32000 μSeconds
}

func (m *MundiClient) GetLaserParameter(group byte) (LaserParameter, error) {
	response, err := m.sendAndReceiveMessage([]byte{getLaserParameter, 0x01, group})

	if err != nil {
		return LaserParameter{}, err
	}
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
	}, nil
}

func (m *MundiClient) SetLaserParameterDuty(group byte, duty byte) error {
	var dutyID byte = 0xE2
	var length byte = 0x06

	message := []byte{setLaserParameter, length, setLaserParameterGroupID, group, 0, dutyID, duty, 0}
	response, err := m.sendAndReceiveMessage(message)

	if err != nil || response[0] != acknowledge {
		return errors.New("Could not set Laser Parameter")
	}
	return nil
}

func (m *MundiClient) SetLaserParameterFrequency(group byte, frequency uint16) error {
	var frequencyID byte = 0xE1
	var length byte = 0x06

	msbFrequency, lsbFrequency := byte(frequency>>8), byte(frequency&0xff)

	message := []byte{setLaserParameter, length, setLaserParameterGroupID, group, 0, frequencyID, lsbFrequency, msbFrequency}
	response, err := m.sendAndReceiveMessage(message)

	if err != nil || response[0] != acknowledge {
		return errors.New("Could not set Laser Parameter")
	}
	return nil
}
