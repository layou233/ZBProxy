// Package proxyprotocol partially implements PROXY protocol proposed by HAProxy,
// aiming at high-performance integration with the rest of ZBProxy.
// The full protocol specification can be found at https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt
package proxyprotocol

import (
	"errors"
	"net/netip"
)

const (
	VersionUnspecified = iota
	Version1
	Version2
)

const (
	maskVersion  = 0xF0
	maskVersion2 = 0x20

	maskCommand  = 0x0F
	CommandLocal = 0x0
	CommandProxy = 0x1

	transportProtocolUnspecified       = 0x00
	maskTransportProtocolAddressFamily = 0xF0
	maskTransportProtocolType          = 0xF
	TransportProtocolIPv4              = 0x10
	TransportProtocolIPv6              = 0x20
	TransportProtocolUnix              = 0x30
	TransportProtocolStream            = 0x1
	TransportProtocolDatagram          = 0x2
)

var ErrNotProxyProtocol = errors.New("not PROXY protocol")

type Header struct {
	Version           uint8
	Command           uint8
	TransportProtocol uint8
	SourceAddress     netip.AddrPort
	//DestinationAddress netip.AddrPort
}

func TransportProtocolByNetwork(network string) uint8 {
	// see documentation of net.Dial
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
		return TransportProtocolStream
	case "udp", "udp4", "udp6", "unixgram":
		return TransportProtocolDatagram
	default:
		return transportProtocolUnspecified
	}
}

func AddressFamilyByAddr(addr netip.Addr) uint8 {
	if addr.Is4() {
		return TransportProtocolIPv4
	}
	return TransportProtocolIPv6
}
