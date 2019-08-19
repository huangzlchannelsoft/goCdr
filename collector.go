// collector
package main

import (
	"log"
	"time"
	"os"
	"bufio"
	"path/filepath"
	"strings"
	"io/ioutil"
)

// 上次读取行号
var lastLineNum = int32(0)
// 上一个文件名称
var lastFileName = ""
// 上一个文件大小
var lastFileSize = int64(0)

func init() {
	log.Println("init collector!")
}

// 单行读取文件
func readFileForLine (path string,fileName string,transmit CdrSend) error {
	log.Println("[INFO] 接收到文件路径：" + path + ",文件名称为：",path,fileName)

	fileBakName := string(gCfg.CdrFileBakPath + "/" + fileName)

	file,err := os.Stat(path)
	bakFile,errBak := os.Stat(fileBakName)
	if errBak == nil && (file.Size() == bakFile.Size()) && (strings.Compare(fileName,lastFileName) == 0){
		log.Println("[WARN]","文件" + fileName + "存在且大小未发生变化，不进行处理！")
	} else {
		log.Println("[WARN]","文件存在且大小发生变化，进行处理！")
		// 复制文本
		input,err1 := ioutil.ReadFile(path)
		if err1 != nil {
			log.Println("[ERR]","复制文本读取文件流异常！",err1)
			return err1
		}
		err2 := ioutil.WriteFile(fileBakName,input,0644)
		if err2 != nil {
			log.Println("[ERR]","复制文本写入文件流异常！",err2)
			return err2
		}

		// 开启文件流
		file,err4 := os.Open(fileBakName)
		if err4 != nil {
			log.Println(" [ERR] open File IO stream ERROR!",err4)
			return err4
		}
		defer func() {
			if err = file.Close(); err != nil {
				log.Println(" [ERR] close File IO stream ERROR!",err)
			}
		}()

		lastFileName = fileName
		fileInfo,err3 := file.Stat()
		if err3 != nil {
			log.Println("[ERR]","获取文件流异常！",err3)
			return err3
		}
		lastFileSize = fileInfo.Size()

		// 逐行读取
		s := bufio.NewScanner(file)
		lineNums := lastLineNum;
		log.Println("lineNums = ",lineNums)

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
			} else {
				// 完整的一段
				log.Println(line)
				cdr := line
				transmit(cdr)
			}
			// 赋值行号
			lastLineNum = curLineNum
		}
	}
	return nil
}

func CollectCdr(path string, transmit CdrSend) {
	log.Println("run collectCdr.")

	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			//check new cdr
			log.Println("[INFO] 收集cdr文件根目录为：",path)
			// 获取文件
			err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				log.Println("[INFO] 遍历得到路径：",path)
				if !info.IsDir() {
					fileName := info.Name()
					log.Println("[WARN] fileName = ",fileName,",lastFileName =",lastFileName)

					// 规避首次文件名为空 此次文件名去上次文件不是同一文件
					if lastFileName != "" && !strings.EqualFold(fileName,lastFileName)  {
						log.Println("[WARN] 返回上一个文件读取")

						// 去上级目录（父目录）下那最后的那个lastFileName文件 处理 最后 break
						fIndex1 := strings.LastIndex(path,"/")
						fatherPath1 := path[:fIndex1]
						fIndex2 := strings.LastIndex(fatherPath1,"/")
						fatherPath2 := fatherPath1[:fIndex2]
						log.Println("父目录地址为: ",fatherPath2)

						err = readFileForLine(fatherPath2 + "/" + lastFileName,lastFileName,transmit)
						if err != nil {
							log.Println("[ERR]",err)
							return err
						}
						// 移除上一个备份文件
						fileBakName := string(gCfg.CdrFileBakPath + "/" + lastFileName)
						errDel := os.Remove(fileBakName)
						if errDel != nil {
							log.Println("[ERR]","备份文件删除异常！")
							return errDel
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
