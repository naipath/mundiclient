package mundiclient

import "time"

const (
	getSystemTime = 0xA8
	getSystemDate = 0xAB
)

func (m MundiClient) GetSystemTime() time.Time {
	lsb, msb := calculateChecksum(getSystemTime, emptyLength)
	message := []byte{startOfText, getSystemTime, emptyLength, lsb, msb, endOfTransmission}

	systemTime := m.sendAndReceive(message)

	lsbDate, msbDate := calculateChecksum(getSystemDate, emptyLength)
	messageDate := []byte{startOfText, getSystemDate, emptyLength, lsbDate, msbDate, endOfTransmission}
	systemDate := m.sendAndReceive(messageDate)

	then, _ := time.Parse("02012006150405", string(systemDate[3:11])+string(systemTime[3:9]))

	return then
}
