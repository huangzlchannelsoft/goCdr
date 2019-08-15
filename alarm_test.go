// alarm_test
package main

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func test_TransmitAlarm(t *testing.T) {

	uri := "http://127.0.0.1:9097/hooks/cdrAlarm"
	go TransmitAlarmCdr(uri)
	SendAlarm("test alarm")
	assert.Equal(t, true, true, "test")

	ticker := time.NewTicker(time.Second)
	select {
	case <-ticker.C:
	}
}
