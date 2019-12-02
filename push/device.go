package push

import (
	"fmt"
)

// Device represent ZK biometric pheripheral
type Device struct {
	// device serial no
	SN string

	// device option
	Option string

	// push service version
	PushVersion string

	// device language (see const languages)
	Language int

	// push comm key
	PushCommKey string
}

// Marshall implement payload.Marshall operation
func (d Device) Marshall() ([]byte, error) {
	buf := acquireBuffer()
	defer releaseBuffer(buf)

	// write SN
	buf.WriteString("SN=" + d.SN)

	buf.Write(keyValueSeparator)
	buf.WriteString("option=" + d.Option)

	buf.Write(keyValueSeparator)
	buf.WriteString("pushver=" + d.PushVersion)

	buf.Write(keyValueSeparator)
	buf.WriteString("language=" + fmt.Sprintf("%d", d.Language))

	buf.Write(keyValueSeparator)
	buf.WriteString("pushcommkey=" + d.PushCommKey)

	return buf.Bytes(), nil
}

// Unmarshall implement payload.Unmarshall interface
func (d *Device) Unmarshall(b []byte) error {
	if b == nil {
		return ErrEmptyPayload
	}

	// extract values
	d.SN = extractValue(b, "SN").ToString()
	d.Option = extractValue(b, "options").ToString()
	d.PushVersion = extractValue(b, "pushver").ToString()
	d.Language = extractValue(b, "language").ToInt()
	d.PushCommKey = extractValue(b, "pushcommkey").ToString()

	return nil
}
