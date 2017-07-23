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
	return m.sendAndReceiveMessage([]byte{getSystemTime, emptyLength})
}

func (m MundiClient) getSystemDate() []byte {
	return m.sendAndReceiveMessage([]byte{getSystemDate, emptyLength})
}
