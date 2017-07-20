package mundiclient

const (
	getStatusMessage = 0x58
)

func (m MundiClient) GetStatusMessage() string {
	lsb, msb := calculateChecksum(getStatusMessage, emptyLength)
	response := m.sendAndReceive([]byte{startOfText, getStatusMessage, emptyLength, lsb, msb, endOfTransmission})
	length := response[2]

	statusMessage := string(response[3 : length+3])

	return statusMessage
}
