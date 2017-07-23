package mundiclient

import "encoding/binary"

const (
	selectMarkingFile     = 0x98
	getCurrentMarkingFile = 0xDB
	getMarkingFiles       = 0x99
	downloadFile          = 0xD3
)

func (m MundiClient) SelectMarkingFile(filename string) {
	data := []byte(filename)
	length := byte(len(data))

	lsb, msb := calculateChecksum(append([]byte{selectMarkingFile, length}, data...)...)

	startOfMessage := []byte{startOfText, selectMarkingFile, length}
	endOfMessage := []byte{lsb, msb, endOfTransmission}

	message := append(append(startOfMessage, data...), endOfMessage...)

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("could not select marking file")
	}
}

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

func (m MundiClient) DownloadFile(markingfilename string) (string, []byte) {
	length := byte(len(markingfilename))
	data := []byte(markingfilename)

	checksum := []byte{downloadFile, length}
	checksum = append(checksum, data...)
	lsb, msb := calculateChecksum(checksum...)

	message := []byte{startOfText, downloadFile, length}
	message = append(message, data...)
	message = append(message, lsb, msb, endOfTransmission)

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("Did not get acknowledge for download")
	}

	fileSize := binary.BigEndian.Uint32(response[3:8])
	fileNameLength := binary.BigEndian.Uint16(response[7:9])
	fileName := string(response[9 : 9+fileNameLength])

	totalBytes := []byte{}
	for {
		response = m.sendAndReceive([]byte{0x06})

		dataLength := binary.BigEndian.Uint16(response[2:4])
		totalBytes = append(totalBytes, response[4:dataLength+4]...)

		if dataLength != 500 {
			break
		}
	}

	if len(totalBytes) != int(fileSize) {
		panic("Did not receive correct amount of bytes")
	}
	return fileName, totalBytes
}
