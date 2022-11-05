package tls

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"github.com/layou233/ZBProxy/common/rw"

	"github.com/xtls/xray-core/common/protocol/tls"
)

var ErrNotTLS = errors.New("not TLS header")

func SniffAndRecordTLS(conn net.Conn) (*tls.SniffHeader, *bytes.Buffer, error) {
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
	h := &tls.SniffHeader{}
	err = tls.ReadClientHello(bs, h)
	if err == nil {
		return h, &recorder, nil
	}
	return nil, &recorder, err
}

func IsValidTLSVersion(major, minor byte) bool {
	if major < 4 {
		if minor > 0 {
			return major > minor
		}
	}
	return false
}
