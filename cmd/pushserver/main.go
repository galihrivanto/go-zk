package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/galihrivanto/go-zk/push"
)

type hook struct{}

func (h hook) OnInitialExchange(d push.Device) *push.ExchangeCommand {
	log.Println("SN:", d.SN)
	log.Println("options:", d.Option)
	log.Println("lang:", push.LangText(d.Language))
	log.Println("push ver:", d.PushVersion)

	return &push.ExchangeCommand{
		SN:          d.SN,
		AttLogStamp: 9999,
		Delay:       10,
		TimeZone:    7,
	}
}

func (h hook) Middlewares(option *push.ServerOption) []push.Middleware {
	return []push.Middleware{
		push.MiddlewareFunc(Verbose),
	}
}

func Verbose(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
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

	// explicitly register device
	s.RegisterDevice("BRM9181260009")

	go s.Start(context.Background())

	<-s.Ready()

	resp, err := s.Do("BRM9181260009", push.Command{CMD: "REBOOT"})
	if err != nil {
		log.Println(err)
	}

	if resp.IsOK() {
		log.Println("Send reboot command successful")
	}

	fmt.Scanln()

}
