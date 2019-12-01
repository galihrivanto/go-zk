package gozk

import "errors"

import "encoding/binary"

import "fmt"

// VerificationKind supported by device
type VerificationKind int

// verification modes
const (
	VerifyWithPassword = iota
	VerifyWithFingerPrint
	VerifyWithCard
)

type EventAttLog struct {
	UID              string
	VerificationKind VerificationKind
	DateString       string
}

// EventAttLogFromEvent decode att log event from event payload
func EventAttLogFromEvent(evt Event) (EventAttLog, error) {
	if evt.Type != EfAttlog || len(evt.Data) < 32 {
		return EventAttLog{}, errors.New("invalid att log event")
	}

	uid := string(evt.Data[0:9])
	verificationKind := VerificationKind(binary.LittleEndian.Uint16(evt.Data[24:26]))
	dateString := fmt.Sprintf("20%d/%d/%d %02d:%02d:%02d",
		evt.Data[26],
		evt.Data[27],
		evt.Data[28],
		evt.Data[29],
		evt.Data[30],
		evt.Data[31],
	)

	return EventAttLog{
		UID:              uid,
		VerificationKind: verificationKind,
		DateString:       dateString,
	}, nil

}
