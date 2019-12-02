package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	rgozk "github.com/galihrivanto/go-zk/remote"
)

func main() {
	rgozk.SetVerbose()

	var host string

	flag.StringVar(&host, "host", "192.168.1.201:4370", "address of zk device")
	flag.Parse()

	term := rgozk.NewTerminal(host)
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

	listener := rgozk.NewEventListener(term)
	events, err := listener.Listen(ctx)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("listening events...")
	go func() {
		for event := range events {
			switch event.Type {
			case rgozk.EfAttlog:
				att, err := rgozk.EventAttLogFromEvent(event)
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
