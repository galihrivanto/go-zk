package push

import (
	"errors"
	"strconv"
)

// lib error
var (
	ErrEmptyPayload    = errors.New("Payload is empty")
	ErrPayloadIsNil    = errors.New("Payload value is nil")
	ErrReceiverInvalid = errors.New("Receiver doesn't implement correct payload interface")
)

// static value
var (
	keyValueSeparator = []byte("&")
	keySeparator      = []byte("=")
	sp                = []byte(" ")
	lf                = []byte("\n")
	ht                = []byte("\t")
)

// Payload define wellknown operation which should be
// supported by communication payload, either from
// request body or request URL
type Payload interface {
	Marshall() ([]byte, error)
	Unmarshall([]byte) error
}

// Marshall encode object into bytes
func Marshall(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, ErrPayloadIsNil
	}

	// target value should implement payload interface
	o, ok := v.(Payload)
	if !ok {
		return nil, ErrReceiverInvalid
	}

	return o.Marshall()
}

// Unmarshall decode raw bytes to given object
func Unmarshall(b []byte, v interface{}) error {
	if v == nil {
		return ErrPayloadIsNil
	}

	// target value should implement payload interface
	o, ok := v.(Payload)
	if !ok {
		return ErrReceiverInvalid
	}

	return o.Unmarshall(b)
}

// value of command or result
type value []byte

func (v value) ToString(vars ...string) string {
	if v == nil {
		if len(vars) > 0 {
			return vars[0]
		}
	}

	return string(v)
}

func (v value) ToInt(vars ...int) int {
	var def int
	if len(vars) > 0 {
		def = vars[0]
	}

	if v == nil {
		return def
	}

	if v, err := strconv.ParseInt(string(v), 10, 64); err != nil {
		return int(v)
	}

	return def
}
