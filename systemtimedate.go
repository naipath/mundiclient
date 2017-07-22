package mundiclient

import "time"

const (
	getSystemTime = 0xA8
	setSystemTime = 0xAA
	getSystemDate = 0xAB
)

func (m MundiClient) GetSystemTime() time.Time {
	systemTime := m.getSystemTime()
	systemDate := m.getSystemDate()

	then, _ := time.Parse("02012006150405", string(systemDate[3:11])+string(systemTime[3:9]))

	return then
}

func (m MundiClient) getSystemTime() []byte {
	lsb, msb := calculateChecksum(getSystemTime, emptyLength)
	message := []byte{startOfText, getSystemTime, emptyLength, lsb, msb, endOfTransmission}
	return m.sendAndReceive(message)
}

func (m MundiClient) getSystemDate() []byte {
	lsb, msb := calculateChecksum(getSystemDate, emptyLength)
	message := []byte{startOfText, getSystemDate, emptyLength, lsb, msb, endOfTransmission}
	return m.sendAndReceive(message)
}
