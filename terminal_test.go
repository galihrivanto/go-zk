package gozk

import "testing"

func TestTerminal(t *testing.T) {
	SetVerbose()

	// test connect and disconnect
	address := getEnvar("ZK_ADDRESS", "192.168.1.201:4370")
	Println(address)
	term := NewTerminal(address)
	if err := term.Connect(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := term.Disable(); err != nil {
		t.Error(err)
	}

	Println("Version:", term.GetVersion())
	Println("Serial No:", term.GetInfo("~SerialNumber"))
	Println("Platform:", term.GetInfo("~Platform"))
	Println("DeviceName:", term.GetInfo("~DeviceName"))
	Println("Device Time:", term.GetTime())

	if err := term.Disconnect(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
