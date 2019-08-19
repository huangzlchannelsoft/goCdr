// parser_test
package main

import (
	"testing"

	"github.com/bmizerany/assert"
)

func Test_CheckCalledNumber(t *testing.T) {
	assert.Equal(t, checkCalledNumber("01088822200"), true, "01088822200")
	assert.Equal(t, checkCalledNumber("010-88822200"), true, "010-88822200")
	assert.Equal(t, checkCalledNumber("057427614632"), true, "057427614632")
	assert.Equal(t, checkCalledNumber("0574-27614632"), true, "0574-27614632")
	assert.Equal(t, checkCalledNumber("13587010754"), true, "13587010754")

	assert.Equal(t, checkCalledNumber("750140112"), false, "750140112")
	assert.Equal(t, checkCalledNumber("804003801613"), false, "804003801613")
}

func Test_ContinueFail(t *testing.T) {
	gCfg.ConAbnormal = 20
	LoadExcelPhoneProductor("NumberShow-201908.xlsx")
	LoadTxtPhoneIsp("phone_area_operators.txt")
	cdrs := []string{
		"000010000,g001,2001,17082794034,02910987654,,,20190816173550,20190816173550,20190816173550,93e28f41-b3c6-4c20-944b-df243f4d48dc,ent001,999111001,TKYD-SIP,",
		"000010001,g001,2002,17082794034,13408290987,,,20190816173550,20190816173550,20190816173550,2c681bbf-27d2-4302-9ede-567b760cf5c8,ent001,999111001,TKYD-SIP,",
		"000010002,g001,2003,17082794034,13435414321,,,20190816173550,20190816173550,20190816173550,78236356-7409-4c06-b08e-bcc4e851163a,ent001,999111001,TKYD-SIP,",
		"000010003,g001,2004,17082794034,085176543210,,,20190816173550,20190816173550,20190816173550,d4fe77ef-7b5a-457b-8605-c0be321c2b07,ent001,999111001,TKYD-SIP,",
		"000010004,g001,2005,17082794034,13411952109,,,20190816173550,20190816173550,20190816173550,ee7edf7f-888e-4eab-ae47-9b2d94fab44e,ent001,999111001,TKYD-SIP,",
		"000010005,g001,2006,17082794034,075454321098,,,20190816173550,20190816173550,20190816173550,5df40bb5-8e36-4af2-8734-7666c15d923f,ent001,999111001,TKYD-SIP,",
		"000010006,g001,2007,17082794034,13450704321,,,20190816173550,20190816173550,20190816173550,b0408ba6-5b34-4707-bafe-ea11bbd7b3b1,ent001,999111001,TKYD-SIP,",
		"000010007,g001,2008,17082794034,02076543210,,,20190816173550,20190816173550,20190816173550,4a9cc06f-131b-4479-80d1-69ef3ceb7f1c,ent001,999111001,TKYD-SIP,",
		"000010008,g001,2009,17082794034,059176543210,,,20190816173550,20190816173550,20190816173550,2b25a3f1-5fc5-4bc0-ab3e-156f2776d013,ent001,999111001,TKYD-SIP,",
		"000010009,g001,2010,17082794034,059254321098,,,20190816173550,20190816173550,20190816173550,34176d21-dfad-4010-b1fb-4e87dea43201,ent001,999111001,TKYD-SIP,",
		"000010010,g001,2011,17082794034,085154321098,,,20190816173550,20190816173550,20190816173550,c0c36037-346a-40d6-a3b7-2344fbe6ee2b,ent001,999111001,TKYD-SIP,",
		"000010011,g001,2012,17082794034,13453008765,,,20190816173550,20190816173550,20190816173550,0844588a-e229-4b81-a422-67cdcfcf2592,ent001,999111001,TKYD-SIP,",
		"000010012,g001,2013,17082794034,13419478765,,,20190816173550,20190816173550,20190816173550,f25879ce-b119-4eb8-8237-9c4b41f50544,ent001,999111001,TKYD-SIP,",
		"000010013,g001,2014,17082794034,13433908765,,,20190816173550,20190816173550,20190816173550,65107fcf-ecf3-4e3d-a2ed-dc3c43932e12,ent001,999111001,TKYD-SIP,",
		"000010014,g001,2015,17082794034,091676543210,,,20190816173550,20190816173550,20190816173550,d1540030-0a26-4c3d-91ce-8f92800ea818,ent001,999111001,TKYD-SIP,",
		"000010015,g001,2016,17082794034,13422546543,,,20190816173550,20190816173550,20190816173550,d0bf04af-37ab-4290-b932-2ff346eb368d,ent001,999111001,TKYD-SIP,",
		"000010016,g001,2017,17082794034,13424052109,,,20190816173550,20190816173550,20190816173550,2abb2e39-8df0-4cff-8b62-52f160ad6a2f,ent001,999111001,TKYD-SIP,",
		"000010017,g001,2018,17082794034,031554321098,,,20190816173550,20190816173550,20190816173550,a2528b5b-0582-418a-9ca2-1d5ac89968f5,ent001,999111001,TKYD-SIP,",
		"000010018,g001,2019,17082794034,13415738765,,,20190816173550,20190816173550,20190816173550,b2eb4b12-780c-4de3-a9ea-b13e584d5c82,ent001,999111001,TKYD-SIP,",
		"000010019,g001,2020,17082794034,082754321098,,,20190816173550,20190816173550,20190816173550,7155bb66-2a6c-47b7-88f5-71729b057239,ent001,999111001,TKYD-SIP,",
		"000010020,g001,2021,17082794034,13450254321,,,20190816173550,20190816173550,20190816173550,3993d016-9a41-486e-94c7-7c91ca8dcd31,ent001,999111001,TKYD-SIP,",
		"000010021,g001,2022,17082794034,047510987654,,,20190816173550,20190816173550,20190816173550,d407b715-8521-4550-be66-abdee666f3d2,ent001,999111001,TKYD-SIP,",
		"000010022,g001,2023,17082794034,097098765432,,,20190816173550,20190816173550,20190816173550,3ac1b916-fb5f-451d-971d-fcb74766ecf0,ent001,999111001,TKYD-SIP,",
		"000010023,g001,2024,17082794034,13466154321,,,20190816173550,20190816173550,20190816173550,708a4b20-21ac-48fc-9b4a-562934bd8cc5,ent001,999111001,TKYD-SIP,",
		"000010024,g001,2025,17082794034,13426338765,,,20190816173550,20190816173550,20190816173550,7080ae40-fce6-4ea7-9894-d36b69383caf,ent001,999111001,TKYD-SIP,",
		"000010025,g001,2026,17082794034,02854321098,,,20190816173550,20190816173550,20190816173550,2af01b45-97d0-4b4a-a6b4-ffd3ba40e1e5,ent001,999111001,TKYD-SIP,",
		"000010026,g001,2027,17082794034,091354321098,,,20190816173550,20190816173551,20190816173650,b9191691-4a84-4572-8389-65f2b69013ce,ent001,999111001,TKYD-SIP,",
		"000010027,g001,2028,17082794034,13453820987,,,20190816173550,20190816173551,20190816173650,32b42b38-3252-4289-98b4-297a258d2fcc,ent001,999111001,TKYD-SIP,",
		"000010028,g001,2029,17082794034,13436926543,,,20190816173550,20190816173551,20190816173650,5513cb09-8cee-487e-a2a8-2bb83a1566c7,ent001,999111001,TKYD-SIP,",
		"000010029,g001,2030,17082794034,13413994321,,,20190816173550,20190816173551,20190816173650,94e9dd78-2239-4f45-99ef-d39312b69631,ent001,999111001,TKYD-SIP,",
	}
	idx := 0
	recvCdr := func() string {
		defer func() {
			idx++
		}()
		if idx == len(cdrs) {
			assert.Equal(t, true, false, "end test")
		}
		return cdrs[idx]
	}
	sendAlarm := func(data string) {
		t.Log(data)
		assert.Equal(t, "", data, "")
	}
	ParseCdr(recvCdr, sendAlarm)
}

func Test_TimePolicy(t *testing.T) {

}
