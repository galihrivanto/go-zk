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
	flag.Parse()

	term := gozk.NewTerminal(host)
	if err := term.Connect(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer term.Disconnect()

	u := gozk.NewUserQuery(term)
	err := u.FindAll(func(users []gozk.User) {
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
