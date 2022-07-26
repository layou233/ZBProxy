package tls

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/layou233/ZBProxy/common"
	"github.com/layou233/ZBProxy/common/set"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/outbound"
	"github.com/layou233/ZBProxy/service/access"
	"net"
)

func NewConnHandler(s *config.ConfigProxyService,
	c net.Conn,
	out outbound.Outbound) (net.Conn, error) {
	header, buf, err := SniffAndRecordTLS(c)
	defer buf.Reset()
	if err != nil {
		if err == ErrNotTLS {
			if s.TLSSniffing.RejectNonTLS {
				return nil, err
			}
			return dialAndWrite(s, buf, out)
		}
		return nil, err
	}
	domain := header.Domain()
	hit := false
	for _, list := range s.TLSSniffing.SNIAllowListTags {
		if hit = common.Must[*set.StringSet](access.GetTargetList(list)).Has(domain); hit {
			break
		}
	}
	if !hit {
		if s.TLSSniffing.RejectIfNonMatch {
			return nil, errors.New("")
		}
		return dialAndWrite(s, buf, out)
	}
	remote, err := out.Dial("tcp", fmt.Sprintf("%s:%v", domain, s.TargetPort))
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteTo(remote)
	if err != nil {
		return nil, err
	}
	return remote, nil
}

func dialAndWrite(s *config.ConfigProxyService, buffer *bytes.Buffer, out outbound.Outbound) (net.Conn, error) {
	conn, err := out.Dial("tcp", fmt.Sprintf("%s:%v", s.TargetAddress, s.TargetPort))
	if err != nil {
		return nil, err
	}
	_, err = buffer.WriteTo(conn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
