package push

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type CommandCallback func(CommandResponse)

type Command struct {
	// command id which correspondence to operation id
	// see const commands
	ID string

	// command id
	CMD string

	// command payload / body
	Payload []byte

	Callback CommandCallback
}

func (c *Command) Marshal() ([]byte, error) {
	buf := acquireBuffer()
	defer releaseBuffer(buf)

	// command id
	buf.WriteString(fmt.Sprintf("C:%s", c.ID))

	// command
	buf.WriteString(fmt.Sprintf(":%s", c.CMD))

	if c.Payload != nil {
		buf.Write(sp)
		buf.Write(c.Payload)
	}

	return buf.Bytes(), nil
}

type CommandResponse struct {
	ID     string
	Return int
	CMD    string

	Payload []byte
}

// IsOK check whether response is valid
// see const response
func (c CommandResponse) IsOK() bool {
	return c.Return == 0
}

// Unmarshall implement PayloadDecoder interface
func (c *CommandResponse) Unmarshall(b []byte) error {
	c.ID = extractValue(b, "ID").ToString()
	c.Return = extractValue(b, "Return").ToInt()
	c.CMD = extractValue(b, "CMD").ToString()

	// TODO: get payload

	return nil
}

func generateUUID() uuid.UUID {
	guid, err := uuid.NewV4()
	if err != nil {
		return uuid.FromBytesOrNil(nil)
	}
	return guid
}

func randomCommandID() string {
	return generateUUID().String()
}
