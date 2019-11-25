package gozk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

type User struct {
	UserSN     int
	UserID     string
	Name       string
	Password   string
	CardNo     int
	NotEnabled int
	AdminLevel int
	Group      int
	Timezones  []int
}

// Marshal encode user info in byte array
func (u User) Marshal() []byte {
	dataset := make([]byte, 72)

	binary.LittleEndian.PutUint16(dataset[0:2], uint16(u.UserSN))
	dataset[2] = byte((u.AdminLevel << 1) | u.NotEnabled)
	copy(dataset[3:3+len(u.Password)], []byte(u.Password))
	copy(dataset[11:11+len(u.Name)], []byte(u.Name))
	binary.LittleEndian.PutUint64(dataset[35:39], uint64(u.CardNo))
	dataset[39] = byte(u.Group)

	if u.Timezones != nil && len(u.Timezones) > 0 {
		copy(dataset[40:42], []byte{0x01, 0x00})

		for i, tz := range u.Timezones {
			binary.LittleEndian.PutUint16(dataset[42+(i*2):44+(i*2)], uint16(tz))
		}
	}

	copy(dataset[48:48+len(u.UserID)], []byte(u.UserID))

	return dataset
}

// Unmarshal decode byte array into user info
func (u *User) Unmarshal(dataset []byte) error {
	// user entry must be 72 bytes long
	if len(dataset) != 72 {
		return errors.New("Invalid user entry length")
	}

	// extract serial number
	u.UserSN = int(binary.LittleEndian.Uint16(dataset[0:2]))

	// extract permission token
	permToken := dataset[2]
	u.AdminLevel = int(permToken) >> 1
	u.NotEnabled = int(permToken) & 1

	// extract password, if it invalid use ""
	if dataset[3] != 0x00 {
		// password is 8 bytes length
		// remove trailing zeros
		u.Password = string(bytes.Trim(dataset[3:11], "\x00"))
	}

	u.Name = string(bytes.Trim(dataset[11:35], "\x00"))
	u.CardNo = int(binary.LittleEndian.Uint16(dataset[35:39]))
	u.Group = int(dataset[39])

	if binary.LittleEndian.Uint16(dataset[40:42]) == 1 {
		u.Timezones = make([]int, 3)
		u.Timezones[0] = int(binary.LittleEndian.Uint16(dataset[42:44]))
		u.Timezones[1] = int(binary.LittleEndian.Uint16(dataset[44:46]))
		u.Timezones[2] = int(binary.LittleEndian.Uint16(dataset[46:48]))
	}

	u.UserID = string(bytes.Trim(dataset[48:57], "\x00"))

	return nil
}

type UserQuery struct {
	t *Terminal

	// local store
	users map[int]User
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
	q.users = make(map[int]User)

	dataset := reply.data
	size := len(dataset)

	if size < 5 {
		return errors.New("User dataset length is invalid")
	}

	// skip first 4 bytes
	i := 4

	// user entry is 72 bytes long
	var user User
	for i+72 < size {
		user = User{}

		if err := user.Unmarshal(dataset[i : i+72]); err != nil {
			return err
		}

		q.users[user.UserSN] = user

		i += 72
	}

	return nil
}

// FindByName return users which match name keyword
func (u *UserQuery) FindByName(keyword string, callback func([]User)) error {
	// ensure user has been loaded
	if u.users == nil {
		if err := u.readAllUsers(); err != nil {
			return err
		}
	}

	var found []User
	for _, v := range u.users {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(keyword)) {
			found = append(found, v)
		}
	}

	callback(found)

	return nil
}

// FindAll return all users
func (u *UserQuery) FindAll(callback func([]User)) error {
	return u.FindByName("", callback)
}

func NewUserQuery(t *Terminal) *UserQuery {
	return &UserQuery{t: t}
}
