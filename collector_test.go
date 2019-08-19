// collector_test
package main

import (
	"log"
	"testing"

	"github.com/bmizerany/assert"
)

func test_CheckCdr(t *testing.T) {
	assert.Equal(t, true, true, "test")
}

func Test_CollectCdr(t *testing.T) {
	gCfg.CdrFileBakPath = "D:/tmp/blackListErrData/bak"
	gCfg.CdrCommaTotal = 14
	send := func(cdr string) {
		log.Println(">", cdr)
	}
	CollectCdr("D:/tmp/blackListErrData/20190728", send)
}
