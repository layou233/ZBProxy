package socks

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/layou233/ZBProxy/common/rw"
)

const (
	version4 byte = 4

	ReplyCode4Granted                     byte = 0x5A
	ReplyCode4RejectedOrFailed            byte = 0x5B
	ReplyCode4CannotConnectToIdentd       byte = 0x5C
	ReplyCode4IdentdReportDifferentUserID byte = 0x5D
)

func (c *Client) handshake4(r io.Reader, w io.Writer, address string) error {
	host, portString, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	port64, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return err
	}
	port := uint16(port64)

	ip := net.ParseIP(host)
	if ip == nil || ip.To4() == nil {
		if ip.To16() != nil {
			return fmt.Errorf("socks: IPv6 is not supported in SOCKS verion 4/4a: %v", ip)
		}
		ipList, err := net.LookupIP(host)
		if err != nil {
			return err
		}
		for _, ip := range ipList {
			if ipv4 := ip.To4(); ipv4 != nil {
				return c.request4(r, w, port, ipv4)
			}
		}
		return fmt.Errorf("socks: can't resolve any IPv4 address from domain %v: %v", host, ipList)
	}
	return c.request4(r, w, port, ip.To4())
}

func (c *Client) request4(r io.Reader, w io.Writer, port uint16, addr []byte) error {
	_, err := w.Write([]byte{version4, CommandConnect})
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, port)
	if err != nil {
		return err
	}
	_, err = w.Write(addr)
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
	return c.handleResponse4(r)
}

func (c *Client) handleResponse4(r io.Reader) error {
	resp, err := rw.ReadBytes(r, 2)
	if err != nil {
		return err
	}
	if resp[0] != 0 && resp[0] != version4 { // compatible with nonstandard implementation
		return fmt.Errorf("socks: expected response version 0, but got: %v", resp[0])
	}
	switch resp[1] {
	case ReplyCode4Granted:
	case ReplyCode4RejectedOrFailed:
		return errors.New("socks: connection request rejected or failed")
	case ReplyCode4CannotConnectToIdentd:
		return errors.New("socks: connection request rejected because SOCKS server cannot connect to identd on the client")
	case ReplyCode4IdentdReportDifferentUserID:
		return errors.New("socks: connection request rejected because the client program and identd report different user-ids")
	default:
		return fmt.Errorf("socks: unknown SOCKS 4 reply code: %v", resp[1])
	}
	_, err = rw.ReadBytes(r, 6)
	if err != nil {
		return err
	}

	return nil
}
