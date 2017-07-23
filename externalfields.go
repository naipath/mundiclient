package mundiclient

const (
	modifyExternalField = 0x4D
)

func (m MundiClient) ModifyExternalField(fieldID byte, data string) {
	message := append([]byte{modifyExternalField, byte(len(data)) * 2, fieldID}, []byte(data)...)
	response := m.sendAndReceiveMessage(message)

	if response[0] != acknowledge {
		panic("Could not modify external field")
	}
}
