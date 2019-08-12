// parser
package main

import (
	"log"
)

func init() {
	log.Println("init parser!")
}

//bitmap: math/big kv/bolt
func ParseCdr(recvCdr CdrRecv, sendAlarm AlarmSend) {
	log.Println("run parseCdr")

	for {
		recvCdr()
		//parse cdr
		alarm := ""
		sendAlarm(alarm)
	}
}
