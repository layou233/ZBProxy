package minecraft

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/access"
	"github.com/layou233/zbproxy/v3/common/buf"
	"github.com/layou233/zbproxy/v3/common/bufio"
	"github.com/layou233/zbproxy/v3/common/mcprotocol"
	"github.com/layou233/zbproxy/v3/common/network"
	"github.com/layou233/zbproxy/v3/common/network/socks"
	"github.com/layou233/zbproxy/v3/common/proxyprotocol"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"
	"github.com/layou233/zbproxy/v3/version"

	"github.com/phuslu/log"
	"github.com/zhangyunhao116/fastrand"
)

var minecraftSRV = &adapter.SRVMetadata{ServiceName: "minecraft"}

type Outbound struct {
	logger *log.Logger
	config *config.Outbound
	router adapter.Router
	dialer network.Dialer

	hostnameAccessLists []set.StringSet
	nameAccessLists     []set.StringSet
	onlineCount         atomic.Int32
}

var (
	_ adapter.Outbound = (*Outbound)(nil)
	_ network.Dialer   = (*Outbound)(nil)
)

func NewOutbound(logger *log.Logger, newConfig *config.Outbound) (*Outbound, error) {
	if newConfig.Minecraft == nil {
		return nil, errors.New("not Minecraft outbound config")
	}
	outbound := &Outbound{
		logger: logger,
		config: newConfig,
	}
	return outbound, nil
}

func (o *Outbound) Name() string {
	if o.config != nil {
		return o.config.Name
	}
	return ""
}

func (o *Outbound) PostInitialize(router adapter.Router, provider adapter.RouteResourceProvider) error {
	var err error
	if o.config.Minecraft.HostnameAccess.Mode != access.DefaultMode {
		o.hostnameAccessLists, err = provider.FindListsByTag(o.config.Minecraft.HostnameAccess.ListTags)
		if err != nil {
			return common.Cause("load access control lists: ", err)
		}
	}
	if o.config.Minecraft.NameAccess.Mode != access.DefaultMode {
		o.nameAccessLists, err = provider.FindListsByTag(o.config.Minecraft.NameAccess.ListTags)
		if err != nil {
			return common.Cause("load access control lists: ", err)
		}
	}
	if o.config.Minecraft.MotdFavicon == "{DEFAULT_MOTD}" {
		o.config.Minecraft.MotdFavicon = defaultMOTD
	}
	o.config.Minecraft.MotdDescription = strings.NewReplacer(
		"{INFO}", "ZBProxy "+version.Version,
		"{NAME}", o.config.Name,
		"{HOST}", o.config.TargetAddress,
		"{PORT}", strconv.Itoa(int(o.config.TargetPort)),
	).Replace(o.config.Minecraft.MotdDescription)

	if samples := o.config.Minecraft.OnlineCount.Sample; samples != nil {
		var convertedSamples []playerSample
		switch samples := samples.(type) {
		case map[string]any:
			convertedSamples = make([]playerSample, 0, len(samples))
			for uuid, name := range samples {
				convertedSamples = append(convertedSamples, playerSample{
					Name: name.(string),
					ID:   uuid,
				})
			}

		case []any:
			convertedSamples = make([]playerSample, 0, len(samples))
			var u [16]byte
			var dst [36]byte
			for i, sample := range samples {
				// generate random UUID with ZBProxy signature
				fastrand.Read(u[:])
				u[0] = byte(i)
				u[1] = '$'
				u[2] = 'Z'
				u[3] = 'B'
				u[4] = '$'

				// marshal UUID string
				const hexTable = "0123456789abcdef"
				dst[8] = '-'
				dst[13] = '-'
				dst[18] = '-'
				dst[23] = '-'
				for i, x := range [16]byte{
					0, 2, 4, 6,
					9, 11,
					14, 16,
					19, 21,
					24, 26, 28, 30, 32, 34,
				} {
					c := u[i]
					dst[x] = hexTable[c>>4]
					dst[x+1] = hexTable[c&0x0F]
				}

				convertedSamples = append(convertedSamples, playerSample{
					Name: sample.(string),
					ID:   string(dst[:]),
				})
			}

		default:
			return fmt.Errorf("unknown player samples type: %T", samples)
		}
		o.config.Minecraft.OnlineCount.Sample = convertedSamples
	}

	if o.config.Dialer != "" {
		if o.config.SocketOptions != nil {
			return errors.New("socket options are not available when dialer is specified")
		}
		o.dialer, err = provider.FindOutboundByName(o.config.Dialer)
		if err != nil {
			return err
		}
	} else {
		o.dialer = network.NewSystemDialer(o.config.SocketOptions)
	}
	switch o.config.ProxyProtocolVersion {
	case proxyprotocol.VersionUnspecified,
		proxyprotocol.Version1,
		proxyprotocol.Version2:
	default:
		return fmt.Errorf("invalid proxy protocol version: %v", o.config.ProxyProtocolVersion)
	}
	switch o.config.ProxyOptions.Type {
	case "socks", "socks5", "socks4a", "socks4":
		o.dialer = &socks.Client{
			Dialer:  o.dialer,
			Version: o.config.ProxyOptions.Type,
			Network: o.config.ProxyOptions.Network,
			Address: o.config.ProxyOptions.Address,
		}
	}
	o.router = router
	return nil
}

