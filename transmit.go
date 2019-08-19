// transmit
package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"time"
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

		log.Println("enter transmit client.")
		trigger := ApplyTrigger("transmit")
		exit := false
		go func() {
			defer func() {
				exit = true
				log.Println("Close connected net.Conn. transmit")
				if conn != nil {
					conn.Close()
				}

				*trigger <- TRIGGER_BYE_OK
			}()

			for {
				x := <-*trigger
				if x == TRIGGER_BYE_BYE {
					log.Println("doing for Process Exit. transmit")
					break
				}
			}
		}()

		defer func() {
			log.Println("exit transmit server.")
		}()

		for !exit {
			if !netOk {
				netOk, conn = connect()
			}

			if netOk {
				connected(conn)
				netOk = false
			} else {
				time.Sleep(time.Second)
			}
		}
	} else if *as == "server" {

		var conns []*net.Conn
		add := func(conn *net.Conn) {
			conns = append(conns, conn)
		}

		connected := func(conn net.Conn) {
			log.Println("connected to ", conn.RemoteAddr())
			add(&conn)
			//defer conn.Close()

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
		//defer listener.Close()

		log.Println("enter transmit server.")
		trigger := ApplyTrigger("transmit")
		go func() {
			defer func() {
				log.Println("Close connected net.Conn. transmit")
				for _, conn := range conns {
					(*conn).Close()
				}
				conns = nil

				log.Println("Close net.listener. transmit")
				listener.Close()

				*trigger <- TRIGGER_BYE_OK
			}()

			for {
				x := <-*trigger
				if x == TRIGGER_BYE_BYE {
					log.Println("doing for Process Exit. transmit")
					break
				}
			}
		}()

		defer func() {
			log.Println("exit transmit server.")
		}()

		for {
			log.Println("listen to ", addr)
			conn, err := listener.Accept()
			if err != nil {
				log.Println("[Err] listener.Accept. transmit.", err.Error())
				break
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
