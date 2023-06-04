package minecraft

import (
	"encoding/json"

	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/service/transfer"
	"github.com/layou233/ZBProxy/version"
)

type sample struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type motdObject struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int      `json:"max"`
		Online int      `json:"online"`
		Sample []sample `json:"sample"`
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

	var samples []sample
	if len(s.Minecraft.OnlineCount.Sample) > 0 {
		samples = make([]sample, 0, len(s.Minecraft.OnlineCount.Sample))
		for id, name := range s.Minecraft.OnlineCount.Sample {
			samples = append(samples, sample{ID: id, Name: name})
		}
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
			Max    int      `json:"max"`
			Online int      `json:"online"`
			Sample []sample `json:"sample"`
		}{
			Max:    s.Minecraft.OnlineCount.Max,
			Online: int(online),
			Sample: samples,
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
