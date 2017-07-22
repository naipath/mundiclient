package mundiclient

const (
	offsetField = 0x60
)

func (m MundiClient) OffsetAllFields(x int16, y int16) {
	m.offsetField(0x00, x, y)
}

func (m MundiClient) OffsetField(fieldId byte, x int16, y int16) {
	m.offsetField(fieldId, x, y)
}

func (m MundiClient) offsetField(fieldID byte, x int16, y int16) {
	msbX, lsbX := byte(x>>8), byte(x&0xff)
	msbY, lsbY := byte(y>>8), byte(y&0xff)

	lsb, msb := calculateChecksum(offsetField, 0x05, fieldID, msbX, msbY, lsbX, lsbY)
	message := []byte{startOfText, offsetField, 0x05, fieldID, msbX, lsbX, msbY, lsbY, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("Could not alter offset of field")
	}
}
