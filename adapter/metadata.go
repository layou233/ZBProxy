package adapter

import (
	"net/netip"
	"strconv"
	"strings"

	"github.com/layou233/zbproxy/v3/common/console/color"

	"github.com/zhangyunhao116/fastrand"
)

type Protocol = uint8

/*var (
	greenPlus = color.Apply(color.FgHiGreen, "[+]")
	redMinus  = color.Apply(color.FgHiRed, "[-]")
)*/

type Metadata struct {
	ConnectionID        string
	ServiceName         string
	SniffedProtocol     Protocol
	SourceAddress       netip.AddrPort
	DestinationHostname string
	DestinationPort     uint16
	SRV                 *SRVMetadata
	Minecraft           *MinecraftMetadata
	TLS                 *TLSMetadata
	Custom              map[string]any
}

func (m *Metadata) GenerateID() {
	id := int64(fastrand.Int31())
	idColor := fastrand.Int31n(int32(len(color.List)))
	m.ConnectionID = color.Apply(color.List[idColor], "["+strconv.FormatInt(id, 10)+"]")
}

type MinecraftMetadata struct {
	ProtocolVersion      uint
	PlayerName           string
	OriginDestination    string
	RewrittenDestination string
	fmlMarkup            string
	OriginPort           uint16
	RewrittenPort        uint16
	UUID                 [16]byte
	NextState            int8
	SniffPosition        int
}

func (m *MinecraftMetadata) Valid() bool {
	return 0 < m.NextState
}

func (m *MinecraftMetadata) IsFML() bool {
	return strings.IndexByte(m.OriginDestination, 0) != -1
}

func (m *MinecraftMetadata) CleanOriginDestination() (clean string) {
	clean, m.fmlMarkup, _ = strings.Cut(m.OriginDestination, "\x00")
	return
}

func (m *MinecraftMetadata) FMLMarkup() string {
	if m.fmlMarkup == "" {
		_, m.fmlMarkup, _ = strings.Cut(m.OriginDestination, "\x00")
	}
	return m.fmlMarkup
}

type TLSMetadata struct {
	SNI string
}

type SRVMetadata struct {
	ServiceName string
}
