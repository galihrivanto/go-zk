package push

import "fmt"

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

	buf.WriteString("GET OPTION FROM:")
	buf.WriteString(c.SN)

	writeStringValue(buf, "ATTLOGStamp", lf, fmt.Sprintf("%d", c.AttLogStamp))
	writeStringValue(buf, "OPERLOGStamp", lf, fmt.Sprintf("%d", c.OperLogStamp))
	writeStringValue(buf, "ATTPHOTOStamp", lf, fmt.Sprintf("%d", c.AttLogStamp))

	writeStringValue(buf, "ErrorDelay", lf, fmt.Sprintf("%d", c.ErrorDelay))
	writeStringValue(buf, "Delay", lf, fmt.Sprintf("%d", c.Delay))
	writeStringValue(buf, "TransTimes", lf, c.TransTimes)
	writeStringValue(buf, "TransInterval", lf, fmt.Sprintf("%d", c.Delay))
	writeStringValue(buf, "TransFlag", lf, c.TransFlag)
	writeStringValue(buf, "Realtime", lf, fmt.Sprintf("%d", c.Realtime))
	writeStringValue(buf, "ServerVer", lf, c.ServerVer)

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
