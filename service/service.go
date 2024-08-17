package service

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"os"
	"strconv"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/access"
	"github.com/layou233/zbproxy/v3/common/bufio"
	"github.com/layou233/zbproxy/v3/common/network"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"
	"github.com/layou233/zbproxy/v3/protocol/minecraft"

	"github.com/phuslu/log"
)

type Service struct {
	tcpListener    *net.TCPListener
	ctx            context.Context
	router         adapter.Router
	logger         *log.Logger
	config         *config.Service
	legacyOutbound adapter.Outbound
	listenAddress  string
	ipAccessLists  []set.StringSet

	// TODO: udp service
}

var _ adapter.Service = (*Service)(nil)

func NewService(logger *log.Logger, newConfig *config.Service) *Service {
	return &Service{
		listenAddress: ":" + strconv.Itoa(int(newConfig.Listen)),
		logger:        logger,
		config:        newConfig,
	}
}

func (s *Service) listenLoop() {
	for {
		conn, err := s.tcpListener.AcceptTCP()
		if err != nil {
			return
		}
		go func() {
			tcpAddress := conn.RemoteAddr().(*net.TCPAddr)
			ipString := tcpAddress.IP.String()
			if s.ipAccessLists != nil &&
				!access.Check(s.ipAccessLists, s.config.IPAccess.Mode, ipString) {
				conn.SetLinger(0)
				conn.Close()
				s.logger.Warn().Str("service", s.config.Name).Str("ip", ipString).Msg("Rejected by access control")
				return
			}
			metadata := &adapter.Metadata{
				ServiceName:         s.config.Name,
				DestinationHostname: s.config.TargetAddress,
				DestinationPort:     s.config.TargetPort,
				SourceAddress:       netip.AddrPortFrom(common.MustOK(netip.AddrFromSlice(tcpAddress.IP)).Unmap(), uint16(tcpAddress.Port)),
			}
			metadata.GenerateID()
			s.logger.Info().Str("id", metadata.ConnectionID).Str("service", s.config.Name).
				Str("ip", ipString).Msg("New inbound connection")
			if s.legacyOutbound != nil {
				defer s.logger.Info().Str("id", metadata.ConnectionID).Str("service", s.config.Name).
					Str("ip", ipString).Msg("Disconnected")
				defer conn.Close()
				switch outbound := s.legacyOutbound.(type) {
				case *minecraft.Outbound:
					bufConn := &bufio.CachedConn{Conn: conn}
					err = minecraft.SniffClientHandshake(bufConn, metadata)
					bufConn.Release()
					if err != nil {
						s.logger.Warn().Str("id", metadata.ConnectionID).Str("service", s.config.Name).
							Str("ip", ipString).Err(err).Msg("Error when reading Minecraft handshake")
						return
					}
					err = outbound.InjectConnection(s.ctx, bufConn, metadata)
					if err != nil {
						s.logger.Info().Str("id", metadata.ConnectionID).Str("service", s.config.Name).
							Str("player", metadata.Minecraft.PlayerName).Err(err).Msg("Handling Minecraft connection")
					}
				}
			} else {
				s.router.HandleConnection(conn, metadata)
			}
		}()
	}
}

func (s *Service) Start(ctx context.Context) error {
	var err error
	// handle legacy modes
	if s.config.Minecraft != nil && s.config.TLSSniffing != nil {
		return errors.New("Minecraft and TLSSniffing are mutually exclusive in legacy mode")
	}
	if s.config.Minecraft != nil {
		s.legacyOutbound, err = minecraft.NewOutbound(s.logger, &config.Outbound{
			Name:          "legacy-" + s.config.Name,
			TargetAddress: s.config.TargetAddress,
			TargetPort:    s.config.TargetPort,
			Minecraft:     s.config.Minecraft,
			SocketOptions: network.ConvertLegacyOutboundOptions(s.config.SocketOptions),
		})
		if err != nil {
			return common.Cause("initialize legacy Minecraft outbound: ", err)
		}
		err = s.legacyOutbound.PostInitialize(s.router)
		if err != nil {
			return common.Cause("post initialize legacy Minecraft outbound: ", err)
		}
	}

	// load legacy IP access control
	if s.config.IPAccess.Mode != access.DefaultMode {
		s.ipAccessLists, err = s.router.FindListsByTag(s.config.IPAccess.ListTags)
		if err != nil {
			return common.Cause("load access control lists: ", err)
		}
	}

	listenConfig := &net.ListenConfig{
		Control: network.NewListenerControlFromOptions(s.config.SocketOptions),
	}
	if s.config.SocketOptions != nil {
		network.SetListenerTCPKeepAlive(listenConfig, s.config.SocketOptions.KeepAliveConfig())
		if s.config.SocketOptions.MultiPathTCP {
			network.SetListenerMultiPathTCP(listenConfig, true)
		}
	}
	listener, err := listenConfig.Listen(ctx, "tcp", s.listenAddress)
	if err != nil {
		return common.Cause("start listening: ", err)
	}
	s.tcpListener = listener.(*net.TCPListener)
	s.ctx = ctx
	s.logger.Info().Str("service", s.config.Name).Msg("Listening on " + s.listenAddress)

	go s.listenLoop()
	return nil
}

func (s *Service) Reload(ctx context.Context, newConfig *config.Service) error {
	if s.tcpListener == nil {
		return os.ErrClosed
	}
	s.Close()
	s.listenAddress = ":" + strconv.Itoa(int(newConfig.Listen))
	s.config = newConfig
	s.legacyOutbound = nil
	s.ipAccessLists = nil
	return s.Start(ctx)
}

func (s *Service) UpdateRouter(router adapter.Router) {
	s.router = router
}

func (s *Service) Close() error {
	if s.tcpListener == nil {
		return os.ErrClosed
	}
	err := s.tcpListener.Close()
	s.tcpListener = nil
	return err
}
