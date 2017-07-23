package mundiclient

const (
	offsetField = 0x60
	allFieldsID = 0x00
)

func (m MundiClient) OffsetAllFields(x int16, y int16) {
	m.offsetField(allFieldsID, x, y)
}

func (m MundiClient) OffsetField(fieldId byte, x int16, y int16) {
	m.offsetField(fieldId, x, y)
}

func (m MundiClient) offsetField(fieldID byte, x int16, y int16) {
	var length byte = 0x05
	msbX, lsbX := byte(x>>8), byte(x&0xFF)
	msbY, lsbY := byte(y>>8), byte(y&0xFF)

	message := constructMessage([]byte{offsetField, length, fieldID, msbX, lsbX, msbY, lsbY})

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("Could not alter offset of field")
	}
}
