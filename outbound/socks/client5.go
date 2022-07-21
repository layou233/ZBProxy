package socks

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/layou233/ZBProxy/common/rw"
	"io"
	"net"
	"strconv"
)

const (
	version5 byte = 5

	AuthTypeNotRequired       byte = 0x00
	AuthTypeGSSAPI            byte = 0x01
	AuthTypeUsernamePassword  byte = 0x02
	AuthTypeNoAcceptedMethods byte = 0xFF

	UsernamePasswordStatusSuccess byte = 0x00
	UsernamePasswordStatusFailure byte = 0x01

	CommandConnect      byte = 0x01
	CommandBind         byte = 0x02
	CommandUDPAssociate byte = 0x03

	AddressTypeIPv4   byte = 0x01
	AddressTypeDomain byte = 0x03
	AddressTypeIPv6   byte = 0x04

	ReplyCode5Success                byte = 0
	ReplyCode5Failure                byte = 1
	ReplyCode5NotAllowed             byte = 2
	ReplyCode5NetworkUnreachable     byte = 3
	ReplyCode5HostUnreachable        byte = 4
	ReplyCode5ConnectionRefused      byte = 5
	ReplyCode5TTLExpired             byte = 6
	ReplyCode5Unsupported            byte = 7
	ReplyCode5AddressTypeUnsupported byte = 8
)

func (c Client) handshake5(r io.Reader, w io.Writer, network, address string) error {
	host, portString, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	port64, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return err
	}
	port := uint16(port64)

	if len(c.Methods) == 0 {
		c.Methods = []byte{AuthTypeNotRequired}
	}
	_, err = w.Write([]byte{version5, byte(len(c.Methods))})
	if err != nil {
		return err
	}
	_, err = w.Write(c.Methods)
	if err != nil {
		return err
	}

	handshakeResp := make([]byte, 2)
	_, err = r.Read(handshakeResp)
	if err != nil {
		return err
	}
	if handshakeResp[0] != version5 {
		return fmt.Errorf("socks: expected server version 5, but got %v", handshakeResp[0])
	}
	switch handshakeResp[1] {
	case AuthTypeNotRequired:
	case AuthTypeUsernamePassword:
		return errors.New("socks: username/password auth is not implemented yet")
	case AuthTypeGSSAPI:
		return errors.New("socks: unsupported SOCKS 5 auth type GSSAPI")
	case AuthTypeNoAcceptedMethods:
		return errors.New("socks: server responded no acceptable methods")
	default:
		return fmt.Errorf("socks: unknown auth method: %v", handshakeResp[1])
	}

	_, err = w.Write([]byte{version5, CommandConnect, 0}) // TODO: Bind & UDPAssociate
	if err != nil {
		return err
	}
	if ip := net.ParseIP(host); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			_, err = w.Write([]byte{AddressTypeIPv4})
			if err != nil {
				return err
			}
			_, err = w.Write(ipv4)
			if err != nil {
				return err
			}
		} else if ipv6 := ip.To16(); ipv6 != nil {
			_, err = w.Write([]byte{AddressTypeIPv6})
			if err != nil {
				return err
			}
			_, err = w.Write(ipv6)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("socks: unknown IP type: %v", ip)
		}
	} else {
		_, err = w.Write([]byte{AddressTypeDomain, byte(len(host))})
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(host))
		if err != nil {
			return err
		}
	}
	err = binary.Write(w, binary.BigEndian, port)
	if err != nil {
		return err
	}

	reply := make([]byte, 4)
	_, err = r.Read(reply)
	if err != nil {
		return err
	}
	if reply[0] != version5 {
		return fmt.Errorf("socks: expected server version 5, but got %v", handshakeResp[0])
	}
	switch reply[1] {
	case ReplyCode5Success:
	default:
		return fmt.Errorf("socks: fail to connect destination through SOCKS 5, reason code: %v", reply[1])
	}
	switch reply[3] {
	case AddressTypeIPv4:
		ipv4 := make([]byte, 4)
		_, err = r.Read(ipv4)
		if err != nil {
			return err
		}
	case AddressTypeDomain:
		l, err := rw.ReadByte(r)
		if err != nil {
			return err
		}
		domain := make([]byte, l)
		_, err = r.Read(domain)
		if err != nil {
			return err
		}
	case AddressTypeIPv6:
		ipv6 := make([]byte, 16)
		_, err = r.Read(ipv6)
		if err != nil {
			return err
		}
	}
	_, err = rw.ReadBytes(r, 2)
	if err != nil {
		return err
	}
	return nil
}
