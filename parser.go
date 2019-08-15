// parser
package main

import (
	"log"
	"strings"
	"regexp"
	"math/big"
	"time"
	"fmt"
	"encoding/json"
)

func init() {
	log.Println("init parser!")
}

// 时间策略
const TIME_STRATEGY int8 = 0
// 异常条数策略
const ABNORMAL_STRATEGY int8 = 1


// cdr状态结构体
type CdrState struct {
	// true 正常 false 异常
	isNormal 		bool
	// 当前时间戳
	curTimestamp 	int64
}

type AlarmInfo struct {
	// 被叫号码所在key
	key 			string
	// 策略类型 0-时间策略 1->异常条数理策略
	strategyType 	int8
	// strategyType == 0 ? 56% : 20
	value 			string
}

var cdrStateMap map[string] []*CdrState

// 分割解析cdr数据
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

func getKeys (callNumber string,calledNumber string) []string {
	log.Println("[INFO]","主叫号码：",callNumber,"被叫号码：",calledNumber)

	keys := []string{"联通|上海|上海","联通|上海","联通"}

	return keys
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

// 连续异常条数状态策略
func abnormalStrategy (stateCdrs []*CdrState) *big.Int {
	log.Println("[WARN]","执行异常条数策略部分")

	// 异常数据发生条数
	abnormalCount := big.NewInt(0)
	for _, stateCdr := range stateCdrs {
		if !stateCdr.isNormal {
			abnormalCount = abnormalCount.Add(abnormalCount,big.NewInt(1))
		}
	}
	return abnormalCount
}

// 时间策略
func timeStrategy (stateCdrs []*CdrState) *big.Float {
	log.Println("[WARN]","执行时间策略部分")

	// 异常信息占比
	percentage := big.NewFloat(0.00)
	// 非正常条数
	abnormalCount := big.NewFloat(0.00)
	cdrsLen := big.NewFloat(float64(len(stateCdrs)))

	// 一定时间必须达到指定条数才会计算告警值
	if len(stateCdrs) >= 200 {
		// 遍历异常条数
		for _,stateCdr := range stateCdrs {
			if !stateCdr.isNormal {
				abnormalCount = abnormalCount.Add(abnormalCount,big.NewFloat(1))
			}
		}
		percentage = percentage.Quo(abnormalCount,cdrsLen)
		log.Println("percentage = ",percentage)
	}
	return percentage
}

//bitmap: math/big kv/bolt
func ParseCdr(recvCdr CdrRecv, sendAlarm AlarmSend) {
	log.Println("run parseCdr")

	for {
		cdr := recvCdr()
		log.Println("[INFO] 接收到cdr数据信息为：",cdr)
		// 解析后cdr映射信息
		resultMap := parsingCdr(cdr)

		if resultMap != nil {
			// 被叫号码
			calledNumber := resultMap["calledNumber"]
			// 判断被叫号码是否标准 验证该条话单为外线
			isOutside := checkCalledNumber(calledNumber)
			if isOutside {
				log.Println("[INFO]","此条cdr为外线话单...")

				callNumber := resultMap["callNumber"]
				createTime := resultMap["createTime"]
				turnOnTime := resultMap["turnOnTime"]
				shutDownTime := resultMap["shutDownTime"]

				// 默认正常的bit位
				isNormal := bool(true)
				// 判断话务信息是否正常
				if (strings.Compare(createTime,turnOnTime) == 0) && (strings.Compare(createTime,shutDownTime) == 0) {
					// 异常数据
					isNormal = false
				}
				log.Println("isNormal = ",isNormal)

				// 获取phoneKeys
				keys := getKeys(callNumber,calledNumber)
				log.Println("info = ",keys)

				// 当前时间戳
				curTimestamp := time.Now().Unix()
				log.Println("[INFO]","curTimestamp = ",curTimestamp)
				for _, key := range keys {
					log.Println("操作状态key值为：",key)

					// 发送至监控系统
					//sendCdrState(key,isNormal)

					// 告警百分比
					percentage := big.NewFloat(0.00)
					// 首次加入数据
					if cdrStateMap[key] == nil || (len(cdrStateMap[key]) == 0) {
						log.Println("[WARN]","当前key首次加入记录！")
						stateCdr := CdrState{isNormal,curTimestamp}
						stateCdrs := [] *CdrState{&stateCdr}
						cdrStateMap = make(map[string] []*CdrState)
						// 追加
						cdrStateMap[key] = stateCdrs
					} else {
						log.Println()
						stateCdr := CdrState{isNormal,curTimestamp}
						// 追加并取出之前已存在数据
						stateCdrs := append(cdrStateMap[key], &stateCdr)
						firstCdrCurTime := stateCdrs[0].curTimestamp

						// 是否达到一定时间策略
						if (curTimestamp - firstCdrCurTime) >= (5 * 60) {
							percentage = timeStrategy(stateCdrs)
							log.Println("[INFO]","时间策略计算返回值为：",percentage)
							// 是否大于告警阈值
							if big.NewFloat(0.5).Cmp(percentage) != 1 {
								finalPercentage := fmt.Sprint("%0.2f",percentage.Mul(percentage,big.NewFloat(100))) + "%"
								// 告警
								alarm := AlarmInfo{key,TIME_STRATEGY,finalPercentage}
								byteAlarm,err := json.Marshal(alarm)
								if err != nil {
									log.Println("[ERR]","err",err)
								}
								log.Println("byte",byteAlarm)
								strAlarm := string(byteAlarm)
								log.Println("[INFO]","strAlarm = ",strAlarm)
								//sendAlarm(strAlarm)
								// 数据清空
								cdrStateMap[key] = append(cdrStateMap[key][:0])
							}
						} else {
							// 此条数据为异常数据时才会执行策略
							if !isNormal {
								abnormalCount := abnormalStrategy(stateCdrs)
								if big.NewInt(30).Cmp(abnormalCount) != 1 {
									// 异常条数告警

									// 数据清空
									cdrStateMap[key] = append(cdrStateMap[key][:0])
								}
							}
						}

					}
				}
			}
		}
	}
}
