package main

import (
	"flag"
	"log"
	"os"

	rgozk "github.com/galihrivanto/go-zk/remote"
)

func main() {
	var host string

	flag.StringVar(&host, "host", "192.168.1.201:4370", "address of zk device")
	flag.Parse()

	term := rgozk.NewTerminal(host)
	if err := term.Connect(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer term.Disconnect()

	u := rgozk.NewUserQuery(term)
	err := u.FindAll(func(users []rgozk.User) {
		for _, u := range users {
			log.Println("ID:", u.UserID)
			log.Println("SN:", u.UserSN)
			log.Println("Name:", u.Name)
			log.Println("Card:", u.CardNo)
			log.Println("=============")
		}
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
