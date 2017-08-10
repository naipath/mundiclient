package mundiclient

import (
	"encoding/binary"
	"errors"
)

const (
	getLineSettings = 0xA7
	setLineSettings = 0xDD

	productPresentNotAvailableWithStaticMarking = 0xC0
	tachoOrMaxLineSpeedRequired                 = 0xC1
	productIntervalTooSmall                     = 0xC2

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

func (m *MundiClient) GetLineSettings() (LineSettings, error) {
	response, err := m.sendAndReceiveMessage([]byte{getLineSettings, emptyLength})
	if err != nil {
		return LineSettings{}, err
	}
	return LineSettings{
		response[4],
		detectorType(response[6]),
		binary.BigEndian.Uint16([]byte{response[9], response[8]}),
	}, nil
}

func (m *MundiClient) SetLineSettingsDelay(delay uint16) error {
	var length byte = 0x03
	var setLineSettingsDelayID byte = 0xC4
	msbDelay, lsbDelay := byte(delay>>8), byte(delay&0xFF)

	message := []byte{setLineSettings, length, setLineSettingsDelayID, lsbDelay, msbDelay}
	response, err := m.sendAndReceiveMessage(message)

	if err != nil {
		return err
	}

	if response[0] != acknowledge {
		if response[0] == negativeAcknowledge {
			switch response[1] {
			case productPresentNotAvailableWithStaticMarking:
				return errors.New("Product present not available with static marking")
			case tachoOrMaxLineSpeedRequired:
				return errors.New("Tacho or max line speed required")
			case productIntervalTooSmall:
				return errors.New("Product interval too small")
			}
		}
		return errors.New("Could not set LineSettings")
	}
	return nil
}
