package main

import (
	"context"
	"flag"
	"log"

	"github.com/galihrivanto/go-zk/push"
)

type hook struct{}

func (h hook) OnInitialExchange(d push.Device) *push.ExchangeCommand {
	log.Println("SN:", d.SN)
	log.Println("options:", d.Option)
	log.Println("lang:", d.Language)

	return nil
}

func main() {
	var (
		host     string
		addr     string
		certFile string
		keyFile  string
	)

	flag.StringVar(&host, "host", "192.168.1.177", "host name")
	flag.StringVar(&addr, "address", ":8081", "http server address")
	flag.StringVar(&certFile, "cert-file", "cert.pem", "TLS cert file")
	flag.StringVar(&keyFile, "key-file", "key.pem", "TLS key file")

	flag.Parse()

	option := &push.ServerOption{
		Name:     host,
		Address:  addr,
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	s := push.NewServer(option, &hook{})
	s.Start(context.Background())
}
