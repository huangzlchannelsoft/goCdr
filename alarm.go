// alarm
package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	alarmCh chan string = make(chan string, 1000)
)

type AlarmSend func(alarm string)

func init() {
	log.Println("init alarm!")
}

func AlarmCdr(uri string) {

	client := http.Client{Timeout: time.Second * 5}

	httpPost := func(uri string, data string) bool {
		request, err := http.NewRequest("POST", uri, strings.NewReader(data))
		if err != nil {
			log.Println("[Err] new request.", err.Error())
			return false
		}
		request.Header.Set("Content-type", "application/json")
		request.Header.Set("charset", "utf-8")

		response, err := client.Do(request)
		if err != nil {
			log.Println("[Err] Do request.", err.Error())
			return false
		}
		if response.StatusCode != 200 {
			log.Println("[Err]", "reponse err.", response.StatusCode)
			return false
		}

		return true
	}

	for {
		alarm := <-alarmCh

		ok := httpPost(uri, alarm)
		if ok { //TODO: send alarm error.

		}
	}
}

/*
 */
func SendAlarm(alarm string) {
	alarmCh <- alarm
}
