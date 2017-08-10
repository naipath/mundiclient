package mundiclient

import "errors"

const (
	offsetField = 0x60
	allFieldsID = 0x00
)

func (m *MundiClient) OffsetAllFields(x int16, y int16) error {
	return m.offsetField(allFieldsID, x, y)
}

func (m *MundiClient) OffsetField(fieldId byte, x int16, y int16) error {
	return m.offsetField(fieldId, x, y)
}

func (m *MundiClient) offsetField(fieldID byte, x int16, y int16) error {
	var length byte = 0x05
	msbX, lsbX := byte(x>>8), byte(x&0xFF)
	msbY, lsbY := byte(y>>8), byte(y&0xFF)

	message := []byte{offsetField, length, fieldID, msbX, lsbX, msbY, lsbY}

	response, err := m.sendAndReceiveMessage(message)

	if err != nil || response[0] != acknowledge {
		return errors.New("Could not alter offset of field")
	}
	return nil
}
