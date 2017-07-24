package mundiclient

import (
	"encoding/binary"
	"errors"
)

const (
	selectMarkingFile     = 0x98
	getCurrentMarkingFile = 0xDB
	getMarkingFiles       = 0x99
	downloadFile          = 0xD3
)

func (m MundiClient) SelectMarkingFile(filename string) error {
	data := []byte(filename)
	length := byte(len(data))

	message := append([]byte{selectMarkingFile, length}, data...)
	response, err := m.sendAndReceiveMessage(message)

	if err != nil || response[0] != acknowledge {
		return errors.New("could not select marking file")
	}
	return nil
}

func (m MundiClient) GetCurrentMarkingFile() (string, error) {
	response, err := m.sendAndReceiveMessage([]byte{getCurrentMarkingFile, emptyLength})
	if err != nil {
		return "", err
	}
	return string(response[3 : response[2]+3]), nil
}

func (m MundiClient) GetMarkingFiles() ([]string, error) {
	response, err := m.sendAndReceiveMessage([]byte{getMarkingFiles, emptyLength})

	if err != nil {
		return nil, err
	}

	var markingFiles []string
	for {
		if response[0] == negativeAcknowledge {
			return nil, errors.New("error occured GetMarkingFiles")
		}
		if response[0] == 0x17 {
			return markingFiles, nil
		}
		markingFiles = append(markingFiles, string(response[5:response[4]+5]))
		response, err = m.sendAndReceive([]byte{0x21})

		if err != nil {
			return nil, err
		}
	}
}

func (m MundiClient) DownloadFile(markingfilename string) (string, []byte, error) {
	length := byte(len(markingfilename))
	data := []byte(markingfilename)

	response, err := m.sendAndReceiveMessage(append([]byte{downloadFile, length}, data...))

	if err != nil || response[0] != acknowledge {
		return "", nil, errors.New("Did not get acknowledge for download")
	}

	fileSize := binary.BigEndian.Uint32(response[3:8])
	fileNameLength := binary.BigEndian.Uint16(response[7:9])
	fileName := string(response[9 : 9+fileNameLength])

	totalBytes := []byte{}
	for {
		response, err = m.sendAndReceive([]byte{acknowledge})

		if err != nil {
			return "", nil, err
		}

		dataLength := binary.BigEndian.Uint16(response[2:4])
		totalBytes = append(totalBytes, response[4:dataLength+4]...)

		if dataLength != 500 {
			break
		}
	}
	if len(totalBytes) != int(fileSize) {
		return "", nil, errors.New("Did not receive correct amount of bytes")
	}
	return fileName, totalBytes, nil
}
