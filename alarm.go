// alarm
package main

import (
	"log"
	"net/http"
	"time"
)

var (
	alarmCh chan string = make(chan string, 1000)
)

type AlarmSend func(alarm string)

func init() {
	log.Println("init alarm!")
}

func TransmitAlarmCdr(uri string) {

	client := http.Client{Timeout: time.Second * 5}

	for {
		alarm := <-alarmCh

		ok := HttpPost(&client, uri, alarm, nil)
		if ok { //TODO: send alarm error.
			log.Printf("send alarm: %s.\n", alarm)
		} else {
			log.Printf("[Err] send alarm: %s.\n", alarm)
		}
	}
}

/*
 */
func SendAlarm(alarm string) {
	alarmCh <- alarm
}
