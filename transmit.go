// transmit
package main

import (
	"log"
	"net"
)

var (
	cdrCh chan string = make(chan string, 1000)
)

type CdrSend func(cdr string)
type CdrRecv func() string

func init() {
	log.Println("init transmit!")
}

func TransmitCdr(mode string, addr string) {
	if *as == "client" {

		connect := func() (bool, net.Conn) {
			conn, err := net.Dial("tcp4", addr)
			if err != nil {
				log.Println("[Err]transmit.", err.Error())
				return false, nil
			} else {
				return true, conn
			}
		}

		connected := func(conn net.Conn) {
			defer conn.Close()

			cdr := RecvCdr()
			_, err := conn.Write([]byte(cdr))

			if err != nil { //TODO: write cdr error
				log.Println("[Err] write cdr.", err.Error())
			}
		}

		netOk := false
		var conn net.Conn = nil
		for {
			if !netOk {
				netOk, conn = connect()
			}

			if netOk {
				connected(conn)
			}
		}
	} else if *as == "server" {
		buf := make([]byte, 1000)

		connected := func(conn net.Conn) {
			defer conn.Close()

			n, err := conn.Read(buf)
			if err != nil {
				log.Println("[Err] read cdr.", err.Error())
			}

			SendCdr(string(buf[:n]))
		}

		listener, err := net.Listen("tcp4", addr)
		if err != nil {
			log.Panic("[Err] transmit.", err.Error())
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Panic("[Err] transmit.", err.Error())
			}

			go connected(conn)
		}
	} else if *as == "singleton" {
		return
	}
}

func SendCdr(cdr string) {
	cdrCh <- cdr
}

func RecvCdr() string {
	return <-cdrCh
}
