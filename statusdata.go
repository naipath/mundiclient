package mundiclient

import "encoding/binary"

const (
	getStatusData    = 0x57
	getStatusMessage = 0x58
)

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

func (m MundiClient) GetStatusData() (StatusData, error) {
	response, err := m.sendAndReceiveMessage([]byte{getStatusData, emptyLength})
	if err != nil {
		return StatusData{}, err
	}
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
	}, nil
}

func (m MundiClient) GetStatusMessage() (string, error) {
	response, err := m.sendAndReceiveMessage([]byte{getStatusMessage, emptyLength})
	if err != nil {
		return "", err
	}
	statusMessageLength := response[2]
	return string(response[3 : statusMessageLength+3]), nil
}
