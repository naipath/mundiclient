package mundiclient

const (
	modifyExternalField = 0x4D
)

func (m MundiClient) ModifyExternalField(fieldID byte, data string) {

	length := byte(len(data)) * 2

	startOfMessage := []byte{modifyExternalField, length, fieldID}
	message := constructMessage(append(startOfMessage, []byte(data)...))

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("Could not modify external field")
	}
}
