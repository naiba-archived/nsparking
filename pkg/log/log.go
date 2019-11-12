package log

import (
	"log"
)

func init() {
	log.SetPrefix("[NsParking] ")
}

// Println ...
func Println(v ...interface{}) {
	log.Println(v...)
}
