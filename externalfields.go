package mundiclient

const (
	modifyExternalField = 0x4D
)

func (m MundiClient) ModifyExternalField(fieldID byte, data string) {

	length := byte(len(data)) * 2

	lsb, msb := calculateChecksum(append([]byte{modifyExternalField, length}, []byte(data)...)...)
	startOfMessage := []byte{startOfText, modifyExternalField, length, fieldID}
	endOfMessage := []byte{lsb, msb, endOfTransmission}

	message := append(append(startOfMessage, []byte(data)...), endOfMessage...)

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("Could not modify external field")
	}
}
