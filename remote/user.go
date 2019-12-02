package remote

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

// define user related error list
var (
	ErrInvalidUserID = errors.New("UserID not found on user list")
)

const maxSN = 10000

// FpData contains information  if user
// fingerprint template
type FpData struct {
	Index    int
	Template []byte
	Flag     int
}

// User represent registered user on zk device
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

	// fingerprint data
	fingerPrints map[int]FpData
}

// SetFpTemplate set user's fingerprint template
func (u User) SetFpTemplate(index int, template []byte, flag int) {
	if u.fingerPrints == nil {
		u.fingerPrints = make(map[int]FpData)
	}

	u.fingerPrints[index] = FpData{
		Index:    index,
		Template: template,
		Flag:     flag,
	}
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

// UserQuery provides access to registered users on device
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

// readAllFingerprintTemplates request all fingerprint templates
func (q *UserQuery) readAllFingerprintTemplates() error {
	cmdData, err := hex.DecodeString("0107000200000000000000")
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

	dataset := reply.data
	size := len(dataset)

	if size < 5 {
		return errors.New("User fingerprint dataset length is invalid")
	}

	// skip first 4 bytes
	i := 4

	// every template entry is 6 bytes + template length
	var (
		templateSize int
		userSN       int
		fpIndex      int
		fpFlag       int
		fpTemplate   []byte
	)

	for i < size {
		// extract template size
		templateSize = int(binary.LittleEndian.Uint16(dataset[i:i+2])) - 6

		// extract user serial no
		userSN = int(binary.LittleEndian.Uint16(dataset[i+2 : i+4]))

		// extract fingerprint index
		fpIndex = int(dataset[i+4])

		// extract fingerprint flag
		fpFlag = int(dataset[i+5])

		// extract template
		copy(fpTemplate, dataset[i+6:i+templateSize+6])

		// check if user sn is exists
		if _, ok := q.users[userSN]; ok {
			q.users[userSN].SetFpTemplate(fpIndex, fpTemplate, fpFlag)
		}

		i += templateSize + 6
	}

	return nil
}

// FindByName return users which match name keyword
func (q *UserQuery) FindByName(keyword string, callback func([]User)) error {
	// ensure user has been loaded
	if q.users == nil {
		if err := q.readAllUsers(); err != nil {
			return err
		}
	}

	var found []User
	for _, v := range q.users {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(keyword)) {
			found = append(found, v)
		}
	}

	callback(found)

	return nil
}

// FindAll return all users
func (q *UserQuery) FindAll(callback func([]User)) error {
	return q.FindByName("", callback)
}

// getUserSN obtains user internal index from user id
func (q *UserQuery) getUserSN(userID string) int {
	if q.users != nil {
		for sn, u := range q.users {
			if strings.EqualFold(u.UserID, userID) {
				return sn
			}
		}
	}

	return -1
}

// GetVerificatinMode get user verification mode
func (q *UserQuery) GetVerificatinMode(userID string) (VerificationKind, error) {
	sn := q.getUserSN(userID)
	if sn == -1 {
		return 0, ErrInvalidUserID
	}

	var response Packet
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(sn))

	if err := q.t.SendAndReceive(CmdVerifyRrq, data, &response); err != nil {
		return 0, err
	}

	if response.data == nil || len(response.data) < 3 {
		return 0, errors.New("Invalid response")
	}

	return VerificationKind(response.data[2]), nil
}

// SetVerificationMode override user verification mode
func (q *UserQuery) SetVerificationMode(userID string, mode VerificationKind) error {
	sn := q.getUserSN(userID)
	if sn == -1 {
		return ErrInvalidUserID
	}

	var response Packet
	data := make([]byte, 24)
	binary.LittleEndian.PutUint16(data[:2], uint16(sn))
	data[2] = byte(mode)

	if err := q.t.SendAndReceive(CmdVerifyWrq, data, &response); err != nil {
		return err
	}

	if !response.OK() {
		return errors.New("Set verification mode failed")
	}

	return nil
}

// getNextSN return next available number
// due possibility user can be freely deleted
// SN number may not in increment order
// therefore we need to find most smaller
// "free" number available
func (q *UserQuery) getNextSN() int {
	for i := 0; i < maxSN; i++ {
		// if sn not used
		if _, ok := q.users[i]; !ok {
			return i
		}
	}

	return 0
}

// CreateUser register new user on device
func (q *UserQuery) CreateUser() error {
	return nil
}

// DeleteUser remover user record from device
func (q *UserQuery) DeleteUser(userID string) error {
	sn := q.getUserSN(userID)
	if sn == -1 {
		return ErrInvalidUserID
	}

	var response Packet
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(sn))

	if err := q.t.SendAndReceive(CmdDeleteUser, data, &response); err != nil {
		return err
	}

	if !response.OK() {
		return errors.New("Delete user failed")
	}

	return nil
}

// DownloadFingerPrint request a fingerprint template for given user
func (q *UserQuery) DownloadFingerPrint(userID string, index int) error {
	sn := q.getUserSN(userID)
	if sn == -1 {
		return ErrInvalidUserID
	}

	data := make([]byte, 3)
	binary.LittleEndian.PutUint16(data[:2], uint16(sn))

	if err := q.t.SendCommand(CmdUsertempRrq, data); err != nil {
		return err
	}

	var response Packet
	if err := q.t.ReceiveLongReply(&response, 1024); err != nil {
		return err
	}

	// put on user data
	// q.users[sn].SetFpTemplate(index)

	return nil
}

// DeleteFingerPrint remove registered fingerprint on user id
func (q *UserQuery) DeleteFingerPrint(userID string, index int) error {
	sn := q.getUserSN(userID)
	if sn == -1 {
		return ErrInvalidUserID
	}

	_, ok := q.users[sn]
	if !ok {
		return ErrInvalidUserID
	}

	// remove on remote termina
	data := make([]byte, 25)
	copy(data[0:len(userID)], []byte(userID))
	data[24] = byte(index)

	var response Packet
	if err := q.t.SendAndReceive(CmdDelFptmp, data, &response); err != nil {
		return err
	}

	if !response.OK() {
		return errors.New("Failed to delete fingerprint")
	}

	// delete local data
	delete(q.users[sn].fingerPrints, index)

	return nil
}

// NewUserQuery initiate user query
func NewUserQuery(t *Terminal) *UserQuery {
	return &UserQuery{t: t}
}
