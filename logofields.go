package mundiclient

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	uploadLogo         = 0x50
	uploadLogoData     = 0x42
	lastUploadLogoData = 0x46
)

func (m MundiClient) UploadLogo(logo *os.File) {

	filestatistics, _ := logo.Stat()

	fileSize := uint32(filestatistics.Size())
	fileSizeLength := convertUInt32ToBytes(fileSize)

	fileName := filepath.Base(logo.Name())
	fileNameLength := []byte{byte(len(fileName) >> 8), byte(len(fileName) & 0xff)}

	message := []byte{uploadLogo}
	message = append(message, fileSizeLength...)
	message = append(message, fileNameLength...)
	message = append(message, []byte(fileName)...)

	response := m.sendAndReceiveMessage(message)

	if response[0] != acknowledge {
		panic("error sending logo")
	}

	b, err := ioutil.ReadAll(logo)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(b)/500; i++ {
		dataToSend := b[i*500 : i*500+500]

		response = m.sendAndReceiveMessage(append([]byte{uploadLogoData, 0x01, 0xF4}, dataToSend...))

		if response[0] != acknowledge {
			panic("error sending part of logo: " + string(i))
		}
	}

	var totalLogo uint32
	for _, element := range b {
		totalLogo += uint32(element)
	}
	logoChecksum := convertUInt32ToBytes(totalLogo)

	lastData := b[len(b)/500*500 : len(b)/500*500+len(b)%500]
	lastDataLength := uint16(len(lastData))
	msbLastMessage, lsbLastMessage := byte(lastDataLength>>8), byte(lastDataLength&0xff)

	lastLogoBlockMessage := []byte{lastUploadLogoData, msbLastMessage, lsbLastMessage}
	lastLogoBlockMessage = append(lastLogoBlockMessage, lastData...)
	lastLogoBlockMessage = append(lastLogoBlockMessage, logoChecksum...)

	response = m.sendAndReceiveMessage(lastLogoBlockMessage)
	if response[0] != acknowledge {
		panic("could not save logo")
	}
}
