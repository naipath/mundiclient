package mundiclient

import (
	"fmt"
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
	fileSizeLength := []byte{byte(fileSize >> 24), byte(fileSize >> 16), byte(fileSize >> 8), byte(fileSize & 0xff)}

	fileName := filepath.Base(logo.Name())
	fileNameLength := []byte{byte(len(fileName) >> 8), byte(len(fileName) & 0xff)}

	checksum := []byte{uploadLogo}
	checksum = append(checksum, fileSizeLength...)
	checksum = append(checksum, fileNameLength...)
	checksum = append(checksum, []byte(fileName)...)

	lsb, msb := calculateChecksum(checksum...)

	message := []byte{startOfText, uploadLogo}
	message = append(message, fileSizeLength...)
	message = append(message, fileNameLength...)
	message = append(message, []byte(fileName)...)
	message = append(message, lsb, msb, endOfTransmission)

	response := m.sendAndReceive(message)

	if response[0] != acknowledge {
		panic("error sending logo")
	}

	b, err := ioutil.ReadAll(logo)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(b)/500; i++ {

		dataToSend := b[i*500 : i*500+500]

		var checksum uint16 = uploadLogoData
		for _, element := range dataToSend {
			checksum += uint16(element)
		}
		checksum += 0x01 + 0xF4

		logoDataMessage := []byte{startOfText, uploadLogoData, 0x01, 0xF4}
		logoDataMessage = append(logoDataMessage, dataToSend...)
		logoDataMessage = append(logoDataMessage, byte(checksum&0xFF), byte(checksum>>8), endOfTransmission)

		fmt.Println("Now at ", i*500)
		fmt.Println("Total is ", fileSize)

		response = m.sendAndReceive(logoDataMessage)

		if response[0] != acknowledge {
			panic("error sending part of logo: " + string(i))
		}
	}

	var totalLogo uint32
	for _, element := range b {
		totalLogo += uint32(element)
	}
	logoChecksum := []byte{byte(totalLogo >> 24), byte(totalLogo >> 16), byte(totalLogo >> 8), byte(totalLogo & 0xff)}

	lastData := b[len(b)/500*500 : len(b)/500*500+len(b)%500]
	lastDataLength := uint16(len(lastData))
	msbLastMessage, lsbLastMessage := byte(lastDataLength>>8), byte(lastDataLength&0xff)

	lastDataChecksum := []byte{lastUploadLogoData, logoChecksum[0], logoChecksum[1], logoChecksum[2], logoChecksum[3], msbLastMessage, lsbLastMessage}
	lastDatalsb, lastDatamsb := calculateChecksum(append(lastDataChecksum, lastData...)...)

	lastLogoBlockMessage := []byte{startOfText, lastUploadLogoData, msbLastMessage, lsbLastMessage}
	lastLogoBlockMessage = append(lastLogoBlockMessage, lastData...)
	lastLogoBlockMessage = append(lastLogoBlockMessage, logoChecksum...)
	lastLogoBlockMessage = append(lastLogoBlockMessage, []byte{lastDatalsb, lastDatamsb, endOfTransmission}...)

	response = m.sendAndReceive(lastLogoBlockMessage)
	if response[0] != acknowledge {
		panic("could not save logo")
	}
}
