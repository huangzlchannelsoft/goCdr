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

func test_CollectCdr(t *testing.T) {
	send := func(cdr string) {
		log.Println(">", cdr)
	}
	CollectCdr("G:\\xbak\\cdr\\cdr", send)
}
