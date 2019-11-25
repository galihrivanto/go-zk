package main

import (
	"flag"
	"log"
	"os"

	gozk "github.com/galihrivanto/go-zk"
)

func main() {
	var host string

	flag.StringVar(&host, "host", "192.168.1.201:4370", "address of zk device")

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
}
