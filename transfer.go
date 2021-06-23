package main

import (
	"bytes"
	"log"
	"net"
)

func forDial(fromConn net.Conn, forAddr string) {
	toConn, err := net.Dial("tcp", forAddr)
	if err != nil {
		log.Printf("[Bad Connection] %s to %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
		toConn.Close()
		return
	}
	log.Printf("[Transfer started] %s to %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
	go transfer(fromConn, toConn, 1024, true)
	go transfer(toConn, fromConn, 1024, false)
}

/*
func toDial(fromConn net.Conn) {
	toConn, err := net.Dial("tcp", fromConn.RemoteAddr().String())
	if err != nil {
		log.Printf("[Closed] %s to %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
		toConn.Close()
		return
	}
	log.Printf("%s to %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
	go transfer(fromConn, toConn, 1024, true)
	go transfer(toConn, fromConn, 1024, false)
}*/

func transfer(f, t net.Conn, n int, isFrom2to bool) {
	firstConn, secondConn := true, false
	onlineConnections++
	defer func() { onlineConnections-- }()
	defer f.Close()
	defer t.Close()

	var buf = make([]byte, n)
	for {
		count, err := f.Read(buf)
		if err != nil {
			break
		}
		if firstConn {
			if isFrom2to && buf[1] == 0 {
				addressLength := DecodeVarint(buf, 3)
				//log.Println(addressLength)
				buf = bytes.Join([][]byte{
					buf[:3],
					{(byte)(len(serverAddr))},
					[]byte(serverAddr),
					{byte(serverPort >> 8), byte(serverPort & 0xff)}, // uint16 to []byte aka []uint8
					buf[3+addressLength+2+1:],                        // 2 is the length of 2* unsigned short (uint16)
				}, make([]byte, 0))
				buf[0] = (byte)((int)(buf[0]) + len(serverAddr) - addressLength)
				count += len(serverAddr) - addressLength
			}
			//log.Println(buf)
			firstConn = false
			secondConn = true
		} else if secondConn {
			log.Printf("[ATTENTION] A new player has joined the game.")
			//log.Println(buf)
			defer func() { log.Printf("[Closed] %s to %s", f.RemoteAddr().String(), t.RemoteAddr().String()) }()
			secondConn = false
		}
		count, err = t.Write(buf[:count])
		if err != nil {
			log.Printf("err: %s", err.Error())
			break
		}
	}
}
