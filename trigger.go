// trigger
package main

import (
	"context"
	"log"
	"time"
)

const (
	TRIGGER_NEW_A_DAY = 0
	TRIGGER_EXIT_PROC = 1
)

var (
	triggerChanList []*chan int
)

func init() {
	log.Println("init trigger!")
}

func TrickerDeamon(ctx context.Context) {

	ticker := time.NewTicker(time.Minute)
	dtime := time.Now().YearDay()

	for {
		select {
		case <-ctx.Done():
			for _, tc := range triggerChanList {
				*tc <- TRIGGER_EXIT_PROC
			}
			return
		case <-ticker.C:
			if dtime != time.Now().YearDay() {
				dtime = time.Now().YearDay()
				for _, tc := range triggerChanList {
					*tc <- TRIGGER_NEW_A_DAY
				}
			}
		}
	}
}

func GetTrigger() *chan int {
	tc := make(chan int, 1)
	triggerChanList = append(triggerChanList, &tc)
	return &tc
}
