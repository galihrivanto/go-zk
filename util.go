package gozk

import (
	"encoding/binary"
	"math"
	"time"
)

// Calculates checksum of packet.
//
// param: data to which the
// checksum is going to be applied.
// return: checksum result given as a number.
func checksum(payload []byte) uint16 {
	if len(payload) == 0 {
		return uint16(0)
	}

	// make odd length packet
	if len(payload)%2 == 1 {
		payload = append(payload, 0x00)
	}

	acc := int64(0)
	for len(payload) > 1 {
		acc += int64(binary.LittleEndian.Uint16(payload[0:2]))
		if acc > math.MaxUint16 {
			acc -= math.MaxUint16
		}

		payload = payload[2:]
	}

	acc = (acc & 0xFFFF) + ((acc & 0xFFFF0000) >> 16)

	for acc < 0 {
		acc += math.MaxUint16
	}

	return uint16(acc ^ 0xFFFF)
}

// decodeTime cecodes time, as given on ZKTeco get/set time commands.
// param: raw data with the time field stored in little endian.
// return: time.Time, with the extracted date.
func decodeTime(raw []byte) time.Time {
	// extract time value
	t := uint(binary.LittleEndian.Uint64(raw))

	Println("raw", t)

	second := int(t % 60)
	minute := int((t / 60) % 60)
	hour := int((t / 3600) % 24)
	day := int((t / (3600 * 24) % 31)) + 1
	month := int((t / (3600 * 24 * 31) % 12)) + 1
	year := int((t/(3600*24))/365) + 2000

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

func encodeTime(t time.Time) []byte {
	b := make([]byte, 8)

	v := ((t.Year()%100)*12*31+
		((int(t.Month())-1)*31)+
		t.Day()-1)*(24*60*60) +
		(t.Hour()*60+t.Minute())*60 +
		t.Second()

	Println("v", v)

	binary.LittleEndian.PutUint64(b, uint64(v))

	return b
}

// isPayloadValid checks if a given packet payload is valid, considering the checksum,
// where the payload is given with the checksum.

// param: byte array with the payload contents.
// return: if the payload is consistent, returns True,
// otherwise returns False.
func isPayloadValid(payload []byte) bool {
	// if the checksum is valid the checksum calculation, without removing the
	// checksum, should be equal to zero
	return checksum(payload) == 0
}
