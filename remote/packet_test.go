package remote

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

type testCase struct {
	command uint16
	data    []byte
	result  []byte
}

var testCases []testCase

func loadPacketTestCases() {
	testCases = make([]testCase, 1)

	// connect
	b, _ := hex.DecodeString("5050827d08000000e80317fc00000000")
	testCases[0] = testCase{
		command: CmdConnect,
		data:    nil,
		result:  b,
	}
}

func TestPacketMarshal(t *testing.T) {
	SetVerbose()
	loadPacketTestCases()

	for _, tc := range testCases {
		p := CreateCommandPacket(tc.command, tc.data, uint16(0), uint16(0))

		t.Log(fmt.Sprintf("%x vs %x\n", tc.result, p.Marshal()))

		if bytes.Compare(tc.result, p.Marshal()) != 0 {
			t.Error("unmarshal not match with result")
			t.FailNow()
		}
	}
}
