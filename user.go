package gozk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

type User struct {
	UserID     string
	UserSN     uint16
	Name       string
	Password   string
	CardNo     int
	AdminLevel uint16
	Enabled    bool
	GroupNo    int
	Timezones  []string
}

type UserQuery struct {
	t *Terminal

	// local store
	users map[uint16]User
}

// readAllUsers fetch all user id into internal memory
func (q *UserQuery) readAllUsers() error {
	cmdData, err := hex.DecodeString("0109000500000000000000")
	if err != nil {
		return err
	}

	if err := q.t.SendCommand(CmdDataWrrq, cmdData); err != nil {
		return err
	}

	var reply Packet
	if err := q.t.ReceiveLongReply(&reply, 1024); err != nil {
		return err
	}

	// re-init user dictionary
	q.users = make(map[uint16]User)

	dataset := reply.data
	size := len(dataset)

	if size < 5 {
		return errors.New("User dataset length is invalid")
	}

	// skip first 4 bytes
	i := 4

	// user entry is 72 bytes long
	for i+72 < size {
		var user User

		// extract serial number
		user.UserSN = binary.LittleEndian.Uint16(dataset[i : i+2])

		// extract permission token
		// permToken := dataset[i+2]

		// extract password, if it invalid use ""
		if dataset[i+3] != 0x00 {
			pass := dataset[i+3 : i+11]

			//remove trailing zeros
			user.Password = string(bytes.Trim(pass, "\x00"))
		}

		user.Name = string(bytes.Trim(dataset[i+11:i+35], "\x00"))
		user.CardNo = int(binary.LittleEndian.Uint64(dataset[i+35 : i+39]))
		user.GroupNo = int(dataset[i+39])

		// TODO: extract timezone

		user.UserID = string(bytes.Trim(dataset[i:48:i+57], "\x00"))

		q.users[user.UserSN] = user

		i += 72
	}

	return nil
}

func NewUserQuery(t *Terminal) *UserQuery {
	return &UserQuery{t: t}
}
