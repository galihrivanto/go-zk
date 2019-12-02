package remote

import "log"

var verbose bool

// SetVerbose enable verbose logging
// useful for debugging
func SetVerbose() {
	verbose = true
}

// Print send output to standard logger
func Print(v ...interface{}) {
	if !verbose {
		return
	}

	log.Print(v...)
}

// Println send output to standard logger
func Println(v ...interface{}) {
	if !verbose {
		return
	}

	log.Println(v...)
}

// Printf send formatted output to standard logger
func Printf(format string, v ...interface{}) {
	if !verbose {
		return
	}

	log.Printf(format, v...)
}