func (o *Outbound) Reload(options adapter.OutboundReloadOptions) error {
	o.config = options.Config
	o.hostnameAccessLists = nil
	o.nameAccessLists = nil
	return o.PostInitialize(o.router, &options)
}

func (o *Outbound) connectServer(ctx context.Context, metadata *adapter.Metadata) (net.Conn, error) {
	if metadata.DestinationHostname == "" {
		metadata.DestinationHostname = o.config.TargetAddress
	}
	if metadata.DestinationPort == 0 {
		metadata.DestinationPort = o.config.TargetPort
	}
	destinationAddress := net.JoinHostPort(metadata.DestinationHostname, strconv.FormatUint(uint64(metadata.DestinationPort), 10))
	if !o.config.Minecraft.IgnoreSRVRedirect {
		metadata.SRV = minecraftSRV
	}
	conn, err := adapter.DialContextWithMetadata(o.dialer, ctx, "tcp", destinationAddress, metadata)
	if err != nil {
		return nil, err
	}
	if o.config.ProxyProtocolVersion != proxyprotocol.VersionUnspecified {
		var localAddress netip.AddrPort
		localAddress, err = netip.ParseAddrPort(conn.LocalAddr().String())
		if err != nil {
			conn.Close()
			return nil, common.Cause("failed to parse local address: ", err)
		}
		err = (&proxyprotocol.Header{
			Version:           uint8(o.config.ProxyProtocolVersion),
			Command:           proxyprotocol.CommandProxy,
			TransportProtocol: proxyprotocol.TransportProtocolStream | proxyprotocol.AddressFamilyByAddr(metadata.SourceAddress.Addr()),
			SourceAddress:     metadata.SourceAddress,
		}).WriteHeader(conn, localAddress)
		if err != nil {
			conn.Close()
			return nil, common.Cause("failed to write PROXY protocol header: ", err)
		}
	}
	return conn, nil
}

