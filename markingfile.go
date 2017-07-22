package mundiclient

const (
	getCurrentMarkingFile = 0xDB
	getMarkingFiles       = 0x99
)

func (m MundiClient) GetCurrentMarkingFile() string {
	lsb, msb := calculateChecksum(getCurrentMarkingFile, emptyLength)
	message := []byte{startOfText, getCurrentMarkingFile, emptyLength, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	return string(response[3 : response[2]+3])
}

func (m MundiClient) GetMarkingFiles() []string {
	lsb, msb := calculateChecksum(getMarkingFiles, emptyLength)

	message := []byte{startOfText, getMarkingFiles, emptyLength, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	var markingFiles []string

	for {
		if response[0] == 0x15 {
			panic("error occured GetMarkingFiles")
		}
		if response[0] == 0x17 {
			return markingFiles
		}
		markingFiles = append(markingFiles, string(response[5:response[4]+5]))
		response = m.sendAndReceive([]byte{0x21})
	}
}
