// phone_test
package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"testing"

	"github.com/bmizerany/assert"
)

func test_PhonelistFile(t *testing.T) {
	phoneList, err := ioutil.ReadFile("phonelist.txt")
	if err != nil {
		t.Log(err.Error())
		return
	}

	for {
		n, phone, err := bufio.ScanWords(phoneList, true)
		if err != nil || n == 0 {
			break
		}
		log.Println(n, string(phone), err)
		phoneList = phoneList[n:]
	}

	assert.Equal(t, true, true, "test")
}

func test_PhoneExslFile(t *testing.T) {
	//LoadExcelPhoneProperty("OutNumber-201908.xlsx")
}

func test_GetPhoneProperty(t *testing.T) {
	ispUri := "http://paas.ccod.com/t/qn-api/phone_area/queryisp/"
	proUri := "http://paas.ccod.com/t/qn-api/phone_area/queryvoip/"
	SetPhonePropertyUri(ispUri, proUri)

	na := []string{"18466909107", "075536563445", "17092399436", "02062167359"}
	nb := []string{"13996370001", "057455833166", "02586306666", "13587010754"}
	keys := GetPhoneProperty(na, nb)

	log.Println(keys)
}
