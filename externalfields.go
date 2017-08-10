package mundiclient

import "errors"

const (
	modifyExternalField = 0x4D
)

func (m *MundiClient) ModifyExternalField(fieldID byte, data string) error {
	message := append([]byte{modifyExternalField, byte(len(data)) * 2, fieldID}, []byte(data)...)
	response, err := m.sendAndReceiveMessage(message)

	if err != nil || response[0] != acknowledge {
		return errors.New("Could not modify external field")
	}
	return nil
}
