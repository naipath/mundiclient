package mundiclient

import "encoding/binary"

const (
	getLineSettings = 0xA7
	setLineSettings = 0xDD

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
	response := m.sendAndReceiveMessage([]byte{getLineSettings, emptyLength})
	return LineSettings{
		response[4],
		detectorType(response[6]),
		binary.BigEndian.Uint16([]byte{response[9], response[8]}),
	}
}

func (m MundiClient) SetLineSettingsDelay(delay uint16) {
	var length byte = 0x03
	var setLineSettingsDelayID byte = 0xC4
	msbDelay, lsbDelay := byte(delay>>8), byte(delay&0xFF)

	message := []byte{setLineSettings, length, setLineSettingsDelayID, lsbDelay, msbDelay}
	response := m.sendAndReceiveMessage(message)

	if response[0] != acknowledge {
		if response[0] == 0x15 {
			switch response[1] {
			case 0xC0:
				panic("Product present not available with static marking")
			case 0xC1:
				panic("Tacho or max line speed required")
			case 0xC2:
				panic("Product interval too small")
			}
		}
		panic("Could not set LineSettings")
	}
}
