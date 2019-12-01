package push

import (
	"bytes"
	"testing"
)

const sampleKeyValues = `SN=123456789&pushver=1.2.3&language=72&pushcommkey=987654321`

func TestExtractValue(t *testing.T) {
	testCases := [][]string{
		[]string{"SN", "123456789"},
		[]string{"pushver", "1.2.3"},
		[]string{"language", "72"},
		[]string{"pushcommkey", "987654321"},
	}

	for _, tc := range testCases {
		v := extractValue([]byte(sampleKeyValues), tc[0])
		if v == nil {
			t.Errorf("result shouldn't nil")
			t.FailNow()
		}

		if bytes.Compare(v, []byte(tc[1])) == -1 {
			t.Errorf("expected %s but returned %s", tc[1], string(v))
		}
	}
}
