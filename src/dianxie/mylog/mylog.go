package mylog

import (
	"log"
	"os"
)

const (
	// Bits or'ed together to control what's printed. There is no control over the
	// order they appear (the order listed here) or the format they present (as
	// described in the comments).  A colon appears after these items:
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota     // the date: 2009/01/23
	Ltime                         // the time: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

var logger *log.Logger

func StartLog() {
	file, err := os.Create("../log/error.log")
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	logger = log.New(file, "", log.LstdFlags)
}
func ErrorLog(text string) {

	//logger := log.New(file, "", log.LstdFlags|log.Llongfile)

	log.Println(text)
	logger.Println(text)

	//log.SetFlags(log.LstdFlags)
	//log.Fatal(text)
}
