// parser
package main

import (
	"log"
	"strings"
	"regexp"
	"math/big"
	"fmt"
	"encoding/json"
	"time"
	"os"
)

func init() {
	log.Println("init parser!")
}

// 时间策略
const TIME_STRATEGY int8 = 0
// 异常条数策略
const ABNORMAL_STRATEGY int8 = 1
// 百分号
const TAGE  = string("%")
// 文件扩展名
const FILE_EXTENSION = string(".txt")
// 换行符
const NEWLINE_SYMBOL  = string("\r\n")
// tab
const TAB = string("\t")
// 符
const ARROW_SYMBOL = string(" ")

// cdr状态映射集
var cdrStateMap = make(map[string] []*CdrState)

// cdr状态结构体
type CdrState struct {
	// true 正常 false 异常
	isNormal 		bool
	// 当前时间戳
	curTimestamp 	int64
}

// 告警结构体
type AlarmInfo struct {
	// 被叫号码所在key
	key 			string
	// 策略类型 0-时间策略 1->异常条数理策略
	strategyType 	int8
	// strategyType == 0 ? 56% : 20
	value 			string
}

// 将统计出的key信息写入文件
func writeKeysInfo (path string,lineInfo string) string {
	log.Println("[INFO]","将要写入文件内容为：",lineInfo)
	now := time.Now()
	year, month, day := now.Date()
	hour := now.Hour()
	minute := now.Minute()
	Second := now.Second()

	nowStr := fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day,hour, minute, Second)
	fileNameStr := fmt.Sprintf("%d%d%d", year, month, day) + FILE_EXTENSION
	filePath := path + "/" + fileNameStr

	// 文件是或否存在
	_,err := os.Stat(filePath)
	if err == nil {
		// 存在的情况下追加
		log.Println("[WARN]","存在的情况下追加")
		file,err2 := os.OpenFile(filePath,os.O_APPEND,os.ModeAppend)
		defer file.Close()
		if err2 != nil {
			log.Println("[ERR]",err2)
		} else {
			file.WriteString(nowStr)
			file.WriteString(TAB)
			file.WriteString(lineInfo)
			file.WriteString(NEWLINE_SYMBOL)
		}
	} else {
		// 不存在创建写入
		log.Println("[WARN]","不存在创建写入")
		file,err := os.Create(path + "/" + fileNameStr)
		defer file.Close()
		if err != nil {
			log.Println("[ERR]",err)
		} else {
			file.WriteString(nowStr)
			file.WriteString(TAB)
			file.WriteString(lineInfo)
			file.WriteString(NEWLINE_SYMBOL)
		}
	}
	return fileNameStr
}



// 分割解析cdr数据
func parsingCdr (cdr string) map[string] string {
	log.Println("[INFO] 接收到cdr信息为：",cdr)
	resultMap := make(map[string] string,6)

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
	isornoMobile,_ := regexp.MatchString(gCfg.MobileReg, calledNumber)
	isornoFixed,_ := regexp.MatchString(gCfg.FixedLineReg, calledNumber)

	// 判断是否合法的手机号或固话
	if isornoMobile || isornoFixed {
		flag = true
	}
	return flag
}

// 连续异常条数状态策略
func abnormalStrategy (stateCdrs []*CdrState) int {
	log.Println("[WARN]","执行异常条数策略部分")

	// 异常数据发生条数
	abnormalCount := int(0)
	// 倒序遍历获取是或否连续异常
	for i := (len(stateCdrs) - 1); i >=0; i-- {
		stateCdr := stateCdrs[i]
		if !stateCdr.isNormal {
			abnormalCount ++
		} else {
			break
		}
	}
	log.Println("[INFO]","发生连续异常数据条数有：",abnormalCount)
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

	// 至少达到一定条数再触发计算
	if len(stateCdrs) >= gCfg.ConAbnormal {
		// 遍历异常条数
		for _,stateCdr := range stateCdrs {
			if !stateCdr.isNormal {
				abnormalCount = abnormalCount.Add(abnormalCount,big.NewFloat(1))
			}
		}
	}

	percentage = percentage.Quo(abnormalCount,cdrsLen)
	log.Println("[INFO]","时间策略计算值为：",percentage)
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
				isNormal := true
				// 判断话务信息是否正常
				if (strings.Compare(createTime,turnOnTime) == 0) && (strings.Compare(createTime,shutDownTime) == 0) {
					// 异常数据
					isNormal = false
				}
				log.Println("isNormal = ",isNormal)

				// 获取phoneKeys
				phonePro := GetPhoneProperty(callNumber,calledNumber)

				// 发送监控系统
				AddCallStat(phonePro,isNormal)

				keys := []string{phonePro.province,phonePro.productor,phonePro.isp,phonePro.area}

				// 当前时间戳
				curTimestamp := time.Now().Unix()
				log.Println("[INFO]","curTimestamp = ",curTimestamp)
				for _, key := range keys {
					log.Println("操作状态key值为：",key)

					// 告警百分比
					percentage := big.NewFloat(0.00)
					// 首次加入数据
					if cdrStateMap[key] == nil || (len(cdrStateMap[key]) == 0) {
						log.Println("[WARN]","当前key首次加入记录！")
						stateCdr := CdrState{isNormal,curTimestamp}
						stateCdrs := [] *CdrState{&stateCdr}
						// 追加
						cdrStateMap[key] = stateCdrs
					} else {
						log.Println("[WARN]","追加keyCdr记录")
						stateCdr := CdrState{isNormal,curTimestamp}
						// 追加并取出之前已存在数据
						stateCdrs := append(cdrStateMap[key], &stateCdr)
						firstCdrCurTime := stateCdrs[0].curTimestamp

						// 是否达到一定时间策略
						if (curTimestamp - firstCdrCurTime) >= (gCfg.TimeMinInterva * 60) {
							percentage = timeStrategy(stateCdrs)
							// 百分比
							finalPercentage := fmt.Sprintf("%0.2f",percentage.Mul(percentage,big.NewFloat(100))) + TAGE
							// 是否大于告警阈值
							if big.NewFloat(gCfg.Percentage).Cmp(percentage) != 1 {
								// 告警
								alarm := AlarmInfo{key,TIME_STRATEGY,finalPercentage}
								byteAlarm,err := json.Marshal(alarm)
								if err != nil {
									log.Println("[ERR]","err",err)
								}
								strAlarm := string(byteAlarm)
								log.Println("[INFO]","strAlarm = ",strAlarm)

								// 发送告警
								sendAlarm(strAlarm)
							}
							keyInfoLine := key + ARROW_SYMBOL + finalPercentage
							// 数据key存储
							writeKeysInfo(gCfg.StrategyInfoPath,keyInfoLine)

							// 数据清空
							cdrStateMap[key] = append(cdrStateMap[key][:0])
						} else {
							// 此条数据为异常数据时才会执行策略
							if !isNormal {
								abnormalCount := abnormalStrategy(stateCdrs)
								if abnormalCount >= gCfg.ConAbnormal {
									// 异常条数告警
									alarm := AlarmInfo{key,ABNORMAL_STRATEGY,string(abnormalCount)}
									byteAlarm,err := json.Marshal(alarm)
									if err != nil {
										log.Println("[ERR]","err",err)
									}
									strAlarm := string(byteAlarm)
									log.Println("[INFO]","strAlarm = ",strAlarm)

									// 发送告警
									sendAlarm(strAlarm)

									// 获取key存储
									keyInfoLine := key + ARROW_SYMBOL + string(abnormalCount)
									// 数据key存储
									writeKeysInfo(gCfg.StrategyInfoPath,keyInfoLine)

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
