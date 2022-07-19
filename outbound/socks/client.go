package socks

import (
	"fmt"
	"io"
	"net"
	"net/url"
)

type Client struct {
	Version            string
	Network            string
	Address            string
	Username, Password string
	Methods            []byte
}

// NewClientFromURL
//
// socks(5)://username:password@127.0.0.1:1080
//
// socks4(a)://userid@127.0.0.1:1080
func NewClientFromURL(s string) (*Client, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	c := &Client{}
	switch u.Scheme {
	case "socks", "socks5":
		c.Version = "5"
		c.Username = u.User.Username()
		c.Password, _ = u.User.Password()
	case "socks4a":
		c.Version = "4a"
		c.Username = u.User.Username()
	case "socks4":
		c.Version = "4"
		c.Username = u.User.Username()
	default:
		return nil, fmt.Errorf("socks: unknown SOCKS version: %v", u.Scheme)
	}
	c.Address = u.Host
	return c, nil
}

func (c Client) Dial(network, address string) (net.Conn, error) {
	conn, err := net.Dial(network, c.Address)
	if err != nil {
		return nil, fmt.Errorf("socks: fail to dial to SOCKS server: %v", err)
	}
	if err = c.Handshake(conn, conn, address); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c Client) DialTCP(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP(network, laddr, raddr)
	if err != nil {
		return nil, fmt.Errorf("socks: fail to dial to SOCKS server: %v", err)
	}
	if err = c.Handshake(conn, conn, raddr.String()); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c Client) Handshake(r io.Reader, w io.Writer, address string) error {
	switch c.Version {
	case "5":
		return c.handshake5(r, w, address)
	}
	return fmt.Errorf("socks: unknown SOCKS version: %v", c.Version)
}