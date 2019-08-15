// mon
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	CallStatDB   = "callstat.db"
	CallStatBulk = "phone2CallStat"
	CallStatJob  = "cdrCallStat"
)

type CallRateStat struct {
	productor         string
	isp               string
	province          string
	area              string
	callCounter       int
	callErrCounter    int
	oldCallCounter    int
	oldCallErrCounter int
}

var (
	property2Stat       map[string]*CallRateStat
	callRateStatList    []*CallRateStat
	oldCallRateStatList []*CallRateStat
	cstLock             sync.Mutex
)

func init() {
	log.Println("init callStat!")
	property2Stat = make(map[string]*CallRateStat)
}

func TransmitMonCdr() {

	setKeyValue := func(k []byte, v []byte) {
		log.Println(string(k))
		pp := Key2PhoneProperty(string(k))

		vs := strings.Split(string(v), "_")
		callCounter, _ := strconv.Atoi(vs[0])
		callErrCounter, _ := strconv.Atoi(vs[1])
		oldCallCounter, _ := strconv.Atoi(vs[2])
		oldCallErrCounter, _ := strconv.Atoi(vs[3])

		crs := &CallRateStat{
			productor:         pp.productor,
			isp:               pp.isp,
			province:          pp.province,
			area:              pp.area,
			callCounter:       callCounter,
			callErrCounter:    callErrCounter,
			oldCallCounter:    oldCallCounter,
			oldCallErrCounter: oldCallErrCounter,
		}

		cstLock.Lock()
		defer cstLock.Unlock()
		property2Stat[string(k)] = crs
		callRateStatList = append(callRateStatList, crs)
	}
	boltEnumKeyValue(CallStatDB, CallStatBulk, setKeyValue)

	cursor := 0
	getKeyValue := func() ([]byte, []byte) {
		if len(callRateStatList) <= cursor {
			crs := callRateStatList[cursor]
			cursor++

			k := PhoneProperty2Key(&PhoneProperty{crs.productor, crs.isp, crs.province, crs.area})
			v := fmt.Sprintf("%d_%d_%d_%d", crs.callCounter, crs.callErrCounter, crs.oldCallCounter, crs.oldCallErrCounter)
			return []byte(k), []byte(v)
		}
		return nil, nil
	}

	trigger := GetTrigger()

	for {
		select {
		case x := <-*trigger:
			cstLock.Lock()
			defer cstLock.Unlock()

			if x == TRIGGER_EXIT_PROC { //save cur status
				cursor = 0
				boltBatchWriteKeyValue(CallStatDB, CallStatBulk, getKeyValue)
				return
			} else if x == TRIGGER_NEW_A_DAY { //reset status
				oldCallRateStatList = callRateStatList
				callRateStatList = nil
				property2Stat = make(map[string]*CallRateStat)
				boltDeleteBucket(CallStatDB, CallStatBulk)
			}
		}
	}
}

/*
 */
func AddCallStat(pp *PhoneProperty, ok bool) {

	crs := property2Stat[PhoneProperty2Key(pp)]
	if crs == nil {

		crs = &CallRateStat{
			productor:         pp.productor,
			isp:               pp.isp,
			province:          pp.province,
			area:              pp.area,
			callCounter:       0,
			callErrCounter:    0,
			oldCallCounter:    0,
			oldCallErrCounter: 0,
		}

		cstLock.Lock()
		defer cstLock.Unlock()
		property2Stat[PhoneProperty2Key(pp)] = crs
		callRateStatList = append(callRateStatList, crs)
	}

	crs.callCounter++
	if !ok {
		crs.callErrCounter++
	}
}

/**mon plugin
 */
func InitCallStat(args ...string) {

}

func SetCallStatMeta(newMetrix NewMetrix, regMetrix RegMetrix) {
	cdrCallStatJob := CallStatJob
	cdrCallStatPrex := "cdr"
	cdrCallStatMetrix := []struct {
		name   string
		typ    int
		labels []string
	}{
		{"calls", METRIX_GAUGEVEC, []string{"productor", "isp", "province", "area"}},

		{"errCalls", METRIX_GAUGEVEC, []string{"productor", "isp", "province", "area"}},
	}

	for i := 0; i < len(cdrCallStatMetrix); i++ {
		//job, metrix type, prex, metrix, labels
		newMetrix(cdrCallStatJob, cdrCallStatMetrix[i].typ, cdrCallStatPrex, cdrCallStatMetrix[i].name, cdrCallStatMetrix[i].labels)
	}
	regMetrix(cdrCallStatJob)
}

func CallCallStat(setMetrix SetMetrix) {
	cstLock.Lock()
	defer cstLock.Unlock()

	cdrCallStatJob := CallStatJob

	crsList := callRateStatList
	if oldCallRateStatList != nil {
		crsList = oldCallRateStatList
		oldCallRateStatList = nil
	}
	for _, crt := range crsList {
		if crt.callCounter != crt.oldCallCounter {
			setMetrix(float64(crt.callCounter), []string{crt.productor, crt.isp, crt.province, crt.area}, cdrCallStatJob, 0)
			setMetrix(float64(crt.callErrCounter), []string{crt.productor, crt.isp, crt.province, crt.area}, cdrCallStatJob, 1)
			crt.oldCallCounter = crt.callCounter
			crt.oldCallErrCounter = crt.callErrCounter
		}
	}
}
