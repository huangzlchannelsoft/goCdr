// phone_test
package main

import (
	"bufio"
	"fmt"
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

func test_LoadExcelPhoneProductor(t *testing.T) {
	LoadExcelPhoneProductor("NumberShow-201908.xlsx")
	assert.Equal(t, phone2Productor["18466909114"].productor, "智博", "")
	assert.Equal(t, phone2Productor["17188539704"].productor, "讯众", "")
	assert.Equal(t, phone2Productor["075536563441"].productor, "COP", "")
	assert.Equal(t, phone2Productor["17045320644"].productor, "蜂云物联", "")
	assert.Equal(t, phone2Productor["02260537322"].productor, "中移在线", "")
}

func test_LoadTxtPhoneIsp(t *testing.T) {
	LoadTxtPhoneIsp("phone_area_operators.txt")
	assert.Equal(t, phone2Isp["1340201"].isp, "中国移动", "")
	assert.Equal(t, phone2Isp["1668811"].isp, "中国联通", "")
	assert.Equal(t, phone2Isp["1859059"].isp, "中国联通", "")
	assert.Equal(t, phone2Isp["1859059"].province, "新疆", "")
	assert.Equal(t, phone2Isp["1859059"].area, "博乐", "")
	assert.Equal(t, phone2Isp["1734867"].isp, "中国电信", "")
	assert.Equal(t, phone2Isp["1734867"].province, "浙江", "")
	assert.Equal(t, phone2Isp["1734867"].area, "嘉兴", "")
	assert.Equal(t, phone2Isp["1559999"].isp, "中国联通", "")
	assert.Equal(t, phone2Isp["1559999"].province, "新疆", "")
	assert.Equal(t, phone2Isp["1559999"].area, "伊犁", "")
}

func test_GetPhoneProperty(t *testing.T) {
	LoadExcelPhoneProductor("NumberShow-201908.xlsx")
	LoadTxtPhoneIsp("phone_area_operators.txt")

	testCases := []struct {
		numberA  string
		numberB  string
		property PhoneProperty
	}{
		{
			"17077115433",
			"15507979869",
			PhoneProperty{"蜂云物联", "中国联通", "江西", "赣州"},
		},
		{
			"02022505577",
			"13165431234",
			PhoneProperty{"COP", "中国联通", "山东", "滨州"},
		},
		{
			"053266014281",
			"06338311777",
			PhoneProperty{"COP", "固话", "山东", "日照"},
		},
		{
			"01053189231",
			"13896313363",
			PhoneProperty{"讯众", "中国移动", "重庆", "万州"},
		},
		{
			"02022505554",
			"020-29899212",
			PhoneProperty{"COP", "固话", "广东", "广州"},
		},
		{
			"17077115463",
			"079182126888",
			PhoneProperty{"蜂云物联", "固话", "江西", "南昌"},
		},
	}

	for _, cas := range testCases {
		pp := GetPhoneProperty(cas.numberA, cas.numberB)
		log.Println(pp.productor, pp.isp, pp.province, pp.area)
		assert.Equal(t, pp.productor, cas.property.productor, fmt.Sprint(cas.numberA, cas.numberB))
		assert.Equal(t, pp.isp, cas.property.isp, fmt.Sprint(cas.numberA, cas.numberB))
		assert.Equal(t, pp.province, cas.property.province, fmt.Sprint(cas.numberA, cas.numberB))
		assert.Equal(t, pp.area, cas.property.area, fmt.Sprint(cas.numberA, cas.numberB))
	}
}

func test_UpdatePhoneProperty(t *testing.T) {
	ispUri := "http://paas.ccod.com/t/qn-api/phone_area/queryisp/"
	proUri := "http://paas.ccod.com/t/qn-api/phone_area/queryvoip/"
	SetPhonePropertyUri(ispUri, proUri)

	na := []string{"18466909107", "075536563445", "17092399436", "02062167359"}
	nb := []string{"13996370001", "057455833166", "02586306666", "13587010754"}
	UpdatePhoneProperty(na, nb)

	for i := 0; i < len(na); i++ {
		pp := GetPhoneProperty("18466909107", "13996370001")
		log.Println(pp.productor, pp.isp, pp.province, pp.area)
	}
}

func test_PhoneAlarm(t *testing.T) {

	sendAlarm := func(data string) {
		t.Log(data)
	}

	LoadExcelPhoneProductor("NumberShow-201908.xlsx")
	LoadTxtPhoneIsp("phone_area_operators.txt")
	SetPhoneAlarmFunc(sendAlarm)

	GetPhoneProperty("12345678901", "400888888888")
}
