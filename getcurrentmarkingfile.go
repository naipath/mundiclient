package mundiclient

const (
	getCurrentMarkingFile = 0xDB
)

func (m MundiClient) GetCurrentMarkingFile() string {
	lsb, msb := calculateChecksum(getCurrentMarkingFile, emptyLength)
	message := []byte{startOfText, getCurrentMarkingFile, emptyLength, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	return string(response[3 : response[2]+3])
}
