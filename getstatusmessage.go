package mundiclient

func (m MundiClient) GetStatusMessage() string {
	lsb, msb := calculateChecksum(0x58, emptyLength)
	response := m.sendAndReceive([]byte{startOfText, 0x58, emptyLength, lsb, msb, endOfTransmission})
	length := response[2]

	statusMessage := string(response[3 : length+3])

	return statusMessage
}
