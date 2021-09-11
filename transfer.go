package main

import (
	"bytes"
	"log"
	"net"
	"time"
)

func forDial(fromConn net.Conn, forAddr string, retryTimes uint8) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic: %s\n", p)
			if retryTimes < 5 {
				forDial(fromConn, forAddr, retryTimes+1)
			}
		}
	}()
	toConn, err := net.DialTimeout("tcp", forAddr, 10*time.Second)
	if err != nil {
		log.Printf("[Bad Connection] %s -> %s", fromConn.LocalAddr().String(), toConn.RemoteAddr().String())
		toConn.Close()
		return
	}
	log.Printf("[Transfer started] %s -> %s", fromConn.RemoteAddr().String(), toConn.RemoteAddr().String())
	go transfer(fromConn, toConn, 4096, true)
	go transfer(toConn, fromConn, 4096, false)
}

func transfer(f, t net.Conn, n int, isFrom2to bool) {
	// firstConn, secondConn := true, false
	onlineConnections++
	defer func() { onlineConnections-- }()
	defer f.Close()
	defer t.Close()

	var buf = make([]byte, n)

	if isFrom2to {
		count, err := f.Read(buf)
		if err != nil {
			log.Printf("err: %s", err.Error())
			return
		}
		packetLength, startIndex := DecodeVarint(buf, 0)
		//log.Println(buf[:count])
		//log.Println(packetLength)
		if buf[startIndex] == 0 {
			/*
				Client first post data to the server.
				And the server address is included in this packet.
				In this situation, we need to locale the server address and change it.
			*/
			var addressLength, protocolVersionLength int
			if (buf[2] & 0b10000000) == 0 {
				addressLength, _ = DecodeVarint(buf, 3)
				protocolVersionLength = 1
			} else {
				addressLength, _ = DecodeVarint(buf, 4)
				protocolVersionLength = 2
			}
			//log.Println(addressLength)
			newPacketLength := packetLength + len(ServerAddr) - addressLength
			buf = bytes.Join([][]byte{
				EncodeVarint(newPacketLength),
				buf[startIndex : startIndex+1+protocolVersionLength], // includes Packet ID and Protocol Version
				{(byte)(len(ServerAddr))},
				[]byte(ServerAddr),
				{byte(ServerPort >> 8), byte(ServerPort & 0xff)}, // uint16 to []byte aka []uint8
				buf[3+addressLength+2+protocolVersionLength:],    // 2 is the length of 2* unsigned short (uint16)
			}, make([]byte, 0))
			count += newPacketLength - packetLength
		} //else { //TODO(MOTD FUNCTION NOT FINISHED YET)
		/*
			Server respond the ping request that requested by client.
			And all the motd information is included in this packet.
			We can rewrite it in order to change the look of the server title.
		*/ /*
					jsonLength, jsonStartIndex := DecodeVarint(buf, startIndex+1)
					jsonStartIndex += startIndex + 1
					motdJson := string(buf[jsonStartIndex:count])
					log.Printf("origin data,%v", motdJson)
					motdJsonLength := len(motdJson)
					motdDescriptionIndex := strings.Index(motdJson, `description":`)
					motdFaviconIndex := strings.Index(motdJson, `favicon":`)
					if IsChangeDescription && IsChangeFavicon {
						motdJson = strings.Join([]string{
							motdJson[:motdDescriptionIndex-1],
							`description":"`,
							MotdDescription,
							`","favicon":"`,
							MotdFavicon,
							`"}`,
						}, "")
					} else if IsChangeDescription {
						motdJson = strings.Join([]string{
							motdJson[:motdDescriptionIndex-1],
							`description":"`,
							MotdDescription,
							`","`,
							motdJson[motdFaviconIndex:],
						}, "")
					} else { // IsChangeFavicon
						motdJson = strings.Join([]string{
							motdJson[:motdFaviconIndex-1],
							`favicon":"`,
							MotdFavicon,
							`"}`,
						}, "")
					}
					lengthDiscrepancy := len(motdJson) - motdJsonLength
					newPacketLengthArray := EncodeVarint(packetLength + lengthDiscrepancy)
					buf = bytes.Join([][]byte{
						newPacketLengthArray,
						{0},
						EncodeVarint(jsonLength + lengthDiscrepancy),
						[]byte(motdJson),
					}, make([]byte, 0))
					count += len(newPacketLengthArray) - startIndex + lengthDiscrepancy
				}
		}*/

		//log.Println(buf[:count])
		_, err = t.Write(buf[:count])
		if err != nil {
			log.Printf("err: %s", err.Error())
			return
		}
	}
	defer func() { log.Printf("[Closed] %s -> %s", f.RemoteAddr().String(), t.RemoteAddr().String()) }()

	for {
		count, err := f.Read(buf)
		if err != nil {
			break
		}
		_, err = t.Write(buf[:count])
		if err != nil {
			log.Printf("err: %s", err.Error())
			break
		}
	}
}
