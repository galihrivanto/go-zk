package gozk

import (
	"encoding/binary"
	"encoding/hex"
	"testing"
	"time"
)

func TestChecksum(t *testing.T) {
	SetVerbose()

	testCases := [][]string{
		{"0b00f38d03005a4b4661636556657273696f6e00", "5a17"},
		{"d007296af38d0a0009", "0000"},
	}

	var res string
	buf := make([]byte, 2)
	for _, tc := range testCases {
		b, err := hex.DecodeString(tc[0])
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		t.Logf(string(b))
		t.Logf("payload: %x\n", b)

		sum := checksum(b)

		t.Logf("checksum: %d\n", sum)

		binary.LittleEndian.PutUint16(buf, sum)

		t.Logf("sum: %x\n", buf)

		res = hex.EncodeToString(buf)

		if res != tc[1] {
			t.Errorf("expected %s but returned %s\n", tc[1], res)
			t.FailNow()
		}
	}
}

func TestEncodeDecodeTime(t *testing.T) {
	SetVerbose()

	var testCases = []time.Time{
		time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(2019, 12, 31, 12, 05, 0, 0, time.Local),
		time.Now(),
	}

	for _, tc := range testCases {
		encoded := encodeTime(tc)
		res := decodeTime(encoded)

		if !res.Equal(tc) {
			t.Errorf("expected %v but returned %v", tc, res)
			t.FailNow()
		}

	}
}
