package remote

import (
	"bytes"
	"context"
	"errors"
	"net"
)

type Event struct {
	Type uint16
	Data []byte
}

type EventListener struct {
	t *Terminal
}

func (e *EventListener) enableRealtime() error {
	// ensure device is enabled
	if err := e.t.Enable(); err != nil {
		return err
	}

	var response Packet
	if err := e.t.SendAndReceive(CmdRegEvent, []byte{0xff, 0xff, 0x00, 0x00}, &response); err != nil {
		return err
	}

	if !response.OK() {
		return errors.New("Failed to register event")
	}

	Printf("%x", bytes.Trim(response.Payload(), "\x00"))

	return nil
}

func (e *EventListener) Listen(ctx context.Context) (<-chan Event, error) {
	if err := e.enableRealtime(); err != nil {
		return nil, err
	}

	ch := make(chan Event)
	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				Println("context cancelled")
				return
			default:
				// get raw packet
				b, err := e.t.receivePacket(4096)
				if err != nil {
					// if error due read timeout, continue listening
					if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
						break
					}

					Println("error while listening event", err)
					return
				}

				var response Packet
				if err := response.Unmarshal(b); err != nil {
					Println(err)
					break
				}

				// ack packet
				if err := e.t.SendCommand(CmdAckOk, nil); err != nil {
					Println("error when acknowledge event")
					break
				}

				ch <- Event{
					Type: response.session,
					Data: response.data,
				}
			}
		}

	}()

	return ch, nil
}

func NewEventListener(t *Terminal) *EventListener {
	return &EventListener{t: t}
}
