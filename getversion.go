package mundiclient

const (
	getVersionID = 0x56
)

func (m MundiClient) GetVersion() string {
	lsb, msb := calculateChecksum(getVersionID, emptyLength)
	message := []byte{startOfText, getVersionID, emptyLength, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	return string(response[3:9])
}
