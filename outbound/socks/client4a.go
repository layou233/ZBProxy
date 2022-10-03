package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

func (c *Client) handshake4A(r io.Reader, w io.Writer, address string) error {
	host, portString, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	port64, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return err
	}
	port := uint16(port64)

	if ip := net.ParseIP(host); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			return c.request4(r, w, port, ipv4)
		} else if ip.To16() != nil {
			return fmt.Errorf("socks: IPv6 is not supported in SOCKS verion 4/4a: %v", ip)
		}
		return fmt.Errorf("socks: unknown IP type: %v", ip)
	}

	// domain
	_, err = w.Write([]byte{version4, CommandConnect})
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, port)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{0, 0, 0, 233}) // magic, used to instruct the server to use 4A
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(c.Username))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{0})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(host)) // write domain
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{0})
	if err != nil {
		return err
	}
	return c.handleResponse4(r)
}
