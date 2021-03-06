package mundiclient

const (
	getFieldContents = 0xC3
	allFieldsNumber  = 0x00

	Unknown            = FieldType(0x00)
	Text               = FieldType(0x01)
	Date               = FieldType(0x02)
	Time               = FieldType(0x03)
	Operator           = FieldType(0x04)
	ExternalText       = FieldType(0x05)
	Incremental        = FieldType(0x06)
	BarcodeText        = FieldType(0x07)
	BarcodeIncremental = FieldType(0x08)
	BarcodeExternal    = FieldType(0x09)
	LogoVLM            = FieldType(0x0A)
	LogoBitmap         = FieldType(0x0B)
	LogoDXF            = FieldType(0x0C)
)

type FieldType byte

type Field struct {
	FieldNumber byte
	TypeOfField FieldType
	FieldText   string
}

func (m *MundiClient) GetAllFieldContents() ([]Field, error) {
	response, err := m.sendAndReceiveMessage(fieldRequest(allFieldsNumber))

	if err != nil {
		return nil, err
	}

	fields := []Field{buildField(response)}
	for i := 2; i <= int(response[3]); i++ {
		response, err = m.sendAndReceiveMessage(fieldRequest(byte(i)))
		if err != nil {
			return nil, err
		}

		fields = append(fields, buildField(response))
	}
	return fields, nil
}

func (m *MundiClient) GetFieldContents(fieldID byte) (Field, error) {
	response, err := m.sendAndReceiveMessage(fieldRequest(fieldID))
	if err != nil {
		return Field{}, err
	}
	return buildField(response), nil
}

func fieldRequest(fieldID byte) []byte {
	return []byte{getFieldContents, 0x01, fieldID}
}

func buildField(input []byte) Field {
	return Field{
		input[4],
		FieldType(input[5]),
		string(input[6: input[2]-3+6]),
	}
}
