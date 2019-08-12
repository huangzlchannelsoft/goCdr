// collector
package main

import (
	"log"
	"time"
)

func init() {
	log.Println("init collector!")
}

func CollectCdr(path string, transmit CdrSend) {
	log.Println("run collectCdr.")

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			//check new cdr
			//read new cdr
			cdr := ""
			transmit(cdr)
		}
	}
}
