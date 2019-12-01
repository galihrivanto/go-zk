package main

import (
	"context"
	"flag"

	"github.com/galihrivanto/go-zk/push"
)

type hook struct{}

func (h hook) OnInitialExchange(d push.Device) bool {
	return true
}

func main() {
	var addr string

	flag.StringVar(&addr, "address", ":8080", "http server address")

	flag.Parse()

	option := &push.ServerOption{
		Name:    "localhost",
		Address: addr,
	}
	s := push.NewServer(option, &hook{})
	s.Start(context.Background())
}
