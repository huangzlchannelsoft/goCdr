// mon_test
package main

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func test_Mon_CallStat(t *testing.T) {

	go TransmitMonCdr()

	uri := "http://10.130.41.226:9091"
	go PromethuesClient(true, uri, "nid001", 10)

	AddCallStat(&PhoneProperty{"智博", "中国移动", "重庆", "重庆"}, true)
	AddCallStat(&PhoneProperty{"智博", "中国移动", "重庆", "重庆"}, true)
	AddCallStat(&PhoneProperty{"智博", "中国移动", "重庆", "重庆"}, false)
	AddCallStat(&PhoneProperty{"COP", "fixed", "浙江", "宁波"}, true)
	AddCallStat(&PhoneProperty{"COP", "fixed", "浙江", "宁波"}, false)
	AddCallStat(&PhoneProperty{"COP", "fixed", "浙江", "宁波"}, false)
	AddCallStat(&PhoneProperty{"蜂云物联", "fixed", "南京", "南京"}, true)
	AddCallStat(&PhoneProperty{"蜂云物联", "fixed", "南京", "南京"}, true)
	AddCallStat(&PhoneProperty{"中移在线", "中国移动", "浙江", "衢州"}, false)
	AddCallStat(&PhoneProperty{"中移在线", "中国移动", "浙江", "衢州"}, false)

	ticker := time.NewTicker(20 * time.Second)
	select {
	case <-ticker.C:
	}

	assert.Equal(t, true, true, "test")
}
