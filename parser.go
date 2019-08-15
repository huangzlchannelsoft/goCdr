// parser
package main

import (
	"log"
	"strings"
	"regexp"
)

func init() {
	log.Println("init parser!")
}


func parsingCdr (cdr string) map[string] string {
	log.Println("[INFO] 接收到cdr信息为：",cdr)
	resultMap := make(map[string] string,5)

	columArrs := strings.Split(cdr,",")
	if len(columArrs) > 0 {
		// 主叫号码
		resultMap["callNumber"] = columArrs[3]
		// 被叫号码
		resultMap["calledNumber"] = columArrs[4]
		// 创建时间
		resultMap["createTime"] = columArrs[7]
		// 接通时间
		resultMap["turnOnTime"] = columArrs[8]
		// 挂断时间
		resultMap["shutDownTime"] = columArrs[9]
		// sessionId
		resultMap["sessionId"] = columArrs[10]
		return resultMap
	}
	return nil
}

// 校验被叫号码
func checkCalledNumber (calledNumber string) bool {
	// 校验标示
	flag := false

	// 手机号固话校验正则表达式
	mobileRegular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	fixedLineRegular := "^0\\d{2,3}-?\\d{7,8}$"

	isornoMobile,_ := regexp.MatchString(mobileRegular, calledNumber)
	isornoFixed,_ := regexp.MatchString(fixedLineRegular, calledNumber)

	// 判断是否合法的手机号或固话
	if isornoMobile || isornoFixed {
		flag = true
	}
	return flag
}

func getInfo (callNumber string) string {
	log.Println("[INFO]","主叫号码：",callNumber)

	info := string("联通|上海|上海")
	return info
}

//bitmap: math/big kv/bolt
func ParseCdr(recvCdr CdrRecv, sendAlarm AlarmSend) {
	log.Println("run parseCdr")

	for {
		cdr := recvCdr()
		log.Println("[INFO] 接收到cdr数据信息为：",cdr)
		// 解析cdr

		resultMap := parsingCdr(cdr)

		if resultMap != nil {
			// 被叫号码
			calledNumber := resultMap["calledNumber"]
			// 判断被叫号码是否标准 验证该条话单为外线
			isOutside := checkCalledNumber(calledNumber)
			if isOutside {
				log.Println("[INFO]","此条cdr为外线话单...")

				callNumber := resultMap["callNumber"]
				info := getInfo(callNumber)
				log.Println("info = ",info)


				createTime := resultMap["createTime"]
				turnOnTime := resultMap["turnOnTime"]
				shutDownTime := resultMap["shutDownTime"]


				// 默认正常的bit位
				vBit := uint8(0)
				// 判断话务信息是否异常
				if (strings.Compare(createTime,turnOnTime) == 0) && (strings.Compare(createTime,shutDownTime) == 0) {
					// 异常数据
					vBit = 1
				}

				/*log.Println("vBit = ",vBit)
				log.Println(" len(bitMap[info]) == ",len(bitMap[info]) == 0)
				if len(bitMap) == 0 || len(bitMap[info].data) == 0 {
					vBitMap := new(Bitmap)
					vBitMap.SetBit(0,vBit)
					bitMap[info] = *vBitMap
				} else {
					log.Println("vbit = ",vBit)
					vBitMap := bitMap[info]
					vBitMap.SetBit(vBitMap.bitsize,vBit)
					bitMap[info] = vBitMap
				}*/

			}
		}

		//parse cdr
		//stat by aera||producter
		alarm := ""
		sendAlarm(alarm)
	}
}
