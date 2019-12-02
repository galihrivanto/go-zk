package push

import (
	"bytes"
	"fmt"
)

func writeValue(b *bytes.Buffer, key string, sep []byte, v []byte) {
	if v == nil || string(v) == "" {
		return
	}

	if sep != nil {
		b.Write(sep)
	}

	b.WriteString(key + "=")
	b.Write(v)
}

func writeStringValue(b *bytes.Buffer, key string, sep []byte, v string) {
	if v == "" {
		return
	}

	if sep != nil {
		b.Write(sep)
	}

	b.WriteString(key + "=")
	b.WriteString(v)
}

func writeIntValue(b *bytes.Buffer, key string, sep []byte, v int, vars ...bool) {
	var useNone bool
	if len(vars) > 0 {
		useNone = vars[0]
	}

	var val string

	if v == 0 && useNone {
		val = "None"
	} else {
		val = fmt.Sprintf("%d", v)
	}

	if sep != nil {
		b.Write(sep)
	}

	b.WriteString(key + "=")
	b.WriteString(val)
}

func extractValue(b []byte, key string, vars ...[]byte) value {
	separator := keyValueSeparator
	if len(vars) > 0 {
		separator = vars[0]
	}

	// scan index key value
	idx := bytes.Index(b, []byte(key))
	if idx == -1 {
		return nil
	}

	// trim from start index
	b = b[idx:]

	// trim by next key value separator if exists
	if sep := bytes.Index(b, separator); sep != -1 {
		b = b[:sep]
	}

	// value start after value separator
	if sep := bytes.Index(b, keySeparator); sep != -1 {
		return b[sep+1:]
	}

	// consider invalid
	return nil
}