func (o *Outbound) InjectConnection(ctx context.Context, conn *bufio.CachedConn, metadata *adapter.Metadata) error {
	if metadata.Minecraft == nil {
		return errors.New("require Minecraft metadata")
	}
	if o.config.Minecraft.HostnameAccess.Mode != access.DefaultMode {
		hostnameClean := metadata.Minecraft.CleanOriginDestination()
		if o.config.Minecraft.HostnameAccess.LowerCase {
			hostnameClean = strings.ToLower(hostnameClean)
		}
		if !access.Check(o.hostnameAccessLists, o.config.Minecraft.HostnameAccess.Mode, hostnameClean) {
			conn.Conn.(*net.TCPConn).SetLinger(0)
			conn.Close()
			return common.Cause("hostname "+o.config.Minecraft.HostnameAccess.Mode+
				" mode, request="+url.QueryEscape(hostnameClean)+": ", access.ErrRejected)
		}
	}
	if metadata.Minecraft.SniffPosition >= 0 {
		conn.Rewind(metadata.Minecraft.SniffPosition)
	}
	switch metadata.Minecraft.NextState {
	case mcprotocol.NextStateStatus:
		// skip Status Request packet
		_, err := conn.Peek(2)
		if err != nil {
			return common.Cause("skip status request: ", err)
		}
		if o.config.Minecraft.MotdFavicon == "" && o.config.Minecraft.MotdDescription == "" {
			// directly proxy MOTD from server
			var remoteConn net.Conn
			remoteConn, err = o.connectServer(ctx, metadata)
			if err != nil {
				return common.Cause("request remote MOTD: ", err)
			}
			//remoteConn.(*net.TCPConn).SetLinger(0) // for some reason
			if metadata.Minecraft.RewrittenDestination == "" {
				metadata.Minecraft.RewrittenDestination = metadata.Minecraft.CleanOriginDestination()
			}
			if metadata.Minecraft.RewrittenPort == 0 {
				metadata.Minecraft.RewrittenPort = metadata.Minecraft.OriginPort
			}
			buffer := buf.New()
			buffer.Reset(mcprotocol.MaxVarIntLen)

			hostname := metadata.Minecraft.RewrittenDestination
			if o.config.Minecraft.EnableHostnameRewrite {
				hostname = o.config.Minecraft.RewrittenHostname
				if hostname == "" {
					hostname = o.config.TargetAddress
				}
			} else if hostname == "" {
				hostname = metadata.Minecraft.CleanOriginDestination()
			}
			if !o.config.Minecraft.IgnoreFMLSuffix && metadata.Minecraft.IsFML() {
				hostname += "\x00" + metadata.Minecraft.FMLMarkup()
			}
			port := metadata.Minecraft.RewrittenPort
			if port <= 0 {
				port = metadata.Minecraft.OriginPort
			}
			// construct handshake packet
			buffer.WriteByte(0) // Server bound : Handshake
			mcprotocol.VarInt(metadata.Minecraft.ProtocolVersion).WriteToBuffer(buffer)
			mcprotocol.WriteString(buffer, hostname)
			binary.BigEndian.PutUint16(buffer.Extend(2), port)
			buffer.WriteByte(mcprotocol.NextStateStatus)
			mcprotocol.AppendPacketLength(buffer, buffer.Len())
			// construct status packet
			buffer.WriteByte(1)
			buffer.WriteByte(0)
			// send 2 packets in 1 write call
			_, err = remoteConn.Write(buffer.Bytes())
			buffer.Release()
			if err != nil {
				return common.Cause("request remote MOTD: ", err)
			}
			return bufio.CopyConn(remoteConn, conn)
		} else {
			motd := generateMOTD(metadata.Minecraft.ProtocolVersion, o.config, &o.onlineCount)
			buffer := buf.New()
			buffer.Reset(mcprotocol.MaxVarIntLen)
			buffer.WriteByte(0) // Client bound : Status Response
			mcprotocol.VarInt(len(motd)).WriteToBuffer(buffer)
			clientMC := mcprotocol.Conn{
				Reader: conn,
				Writer: common.UnwrapWriter(conn), // unwrap to make writev syscall possible
				Conn:   conn,
			}
			err = clientMC.WriteVectorizedPacket(buffer, motd)
			if err != nil {
				return common.Cause("respond MOTD: ", err)
			}

			switch o.config.Minecraft.PingMode {
			case pingModeDisconnect:
				// do nothing and disconnect
			case pingMode0ms:
				buffer.WriteByte(1)  // Client bound : Ping Response
				buffer.WriteZeroN(8) // size of int64 timestamp
				err = clientMC.WritePacket(buffer)
				buffer.Release()
				if err != nil {
					return common.Cause("respond 0ms ping: ", err)
				}
			default:
				err = clientMC.ReadLimitedPacket(buffer, 9)
				if err != nil {
					buffer.Release()
					return common.Cause("read ping request: ", err)
				}
				err = clientMC.WritePacket(buffer)
				buffer.Release()
				if err != nil {
					return common.Cause("respond ping request: ", err)
				}
			}
			o.logger.Info().Str("id", metadata.ConnectionID).Str("outbound", o.config.Name).Msg("Responded MOTD")
			return nil
		}

	case mcprotocol.NextStateLogin:
		buffer := buf.New()
		buffer.Reset(mcprotocol.MaxVarIntLen)
		if o.config.Minecraft.NameAccess.Mode != access.DefaultMode {
			name := metadata.Minecraft.PlayerName
			if o.config.Minecraft.NameAccess.LowerCase {
				name = strings.ToLower(metadata.Minecraft.PlayerName)
			}
			if !access.Check(o.nameAccessLists, o.config.Minecraft.NameAccess.Mode, name) {
				msg, err := generateKickMessage(o.config, metadata.Minecraft.PlayerName).MarshalJSON()
				if err != nil { // almost impossible
					buffer.Release()
					return common.Cause("generate kick message: ", err)
				}
				buffer.WriteByte(0) // Client bound : Disconnect (login)
				mcprotocol.VarInt(len(msg)).WriteToBuffer(buffer)
				err = mcprotocol.Conn{Writer: common.UnwrapWriter(conn)}.WriteVectorizedPacket(buffer, msg)
				if err != nil {
					buffer.Release()
					return common.Cause("send kick packet: ", err)
				}
				o.logger.Warn().Str("id", metadata.ConnectionID).Str("outbound", o.config.Name).
					Str("player", metadata.Minecraft.PlayerName).Msg("Kicked by name access control")
				conn.Conn.(*net.TCPConn).SetLinger(10)
				buffer.Release()
				return nil
			}
		}
		if o.config.Minecraft.OnlineCount.EnableMaxLimit &&
			o.config.Minecraft.OnlineCount.Max <= o.onlineCount.Load() {
			msg, err := generatePlayerNumberLimitExceededMessage(o.config, metadata.Minecraft.PlayerName).MarshalJSON()
			if err != nil {
				buffer.Release()
				return common.Cause("generate player number limit exceeded packet: ", err)
			}
			buffer.WriteByte(0)
			mcprotocol.VarInt(len(msg)).WriteToBuffer(buffer)
			err = mcprotocol.Conn{Writer: common.UnwrapWriter(conn)}.WriteVectorizedPacket(buffer, msg)
			if err != nil {
				buffer.Release()
				return common.Cause("send player number limit exceeded packet: ", err)
			}
			o.logger.Warn().Str("id", metadata.ConnectionID).Str("outbound", o.config.Name).
				Str("player", metadata.Minecraft.PlayerName).Msg("Kicked by player number limiter")
			conn.Conn.(*net.TCPConn).SetLinger(10)
			buffer.Release()
			return nil
		}

		serverConn, err := o.connectServer(ctx, metadata)
		if err != nil {
			buffer.Release()
			return common.Cause("connect server: ", err)
		}
		hostname := metadata.Minecraft.RewrittenDestination
		if o.config.Minecraft.EnableHostnameRewrite {
			hostname = o.config.Minecraft.RewrittenHostname
			if hostname == "" {
				hostname = o.config.TargetAddress
			}
		} else if hostname == "" {
			hostname = metadata.Minecraft.CleanOriginDestination()
		}
		if !o.config.Minecraft.IgnoreFMLSuffix && metadata.Minecraft.IsFML() {
			hostname += "\x00" + metadata.Minecraft.FMLMarkup()
		}
		port := metadata.Minecraft.RewrittenPort
		if port <= 0 {
			port = metadata.Minecraft.OriginPort
		}
		// construct handshake packet
		buffer.WriteByte(0) // Server bound : Handshake
		mcprotocol.VarInt(metadata.Minecraft.ProtocolVersion).WriteToBuffer(buffer)
		mcprotocol.WriteString(buffer, hostname)
		binary.BigEndian.PutUint16(buffer.Extend(2), port)
		buffer.WriteByte(mcprotocol.NextStateLogin)
		mcprotocol.AppendPacketLength(buffer, buffer.Len())
		// write handshake and login packet
		cache := conn.Cache()
		vector := net.Buffers{buffer.Bytes(), cache.Bytes()}
		_, err = vector.WriteTo(serverConn)
		buffer.Release()
		if err != nil {
			serverConn.Close()
			return common.Cause("server handshake: ", err)
		}
		cache.Advance(cache.Len()) // all written
		o.logger.Info().Str("id", metadata.ConnectionID).Str("outbound", o.config.Name).
			Str("player", metadata.Minecraft.PlayerName).Msg("Created Minecraft connection")
		o.onlineCount.Add(1)
		err = bufio.CopyConn(serverConn, conn)
		o.onlineCount.Add(-1)
		return err

	case mcprotocol.NextStateTransfer:
		// TODO: Minecraft transfer support
		conn.Conn.(*net.TCPConn).SetLinger(0)
		return conn.Close()

	default:
		return errors.New("unknown next state")
	}
}

func (o *Outbound) DialContext(context.Context, string, string) (net.Conn, error) {
	return nil, adapter.ErrInjectionRequired
}
