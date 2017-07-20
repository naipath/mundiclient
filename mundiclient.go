package mundiclient

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	startOfText       = 0x02
	endOfTransmission = 0x04
	getVersionID      = 0x56
	emptyLength       = 0x00
)

type MundiClient struct {
	conn net.Conn
}

func New(ip string, port int) *MundiClient {
	return &MundiClient{createConnection(ip, port)}
}

func (m MundiClient) Close() {
	m.conn.Close()
}

func (m MundiClient) GetVersion() string {
	lsb, msb := calculateChecksum(getVersionID, emptyLength)
	message := []byte{startOfText, getVersionID, emptyLength, lsb, msb, endOfTransmission}

	response := m.sendAndReceive(message)

	return string(response[3:9])
}

type Counters struct {
	Lifetime uint32
	Recent   uint32
}

func (m MundiClient) GetCounters() Counters {

	lsb, msb := calculateChecksum(0x47, emptyLength)
	message := []byte{startOfText, 0x47, emptyLength, lsb, msb, endOfTransmission}
	response := m.sendAndReceive(message)

	lifetime := binary.BigEndian.Uint32(response[3:7])
	recent := binary.BigEndian.Uint32(response[7:11])

	return Counters{lifetime, recent}
}

func (m MundiClient) GetStatusMessage() string {
	lsb, msb := calculateChecksum(0x58, emptyLength)
	response := m.sendAndReceive([]byte{startOfText, 0x58, emptyLength, lsb, msb, endOfTransmission})
	length := response[2]

	statusMessage := string(response[3 : length+3])

	return statusMessage
}

type StatusData struct {
	Version                            uint16
	PrintStatus                        bool
	FailureStatus                      bool
	WarningStatus                      bool
	MaintenanceStatus                  bool
	DetectorFault                      bool
	MarkNotLoaded                      bool
	CodeParametersNotLoaded            bool
	TestingLaser                       bool
	DisableShutter                     bool
	ExternalUpdateSingleShotNotUpdated bool
	ComPortDisconnected                bool
	BarcodeError                       bool
	EStop                              bool
	ExternalInterlocks                 bool
	CoolantTemperature                 bool
	DC24Volts                          bool
	ShutterKeyswitch                   bool
	Keyswitch                          bool
	DC48Volts                          bool
	CoolantFlow                        bool
	EmissionIndicator                  bool
	ShutterPrint                       bool
	ShutterStandby                     bool
	VSWRError                          bool
	OverModulation                     bool
	TachoFault                         bool
	RFStatus                           bool
	SystemEnable                       bool
	VapourExtractorFault               bool
	GalvoPower                         bool
	GalvoTemperature                   bool
	GalvoCableDisconnected             bool
}

func (m MundiClient) GetStatusData() StatusData {
	lsb, msb := calculateChecksum(0x57, emptyLength)
	response := m.sendAndReceive([]byte{startOfText, 0x57, emptyLength, lsb, msb, endOfTransmission})

	return StatusData{
		binary.BigEndian.Uint16(response[3:5]), //Version
		response[5]&128 != 0,                   //PrintStatus
		response[5]&64 != 0,                    //FailureStatus
		response[5]&32 != 0,                    //WarningStatus
		response[5]&16 != 0,                    //MaintenanceStatus
		response[6]&128 != 0,                   //DetectorFault
		response[6]&64 != 0,                    //MarkNotLoaded
		response[6]&32 != 0,                    //CodeParametersNotLoaded
		response[6]&16 != 0,                    //TestingLaser
		response[6]&8 != 0,                     //DisableShutter
		response[6]&4 != 0,                     //ExternalUpdateSingleShotNotUpdated
		response[6]&2 != 0,                     //ComPortDisconnected
		response[6]&1 != 0,                     //BarcodeError
		response[7]&128 != 0,                   //EStop
		response[7]&64 != 0,                    //ExternalInterlocks
		response[7]&32 != 0,                    //CoolantTemperature
		response[7]&16 != 0,                    //DC24Volts
		response[7]&8 != 0,                     //ShutterKeyswitch
		response[7]&4 != 0,                     //Keyswitch
		response[7]&2 != 0,                     //DC48Volts
		response[8]&128 != 0,                   //CoolantFlow
		response[8]&64 != 0,                    //EmissionIndicator
		response[8]&32 != 0,                    //ShutterPrint
		response[8]&16 != 0,                    //ShutterStandby
		response[8]&8 != 0,                     //VSWRError
		response[8]&4 != 0,                     //OverModulation
		response[8]&2 != 0,                     //TachoFault
		response[8]&1 != 0,                     //RFStatus
		response[9]&128 != 0,                   //SystemEnable
		response[9]&64 != 0,                    //VapourExtractorFault
		response[9]&32 != 0,                    //GalvoPower
		response[9]&16 != 0,                    //GalvoTemperature
		response[9]&8 != 0,                     //GalvoCableDisconnected
	}
}

func (m MundiClient) sendAndReceive(message []byte) []byte {
	m.conn.Write(message)
	result, err := bufio.NewReader(m.conn).ReadBytes(endOfTransmission)
	fmt.Printf("Got the following response:\n%08b\n", result)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}

func createConnection(ip string, port int) net.Conn {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))

	if err != nil {
		fmt.Print("error opening connection\n", err)
		os.Exit(1)
	}
	return conn
}

func calculateChecksum(input ...byte) (byte, byte) {
	var total uint16
	for _, value := range input {
		total += uint16(value)
	}
	lsb := byte(total & 0xFF)
	msb := byte((total >> 8) & 0xFF)
	return lsb, msb
}
