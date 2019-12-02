package push

// ExchangeCommand represent response from
// server to device upon initial exchange
type ExchangeCommand struct {
	SN string

	AttLogStamp   int
	OperLogStamp  int
	AttPhotoStamp int

	ErrorDelay    int
	Delay         int
	TransTimes    string
	TransInterval int
	TransFlag     string
	TimeZone      int
	Realtime      int
	Encrypt       int
	ServerVer     string
}

// Marshall implement payload.Marshall interface
func (c ExchangeCommand) Marshall() ([]byte, error) {
	buf := acquireBuffer()
	defer releaseBuffer(buf)

	buf.WriteString("GET OPTION FROM: ")
	buf.WriteString(c.SN)

	writeIntValue(buf, "ATTLOGStamp", lf, c.AttLogStamp, true)
	writeIntValue(buf, "OPERLOGStamp", lf, c.OperLogStamp, true)
	writeIntValue(buf, "ATTPHOTOStamp", lf, c.AttPhotoStamp, true)

	writeIntValue(buf, "ErrorDelay", lf, c.ErrorDelay)
	writeIntValue(buf, "Delay", lf, c.Delay)
	writeStringValue(buf, "TransTimes", lf, c.TransTimes)
	writeIntValue(buf, "TransInterval", lf, c.TransInterval)
	writeStringValue(buf, "TransFlag", lf, c.TransFlag)
	writeIntValue(buf, "TimeZone", lf, c.TimeZone)
	writeIntValue(buf, "Realtime", lf, c.Realtime)
	writeStringValue(buf, "ServerVer", lf, c.ServerVer)
	writeIntValue(buf, "Encrypt", lf, c.Encrypt, true)

	return buf.Bytes(), nil
}

// Unmarshall implement payload.Unmarshall interface
func (c *ExchangeCommand) Unmarshall(b []byte) error {
	if b == nil {
		return ErrEmptyPayload
	}

	// extract values
	c.SN = extractValue(b, "SN", []byte(":")).ToString()

	c.AttLogStamp = extractValue(b, "ATTLOGStamp").ToInt()
	c.OperLogStamp = extractValue(b, "OPERLOGStamp").ToInt()
	c.AttPhotoStamp = extractValue(b, "ATTPHOTOStamp").ToInt()

	c.ErrorDelay = extractValue(b, "ErrorDelay").ToInt()
	c.Delay = extractValue(b, "Delay").ToInt()
	c.TransTimes = extractValue(b, "TransTimes").ToString()
	c.TransInterval = extractValue(b, "TransInterval").ToInt()
	c.TransFlag = extractValue(b, "TransFlag").ToString()
	c.Realtime = extractValue(b, "Realtime").ToInt()
	c.ServerVer = extractValue(b, "ServerVer").ToString()

	return nil
}
