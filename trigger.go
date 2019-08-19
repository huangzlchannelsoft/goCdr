// trigger
package main

import (
	"context"
	"log"
	"time"
)

const (
	TRIGGER_NEW_A_DAY = 0
	TRIGGER_BYE_BYE   = 1
	TRIGGER_BYE_OK    = 2
)

var (
	triggerChanList []*chan int
	offsetHour_     int
)

func init() {
	log.Println("init trigger!")
}

func curCnYearDay() int {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Println("[Err] curCnYearDay", err.Error())
		return time.Now().YearDay()
	}
	return time.Now().Add(time.Duration(offsetHour_) * time.Hour).In(loc).YearDay()
}

func TrickerDeamon(ctx context.Context, offsetHour int) {
	log.Println("enter TrickerDeamon.")
	offsetHour_ = offsetHour

	defer func() {
		log.Println("exit TrickerDeamon.")
	}()

	ticker := time.NewTicker(time.Minute)
	dtime := curCnYearDay() //time.Now().YearDay()

	for {
		select {
		case <-ctx.Done():
			for _, tc := range triggerChanList {
				*tc <- TRIGGER_BYE_BYE
			}
			return
		case <-ticker.C:
			if dtime != curCnYearDay() {
				dtime = curCnYearDay()
				for _, tc := range triggerChanList {
					*tc <- TRIGGER_NEW_A_DAY
				}
			}
		}
	}
}

func ApplyTrigger(id string) *chan int {
	tc := make(chan int, 1)
	triggerChanList = append(triggerChanList, &tc)
	return &tc
}

func TriggerExited() {
	for _, tc := range triggerChanList {
		<-*tc
		close(*tc)
	}
	triggerChanList = nil
}
