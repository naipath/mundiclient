package mundiclient

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"fmt"
)

const (
	uploadLogo                 = 0x50
	uploadLogoData             = 0x42
	lastUploadLogoData         = 0x46
	requestOverWritePermission = 0x21
	notSaved                   = 0x0D
	failedToCreateLogoFile     = 0x0F
)

func (m *MundiClient) UploadLogo(logo *os.File) error {

	filestatistics, _ := logo.Stat()

	fileSize := uint32(filestatistics.Size())
	fileSizeLength := convertUInt32ToBytes(fileSize)

	fileName := filepath.Base(logo.Name())
	fileNameLength := []byte{byte(len(fileName) >> 8), byte(len(fileName) & 0xff)}

	message := []byte{uploadLogo}
	message = append(message, fileSizeLength...)
	message = append(message, fileNameLength...)
	message = append(message, []byte(fileName)...)

	response, err := m.sendAndReceiveMessage(message)

	if err != nil || response[0] != acknowledge {
		return errors.New("error sending logo")
	}

	b, err := ioutil.ReadAll(logo)
	if err != nil {
		return err
	}

	for i := 0; i < len(b)/500; i++ {
		dataToSend := b[i*500: i*500+500]

		response, err = m.sendAndReceiveMessage(append([]byte{uploadLogoData, 0x01, 0xF4}, dataToSend...))

		if err != nil || response[0] != acknowledge {
			return errors.New("error sending part of logo: " + string(i))
		}
		fmt.Println("Current iterations", i)
		fmt.Println("Total iterations", len(b)/500)
	}
	var totalLogo uint32
	for _, element := range b {
		totalLogo += uint32(element)
	}
	logoChecksum := convertUInt32ToBytes(totalLogo)

	lastData := b[len(b)/500*500: len(b)/500*500+len(b)%500]
	lastDataLength := uint16(len(lastData))
	msbLastMessage, lsbLastMessage := byte(lastDataLength>>8), byte(lastDataLength&0xff)

	lastLogoBlockMessage := []byte{lastUploadLogoData, msbLastMessage, lsbLastMessage}
	lastLogoBlockMessage = append(lastLogoBlockMessage, lastData...)
	lastLogoBlockMessage = append(lastLogoBlockMessage, logoChecksum...)

	response, err = m.sendAndReceiveMessage(lastLogoBlockMessage)

	if response[0] == requestOverWritePermission {
		response, err = m.sendAndReceive([]byte{requestOverWritePermission})
	}
	if response[0] == notSaved {
		return errors.New("Failed to save logo file")
	}
	if response[0] == failedToCreateLogoFile {
		return errors.New("Failed to create logo file")
	}
	if err != nil || response[0] != acknowledge {
		return errors.New("Could not save logo")
	}
	return nil
}
