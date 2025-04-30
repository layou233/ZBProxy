package proxyprotocol

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/layou233/zbproxy/v3/common/buf"
)

func TestClientV2(t *testing.T) {
	tests := []struct {
		name        string
		header      *Header
		destination netip.AddrPort
		expect      []byte
	}{
		{
			name: "TCP4 127.0.0.1",
			header: &Header{
				Version:           Version2,
				Command:           CommandProxy,
				TransportProtocol: TransportProtocolStream | TransportProtocolIPv4,
				SourceAddress:     netip.MustParseAddrPort("127.0.0.1:51755"),
			},
			destination: netip.MustParseAddrPort("127.0.0.1:1025"),
			//                                                                                     VER  IP/TCP LENGTH
			expect: []byte{
				0x0D, 0x0A, 0x0D, 0x0A, 0x00, 0x0D, 0x0A, 0x51, 0x55, 0x49, 0x54, 0x0A, 0x21, 0x11, 0x00, 0x0C,
				// IPV4 -------------|  IPV4 ----------------|   SRC PORT   DEST PORT
				0x7F, 0x00, 0x00, 0x01, 0x7F, 0x00, 0x00, 0x01, 0xCA, 0x2B, 0x04, 0x01,
			},
		},
		{
			name: "UDP4 127.0.0.1",
			header: &Header{
				Version:           Version2,
				Command:           CommandProxy,
				TransportProtocol: TransportProtocolDatagram | TransportProtocolIPv4,
				SourceAddress:     netip.MustParseAddrPort("127.0.0.1:51755"),
			},
			destination: netip.MustParseAddrPort("127.0.0.1:1025"),
			//                                                                                          IP/UDP
			expect: []byte{
				0x0D, 0x0A, 0x0D, 0x0A, 0x00, 0x0D, 0x0A, 0x51, 0x55, 0x49, 0x54, 0x0A, 0x21, 0x12, 0x00, 0x0C,
				0x7F, 0x00, 0x00, 0x01, 0x7F, 0x00, 0x00, 0x01, 0xCA, 0x2B, 0x04, 0x01,
			},
		},
		{
			name: "TCP6 Proxy for TCP4 127.0.0.1",
			header: &Header{
				Version:           Version2,
				Command:           CommandProxy,
				TransportProtocol: TransportProtocolStream | TransportProtocolIPv6,
				SourceAddress:     netip.MustParseAddrPort("[::ffff:127.0.0.1]:52300"),
			},
			destination: netip.MustParseAddrPort("[::ffff:127.0.0.1]:1025"),
			//                                                                                     VER  IP/TCP   LENGTH
			expect: []byte{
				0x0D, 0x0A, 0x0D, 0x0A, 0x00, 0x0D, 0x0A, 0x51, 0x55, 0x49, 0x54, 0x0A, 0x21, 0x21, 0x00, 0x24,
				// IPV6 -------------------------------------------------------------------------------------|
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x7F, 0x00, 0x00, 0x01,
				// IPV6 -------------------------------------------------------------------------------------|   SRC PORT   DEST PORT
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x7F, 0x00, 0x00, 0x01, 0xCC, 0x4C, 0x04, 0x01,
			},
		},
		{
			name: "TCP6 Maximal",
			header: &Header{
				Version:           Version2,
				Command:           CommandProxy,
				TransportProtocol: TransportProtocolStream | TransportProtocolIPv6,
				SourceAddress:     netip.MustParseAddrPort("[FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF]:65535"),
			},
			destination: netip.MustParseAddrPort("[FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF]:65535"),
			//                                                                                     VER  IP/TCP   LENGTH
			expect: []byte{
				0x0D, 0x0A, 0x0D, 0x0A, 0x00, 0x0D, 0x0A, 0x51, 0x55, 0x49, 0x54, 0x0A, 0x21, 0x21, 0x00, 0x24,
				// IPV6 -------------------------------------------------------------------------------------|
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				// IPV6 -------------------------------------------------------------------------------------|   SRC PORT   DEST PORT
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			},
		},
		{
			name: "TCP6 Proxy for TCP6 ::1",
			header: &Header{
				Version:           Version2,
				Command:           CommandProxy,
				TransportProtocol: TransportProtocolStream | TransportProtocolIPv6,
				SourceAddress:     netip.MustParseAddrPort("[::1]:53135"),
			},
			destination: netip.MustParseAddrPort("[::1]:1025"),
			//                                                                                     VER  IP/TCP   LENGTH
			expect: []byte{
				0x0D, 0x0A, 0x0D, 0x0A, 0x00, 0x0D, 0x0A, 0x51, 0x55, 0x49, 0x54, 0x0A, 0x21, 0x21, 0x00, 0x24,
				// IPV6 -------------------------------------------------------------------------------------|
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				// IPV6 -------------------------------------------------------------------------------------|   SRC PORT   DEST PORT
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xCF, 0x8F, 0x04, 0x01,
			},
		},
		{
			name: "UDP6 Proxy for UDP6 ::1",
			header: &Header{
				Version:           Version2,
				Command:           CommandProxy,
				TransportProtocol: TransportProtocolDatagram | TransportProtocolIPv6,
				SourceAddress:     netip.MustParseAddrPort("[::1]:53135"),
			},
			destination: netip.MustParseAddrPort("[::1]:1025"),
			//                                                                                     VER  IP/UDP   LENGTH
			expect: []byte{
				0x0D, 0x0A, 0x0D, 0x0A, 0x00, 0x0D, 0x0A, 0x51, 0x55, 0x49, 0x54, 0x0A, 0x21, 0x22, 0x00, 0x24,
				// IPV6 -------------------------------------------------------------------------------------|
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				// IPV6 -------------------------------------------------------------------------------------|   SRC PORT   DEST PORT
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xCF, 0x8F, 0x04, 0x01,
			},
		},
		{
			name: "Local with no trailing bytes",
			header: &Header{
				Version: Version2,
				Command: CommandLocal,
			},
			expect: []byte{0x0D, 0x0A, 0x0D, 0x0A, 0x00, 0x0D, 0x0A, 0x51, 0x55, 0x49, 0x54, 0x0A, 0x20, 0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buf.New()
			defer buffer.Release()
			err := tt.header.WriteHeader(buffer, tt.destination)
			if err != nil {
				t.Fatalf("failed to write v2 header: %v", err)
			}
			if !bytes.Equal(buffer.Bytes(), tt.expect) {
				t.Errorf("got=%x, expect=%x", buffer.Bytes(), tt.expect)
			}
		})
	}
}
