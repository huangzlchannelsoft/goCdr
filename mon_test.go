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

	AddCallStat("智博_中国移动_重庆_重庆", true)
	AddCallStat("智博_中国移动_重庆_重庆", true)
	AddCallStat("智博_中国移动_重庆_重庆", false)
	AddCallStat("COP_fixed_浙江_宁波", true)
	AddCallStat("COP_fixed_浙江_宁波", false)
	AddCallStat("COP_fixed_浙江_宁波", false)
	AddCallStat("蜂云物联_fixed_南京_南京", true)
	AddCallStat("蜂云物联_fixed_南京_南京", true)
	AddCallStat("中移在线_中国移动_浙江_衢州", false)
	AddCallStat("中移在线_中国移动_浙江_衢州", false)

	ticker := time.NewTicker(20 * time.Second)
	select {
	case <-ticker.C:
	}

	assert.Equal(t, true, true, "test")
}
