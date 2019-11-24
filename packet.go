package gozk

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// packet error definition
var (
	ErrInvalidPacketLength = errors.New("Packet length less than 16")
	ErrBadStartTag         = errors.New("Bad start tag")
	ErrInvalidChecksum     = errors.New("Checksum not valid")
)

// Packet represent object which being transferred
// between client and server (zk device)
type Packet struct {
	command uint16
	data    []byte

	size uint16

	session      uint16
	reply        uint16
	replyCounter uint16
}

// Marshal encode packet to byte array data
func (p Packet) Marshal() []byte {
	if p.data == nil {
		p.data = make([]byte, 0)
	}

	zkp := make([]byte, 16+len(p.data))
	copy(zkp[:4], StartTag)

	binary.LittleEndian.PutUint16(zkp[6:8], 0)
	binary.LittleEndian.PutUint16(zkp[8:10], p.command)
	binary.LittleEndian.PutUint16(zkp[12:14], p.session)
	binary.LittleEndian.PutUint16(zkp[14:16], p.reply)

	copy(zkp[16:], p.data)

	// calculate payload length
	binary.LittleEndian.PutUint16(zkp[4:6], uint16(len(zkp)-8))

	// calculate checksum
	binary.LittleEndian.PutUint16(zkp[10:12], checksum(zkp[8:]))

	return zkp[:]
}

// Unmarshal decode byte array data to zk packet
func (p *Packet) Unmarshal(zkp []byte) error {
	// minimum valid response packet
	if len(zkp) < 16 {
		return ErrInvalidPacketLength
	}

	// check start tag
	if bytes.Compare(zkp[:4], StartTag) != 0 {
		return ErrBadStartTag
	}

	// extract packet size
	p.size = binary.LittleEndian.Uint16(zkp[4:8])

	// validate checksum
	if !isPayloadValid(zkp[8:]) {
		return ErrInvalidChecksum
	}

	p.reply = binary.LittleEndian.Uint16(zkp[8:10])
	p.session = binary.LittleEndian.Uint16(zkp[12:14])
	p.replyCounter = binary.LittleEndian.Uint16(zkp[14:16])
	p.data = zkp[16:]

	return nil
}

// OK check if reply is valid
func (p Packet) OK() bool {
	Println("reply", p.reply, CmdAckOk, p.reply == CmdAckOk)
	return p.reply == CmdAckOk
}

// CreateCommandPacket construct request packet
func CreateCommandPacket(cmd uint16, data []byte, sessionID, replyNo uint16) Packet {
	return Packet{
		command: cmd,
		data:    data,
		session: sessionID,
		reply:   replyNo,
	}
}
