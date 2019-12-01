package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	gozk "github.com/galihrivanto/go-zk"
)

func main() {
	gozk.SetVerbose()

	var host string

	flag.StringVar(&host, "host", "192.168.1.201:4370", "address of zk device")
	flag.Parse()

	term := gozk.NewTerminal(host)
	if err := term.Connect(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer term.Disconnect()

	log.Println("Version:", term.GetVersion())
	log.Println("Serial No:", term.GetInfo("~SerialNumber"))
	log.Println("Platform:", term.GetInfo("~Platform"))
	log.Println("DeviceName:", term.GetInfo("~DeviceName"))
	log.Println("Device Time:", term.GetTime())

	ctx, cancel := context.WithCancel(context.Background())

	listener := gozk.NewEventListener(term)
	events, err := listener.Listen(ctx)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("listening events...")
	go func() {
		for event := range events {
			switch event.Type {
			case gozk.EfAttlog:
				att, err := gozk.EventAttLogFromEvent(event)
				if err != nil {
					break
				}

				log.Printf("event: att log, uid: %s, kind: %d, date: %s", att.UID, att.VerificationKind, att.DateString)
			default:
				log.Printf("event: %v - %x\n", event.Type, bytes.Trim(event.Data, "\x00"))
			}

		}
	}()

	fmt.Scanln()
	fmt.Println("closing")
	cancel()
}
