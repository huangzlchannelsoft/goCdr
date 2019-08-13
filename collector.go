// collector
package main

import (
	"log"
	"time"
	"os"
	"bufio"
	"path/filepath"
	"strings"
)

// 上次读取行号
var lastLineNum = int32(0)
// 上一个文件名称
var lastFileName = ""

func init() {
	log.Println("init collector!")
}

// 单行读取文件
func readFileForLine (path string,fileName string,transmit CdrSend) error {
	log.Printf("[INFO] 接收到文件路径：<%d>,文件名称为：[%d]",path,fileName)
	// 开启文件流
	file,err := os.Open(path)
	if err != nil {
		log.Println(" [ERR] open File IO stream ERROR!",err)
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Println(" [ERR] close File IO stream ERROR!",err)
		}
	}()

	lastFileName = fileName

	// 逐行读取
	s := bufio.NewScanner(file)
	lineNums := lastLineNum;
	log.Printf("lineNums = %d",lineNums)

	// 当前行号
	curLineNum := int32(0)
	for s.Scan() {
		curLineNum ++
		if lineNums >= curLineNum {
			continue
		}
		line := string(s.Text())

		if strings.Count(line,",") != gCfg.CdrCommaTotal {
			log.Println("[WARN] 这是不完整的一行记录")
			log.Println("[WARN] ",line)
			curLineNum --
			lastLineNum = curLineNum
			break
		} else {
			// 完整的一段
			log.Println(line)
			cdr := line
			transmit(cdr)
		}
	}
	return nil
}

func CollectCdr(path string, transmit CdrSend) {
	log.Println("run collectCdr.")

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			//check new cdr
			log.Printf("[INFO] 收集cdr文件根目录为：<%d>",path)
			// 获取文件
			err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				log.Printf("[INFO] 遍历得到路径：<%d>",path)
				if !info.IsDir() {
					fileName := info.Name()
					log.Printf("[WARN] fileName = %d,lastFileName = %d",fileName,lastFileName)

					// 规避首次文件名为空 此次文件名去上次文件不是同一文件
					if lastFileName != "" && !strings.EqualFold(fileName,lastFileName)  {
						log.Println("[WARN] 返回上一个文件读取")

						// 去上级目录（父目录）下那最后的那个lastFileName文件 处理 最后 break
						fIndex1 := strings.LastIndex(path,"\\")
						fatherPath1 := path[:fIndex1]
						fIndex2 := strings.LastIndex(fatherPath1,"\\")
						fatherPath2 := fatherPath1[:fIndex2]
						log.Printf("父目录地址为: <%d>",fatherPath2)

						err = readFileForLine(fatherPath2 + "\\" + lastFileName,lastFileName,transmit)
						if err != nil {
							log.Println("[ERR]",err)
						}
						// 此次文件名称赋值 保证下次在当前子目录下查询
						lastFileName = fileName
						// 新文件初始行
						lastLineNum = 0
					} else {
						err := readFileForLine(path,fileName,transmit)
						if err != nil {
							log.Println("[ERR]",err)
						}
					}
				}
				return nil
			})
			if err != nil {
				log.Println("[ERR] 获取cdr话单根目录文件异常",err)
			}
		}
	}
}
