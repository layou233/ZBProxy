package tls

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"github.com/layou233/ZBProxy/common/rw"
)

var ErrNotTLS = errors.New("not TLS header")

func SniffAndRecordTLS(conn net.Conn) (*SniffHeader, *bytes.Buffer, error) {
	var recorder bytes.Buffer
	recorder.Grow(64)
	b, err := rw.ReadByte(conn)
	if err != nil {
		recorder.Reset()
		return nil, nil, err
	}
	err = recorder.WriteByte(b)
	if err != nil {
		recorder.Reset()
		return nil, nil, err
	}
	if b != 0x16 { // TLS Handshake
		return nil, &recorder, ErrNotTLS
	}
	bs, err := rw.ReadBytes(conn, 2)
	if err != nil {
		recorder.Reset()
		return nil, nil, err
	}
	_, err = recorder.Write(bs)
	if err != nil {
		recorder.Reset()
		return nil, nil, err
	}
	if !IsValidTLSVersion(bs[0], bs[1]) {
		return nil, &recorder, ErrNotTLS
	}
	bs, err = rw.ReadBytes(conn, 2)
	if err != nil {
		recorder.Reset()
		return nil, nil, err
	}
	_, err = recorder.Write(bs)
	if err != nil {
		recorder.Reset()
		return nil, nil, err
	}
	headerLen := int(binary.BigEndian.Uint16(bs))
	bs, err = rw.ReadBytes(conn, headerLen)
	if err != nil {
		recorder.Reset()
		return nil, nil, err
	}
	_, err = recorder.Write(bs)
	if err != nil {
		recorder.Reset()
		return nil, nil, err
	}
	h := &SniffHeader{}
	err = ReadClientHello(bs, h)
	if err == nil {
		return h, &recorder, nil
	}
	return nil, &recorder, err
}

func IsValidTLSVersion(major, minor byte) bool {
	if major == 3 {
		return minor < 4 && minor > 0
	}
	return false
}
