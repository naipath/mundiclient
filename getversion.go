package mundiclient

const (
	getVersionID = 0x56
)

func (m MundiClient) GetVersion() (string, error) {
	response, err := m.sendAndReceiveMessage([]byte{getVersionID, emptyLength})
	if err != nil {
		return "", err
	}
	return string(response[3:9]), nil
}
