// main
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"gopkg.in/natefinch/lumberjack.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	as            = flag.String("as", "singleton", "run as client || server || singleton.")
	cdrPath       = flag.String("cdrPath", "", "cdr filepath that collected from.")
	svrAddr       = flag.String("svrAddr", "", "ip:port.")
	alarmUri      = flag.String("alarmUri", "", "send alarm via...")
	phoneCheckUri = flag.String("phoneCheckUri", "", "check phone's area|manuf via...")
)

type Config struct {
	Version string `yaml:"version"`
	Logfile bool   `yaml:"logfile"`
	Bakdays int    `yaml:"bakdays"`
	Nid     string `yaml:"nid"`
}

var gCfg Config

func init() {
	buf, err := ioutil.ReadFile("cfg.yaml")
	if err != nil {
		log.Panic("load cfg.yaml error.", err)
	}

	err = yaml.Unmarshal(buf, &gCfg)
	if err != nil {
		log.Panic("parse cfg.yaml error.", err)
	}

	if gCfg.Logfile {
		log.SetOutput(&lumberjack.Logger{
			Filename:   "log/foo.log",
			MaxSize:    100, // megabytes
			MaxBackups: 0,
			MaxAge:     gCfg.Bakdays, //days
			Compress:   false,        // disabled by default
			LocalTime:  true,
		})
	}
}

func main() {
	flag.Parse()

	go TransmitCdr(*as, *svrAddr)
	go AlarmCdr(*alarmUri)

	if *as == "client" {
		go CollectCdr(*cdrPath, SendCdr)
	} else if *as == "server" {
		go ParseCdr(RecvCdr, SendAlarm)
	} else if *as == "singleton" {
		go CollectCdr(*cdrPath, SendCdr)
		go ParseCdr(RecvCdr, SendAlarm)
	}

	c := make(chan os.Signal)
	signal.Notify(c)
	select {
	case s := <-c:
		log.Println("process received signal:", s.String())
	}
}
