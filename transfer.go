package main

import (
	"log"
	"net"
)

func forDial(fromConn net.Conn, forAddr string) {
	toConn, err := net.Dial("tcp", forAddr)
	if err != nil {
		log.Printf("[Closed] %s to %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
		toConn.Close()
		return
	}
	log.Printf("[Transfer started] %s to %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
	go transfer(fromConn, toConn, 1024)
	go transfer(toConn, fromConn, 1024)
}

func toDial(fromConn net.Conn) {
	toAddr := fromConn.RemoteAddr()
	toConn, err := net.Dial("tcp", toAddr.String())
	if err != nil {
		log.Printf("[Closed] %s to %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
		toConn.Close()
		return
	}
	log.Printf("%s to %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
	go transfer(fromConn, toConn, 1024)
	go transfer(toConn, fromConn, 1024)
}

func transfer(f, t net.Conn, n int) {
	firstConn := true
	onlinePlayers++
	defer func() { onlinePlayers-- }()
	defer f.Close()
	defer t.Close()
	log.Printf("[Closed] %s to %s", f.RemoteAddr().String(), t.RemoteAddr().String())

	var buf = make([]byte, n)

	for {
		count, err := f.Read(buf)
		if err != nil {
			break
		}
		if firstConn {
			log.Println(buf)
			firstConn = false
		}
		count, err = t.Write(buf[:count])
		if err != nil {
			log.Printf("fault,err: %s", err.Error())
			break
		}
	}
}
