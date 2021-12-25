package minecraft

import (
	"ZBProxy/version"
	"encoding/json"
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

type motdObject struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		/*		Sample []struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			} `json:"sample"`*/
	} `json:"players"`
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	Favicon string `json:"favicon"`
}

func generateMotdPacket(protocolVersion int, motdFavicon, motdDescription string) packet.Packet {
	motd, _ := json.Marshal(motdObject{
		Version: struct {
			Name     string `json:"name"`
			Protocol int    `json:"protocol"`
		}{
			Name:     "ZBProxy " + version.Version,
			Protocol: protocolVersion,
		},
		Players: struct {
			Max    int `json:"max"`
			Online int `json:"online"`
		}{ // TODO Show online players in server list.
			Max:    1,
			Online: 0,
		},
		Description: struct {
			Text string `json:"text"`
		}{
			Text: motdDescription,
		},
		Favicon: motdFavicon,
	})
	return packet.Marshal(packetid.ServerInfo, packet.String(motd))
}
