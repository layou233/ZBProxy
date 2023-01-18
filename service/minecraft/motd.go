package minecraft

import (
	"encoding/json"

	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/service/transfer"
	"github.com/layou233/ZBProxy/version"
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

func generateMOTD(protocolVersion int, s *config.ConfigProxyService, options *transfer.Options) []byte {
	online := s.Minecraft.OnlineCount.Online
	if online < 0 {
		online = options.OnlineCount.Load()
	}
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
		}{
			Max:    s.Minecraft.OnlineCount.Max,
			Online: int(online),
		},
		Description: struct {
			Text string `json:"text"`
		}{
			Text: s.Minecraft.MotdDescription,
		},
		Favicon: s.Minecraft.MotdFavicon,
	})
	return motd
}
