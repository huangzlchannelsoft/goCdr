// transmit
package main

import (
	"bufio"
	"bytes"
	"io"
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
			log.Println("connected to ", conn.RemoteAddr())
			defer conn.Close()

			for {
				cdr := RecvCdr()
				_, err := conn.Write([]byte(cdr))

				if err != nil { //TODO: write cdr error
					log.Println("[Err] write cdr.", err.Error())
					break
				}
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

		connected := func(conn net.Conn) {
			log.Println("connected to ", conn.RemoteAddr())
			defer conn.Close()

			//split tcp stream
			if true {
				totel := 0
				split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
					if totel < 14 {
						totel = bytes.Count(data, []byte{','})
					}

					if totel >= 14 {
						c := 14
						l := 0
						p := data
						for c > 0 {
							n := bytes.IndexByte(p, ',') + 1
							p = p[n:]
							l += n
							c--
							totel--
						}
						return l, data[:l], nil
					} else if atEOF {
						return 0, nil, io.EOF
					}
					return 0, nil, nil
				}

				scanner := bufio.NewScanner(conn)
				scanner.Split(split)
				for scanner.Scan() {
					cdr := scanner.Text()
					SendCdr(cdr)
					//log.Println(cdr)
				}
			}

			// for {
			// 	buf := make([]byte, 1000)
			// 	n, err := conn.Read(buf)
			// 	if err != nil {
			// 		log.Println("[Err] read cdr.", err.Error())
			// 		break
			// 	}
			// 	//log.Println(string(buf[:n]))
			// 	SendCdr(string(buf[:n]))
			// }
		}

		listener, err := net.Listen("tcp4", addr)
		if err != nil {
			log.Panic("[Err] transmit.", err.Error())
		}
		for {
			log.Println("listen to ", addr)
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
