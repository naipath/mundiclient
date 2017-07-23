package mundiclient

const (
	getVersionID = 0x56
)

func (m MundiClient) GetVersion() string {
	response := m.sendAndReceive(constructMessage([]byte{getVersionID, emptyLength}))
	return string(response[3:9])
}
