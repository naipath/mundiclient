package mundiclient

import "time"

const (
	getSystemTime = 0xA8
	getSystemDate = 0xAB
)

func (m *MundiClient) GetSystemTime() (time.Time, error) {
	systemTime, errTime := m.getSystemTime()
	if errTime != nil {
		return time.Time{}, errTime
	}
	systemDate, errDate := m.getSystemDate()
	if errDate != nil {
		return time.Time{}, errDate
	}

	then, _ := time.Parse("02012006150405", string(systemDate[3:11])+string(systemTime[3:9]))
	return then, nil
}

func (m *MundiClient) getSystemTime() ([]byte, error) {
	return m.sendAndReceiveMessage([]byte{getSystemTime, emptyLength})
}

func (m *MundiClient) getSystemDate() ([]byte, error) {
	return m.sendAndReceiveMessage([]byte{getSystemDate, emptyLength})
}
