package gozk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"strings"
	"time"
)

// terminal related error
var (
	ErrNotConnected       = errors.New("Terminal not connected. Please call Connect()")
	ErrConnectionFailed   = errors.New("Failed to connect to remote terminal")
	ErrNullPacketReceiver = errors.New("Packet recevier should not nil")
)

// Terminal represent server / zk device
type Terminal struct {
	address string
	timeout time.Duration

	// state
	conn net.Conn

	lastPacket *Packet
}

func (t *Terminal) sendPacket(cmd uint16, data []byte) error {
	if t.conn == nil {
		return ErrNotConnected
	}

	var session, replyCounter uint16

	if t.lastPacket != nil {
		session = t.lastPacket.session
		replyCounter = t.lastPacket.replyCounter
	}

	Println("session", session)
	Println("reply", replyCounter)

	// increment reply id
	replyCounter++
	if replyCounter >= math.MaxUint16 {
		replyCounter -= math.MaxUint16
	}

	p := CreateCommandPacket(cmd, data, session, replyCounter)

	Printf("send: %x\n", p.Marshal())

	// send
	_, err := t.conn.Write(p.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (t *Terminal) receivePacket(bufSize int, vars ...int) ([]byte, error) {
	if t.conn == nil {
		return nil, ErrNotConnected
	}

	if len(vars) > 0 {
		deadline := time.Now().Add(time.Second * time.Duration(vars[0]))
		t.conn.SetReadDeadline(deadline)
	} else {
		t.conn.SetReadDeadline(time.Time{})
	}

	if bufSize == 0 {
		bufSize = 1024
	}

	buf := make([]byte, bufSize)
	_, err := t.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	Printf("received: %x\n", bytes.Trim(buf, "\x00"))

	return buf, nil
}

// SendCommand send packet command to undelying connection
func (t *Terminal) SendCommand(cmd uint16, data []byte) error {
	return t.sendPacket(cmd, data)
}

// ReceiveReply wait reply from server (zk device) and decode into
// given packet parameter
func (t *Terminal) ReceiveReply(reply *Packet, vars ...int) error {
	if reply == nil {
		reply = new(Packet)
	}

	var bufSize = 1024
	if len(vars) > 0 && vars[0] > 0 {
		bufSize = vars[0]
	}

	b, err := t.receivePacket(bufSize)
	if err != nil {
		return err
	}

	// decode
	if err := reply.Unmarshal(b); err != nil {
		return err
	}

	// store last received packet
	t.lastPacket = reply

	return nil
}

// ReceiveLongReply wait large dataset reply from server (zk device) and decode into
// given packet parameter
func (t *Terminal) ReceiveLongReply(reply *Packet, vars ...int) error {
	if reply == nil {
		reply = new(Packet)
	}

	var bufSize = 1024
	if len(vars) > 0 && vars[0] > 0 {
		bufSize = vars[0]
	}

	b, err := t.receivePacket(bufSize)
	if err != nil {
		return err
	}

	// decode initial packet
	if err := reply.Unmarshal(b); err != nil {
		return err
	}

	t.lastPacket = reply

	// prepare packet
	var packet []byte

	if reply.reply == CmdData {
		// device sent the dataset immediately, i.e. short dataset
		return nil
	} else if reply.reply == CmdPrepareData {
		// receives first part of the packet with the long dataset
		b, err = t.receivePacket(16)
		if err != nil {
			return err
		}

		packet = append(packet, b...)

		// extract size of total packet
		size := 8 + int(binary.LittleEndian.Uint16(b[4:6]))
		remaining := size - len(packet)

		// keep reading until received completed dataset
		for remaining > 0 {
			b, err = t.receivePacket(bufSize)
			if err != nil {
				return err
			}

			packet = append(packet, b...)
			remaining = size - len(packet)
		}

		// decode complete packet
		if err := reply.Unmarshal(packet); err != nil {
			return err
		}

		return nil
	} else if reply.reply == CmdAckOk {
		// device sent the dataset with additional commands, i.e. longer
		// dataset, see ex_data spec
		size := int(binary.LittleEndian.Uint64(b[1:5]))

		// create data for "ready for data" command
		ready := make([]byte, 4)
		binary.LittleEndian.PutUint64(ready, uint64(size))

		if err := t.SendCommand(CmdDataRdy, ready); err != nil {
			return err
		}

		// receives the prepare data reply
		b, err = t.receivePacket(24)
		if err != nil {
			return err
		}

		// receives the first part of the packet with the long dataset
		b, err = t.receivePacket(16)
		if err != nil {
			return err
		}

		packet = append(packet, b...)

		// extract size of total packet
		size = 8 + int(binary.LittleEndian.Uint16(b[4:6]))
		remaining := size - len(packet)

		// keep reading until received completed dataset
		for remaining > 0 {
			b, err = t.receivePacket(bufSize)
			if err != nil {
				return err
			}

			packet = append(packet, b...)
			remaining = size - len(packet)
		}

		// decode complete packet
		if err := reply.Unmarshal(packet); err != nil {
			return err
		}

		// receives the acknowledge after the dataset packet
		if _, err := t.receivePacket(bufSize); err != nil {
			return err
		}

		t.lastPacket = reply

		// send free data command
		if err := t.SendCommand(CmdFreeData, nil); err != nil {
			return err
		}

		// receives ack
		if _, err := t.receivePacket(bufSize); err != nil {
			return err
		}

		reply.reply++

		return nil

	}

	return nil
}

// SendAndReceive is convenience wrapper around send command
// and receive reply
func (t *Terminal) SendAndReceive(cmd uint16, data []byte, reply *Packet, vars ...int) error {
	if err := t.SendCommand(cmd, data); err != nil {
		return err
	}

	return t.ReceiveReply(reply, vars...)
}

// Connect establish connection to target terminal
// by send connect command and wait it reply
func (t *Terminal) Connect() error {
	if t.conn != nil {
		return nil
	}

	conn, err := net.Dial("tcp", t.address)
	if err != nil {
		return err
	}

	t.conn = conn
	if err := t.conn.SetReadDeadline(time.Now().Add(t.timeout)); err != nil {
		return err
	}

	// send connect command and wait for reply
	var response Packet
	if err := t.SendAndReceive(CmdConnect, nil, &response); err != nil {
		return err
	}

	// set SDKBuild variable of the device
	if err := t.SetInfo("SDKBuild", "1"); err != nil {
		return err
	}

	return nil
}

// Disconnect close connection from remote terminal
func (t *Terminal) Disconnect() error {
	if t.conn != nil {
		defer t.conn.Close()

		// send connect command and wait for reply
		var response Packet
		if err := t.SendAndReceive(CmdExit, nil, &response); err != nil {
			return err
		}

		// ensure ack OK
		if !response.OK() {
			return errors.New("Close error")
		}
	}

	return nil
}

// Enable set device state to enable
func (t *Terminal) Enable() error {
	var response Packet
	if err := t.SendAndReceive(CmdEnabledevice, nil, &response); err != nil {
		return err
	}

	if !response.OK() {
		return errors.New("Enable device failed")
	}

	return nil
}

// Disable set device state to disable until given timeout
func (t *Terminal) Disable(vars ...time.Duration) error {
	var data []byte = nil
	if len(vars) > 0 && vars[0] > 0 {
		data = make([]byte, 2)
		binary.LittleEndian.PutUint16(data, uint16(vars[0].Seconds()))
	}

	var response Packet
	if err := t.SendAndReceive(CmdDisabledevice, data, &response); err != nil {
		return err
	}

	if !response.OK() {
		return errors.New("Disable device failed")
	}

	return nil
}

// GetTime return decoded time of the device
func (t *Terminal) GetTime() time.Time {
	var response Packet
	if err := t.SendAndReceive(CmdGetTime, nil, &response); err != nil {
		return time.Time{}
	}

	return decodeTime(response.data)
}

// SetTime set time of device
func (t *Terminal) SetTime(datetime time.Time) error {
	if datetime.IsZero() {
		datetime = time.Now()
	}

	var response Packet
	if err := t.SendAndReceive(CmdSetTime, encodeTime(datetime), &response); err != nil {
		return err
	}

	if !response.OK() {
		return errors.New("Failed to set device time")
	}

	return nil
}

// GetVersion inquiry device version
func (t *Terminal) GetVersion() string {
	var response Packet
	if err := t.SendAndReceive(CmdGetVersion, nil, &response); err != nil {
		Println(err)
		return ""
	}

	return string(response.data)
}

// GetInfo inquiry device info for given key
func (t *Terminal) GetInfo(key string) string {
	var response Packet
	if err := t.SendAndReceive(CmdOptionsRrq, []byte(key), &response); err != nil {
		Println(err)
		return ""
	}

	parts := strings.SplitN(string(response.data), "=", 2)
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

// SetInfo set device info for given key
func (t *Terminal) SetInfo(key, value string) error {
	var response Packet
	if err := t.SendAndReceive(CmdOptionsWrq, []byte(fmt.Sprintf("%s=%s\x00", key, value)), &response); err != nil {
		return err
	}

	if !response.OK() {
		return errors.New("Set device info failed")
	}

	if err := t.SendAndReceive(CmdRefreshoption, nil, &response); err != nil {
		return err
	}

	// NOTE: for now, ignore 2nd ack
	// if !response.OK() {
	// 	return errors.New("Set device info failed")
	// }

	return nil
}

// NewTerminal initiate new remote terminal
func NewTerminal(address string, opt ...time.Duration) *Terminal {
	timeout := time.Second * 5
	if len(opt) > 0 && opt[0] > 0 {
		timeout = opt[0]
	}

	return &Terminal{
		address: address,
		timeout: timeout,
	}
}
