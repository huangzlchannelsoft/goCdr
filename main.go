// main
package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/natefinch/lumberjack.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	as             = flag.String("as", "", "run as client || server || singleton.")
	cdrPath        = flag.String("cdrPath", "", "cdr filepath that collected from.")
	svrAddr        = flag.String("svrAddr", "", "ip:port.")
	alarmUri       = flag.String("alarmUri", "", "send alarm via...")
	pushGateWayUri = flag.String("pushGateWayUri", "", "send alarm via...")
	phoneIspUri    = flag.String("phoneIspUri", "", "check phone's area via...")
	phoneProUri    = flag.String("phoneProUri", "", "check phone's manul via...")
)

type Config struct {
	TimeMinInterva	int64	`yaml:"timeMinInterva"`
	CdrCommaTotal 	int 	`yaml:"cdrCommaTotal"`
	FixedLineReg 	string 	`yaml:"fixedLineReg"`
	ConAbnormal 	int		`yaml:"conAbnormal"`
	Percentage      float64 `yaml:"percentage"`
	MobileReg 		string 	`yaml:"mobileReg"`
	Version 		string 	`yaml:"version"`
	Logfile 		bool   	`yaml:"logfile"`
	Bakdays 		int    	`yaml:"bakdays"`
	Nid     		string 	`yaml:"nid"`
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

/*
CollectCdr -TransmitCdr-> ParseCdr -TransmitAlarmCdr-
                                  |                  |->  [company]
                                   -TransmitMonCdr---
*/
func main() {
	flag.Parse()

	SetPhonePropertyUri(*phoneIspUri, *phoneProUri)

	ctx, cancel := context.WithCancel(context.Background())
	go TrickerDeamon(ctx)
	go PromethuesClient(true, *pushGateWayUri, gCfg.Nid, 60)

	go TransmitCdr(*as, *svrAddr)
	go TransmitAlarmCdr(*alarmUri)
	go TransmitMonCdr()

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
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			cancel()
			return
		}
	}
}
