package minecraft

import (
	"fmt"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net/packet"
	"github.com/layou233/ZBProxy/config"
	"time"
)

func generateKickMessage(s *config.ConfigProxyService, name packet.String) chat.Message {
	return chat.Message{
		Color: chat.White,
		Extra: []chat.Message{
			{Bold: true, Color: chat.Red, Text: "ZB"},
			{Bold: true, Text: "Proxy"},
			{Text: " - "},
			{Bold: true, Color: chat.Gold, Text: "Connection Rejected\n"},

			{Text: "Your connection request is refused by ZBProxy.\n"},
			{Text: "Reason: "},
			{Color: chat.LightPurple, Text: "You don't have permission to access this service.\n"},
			{Text: "Please contact the Administrators for help.\n\n"},

			{
				Color: chat.Gray,
				Text: fmt.Sprintf("Timestamp: %d | Player Name: %s | Service: %s\n",
					time.Now().UnixMilli(), name, s.Name),
			},
			{Text: "GitHub: "},
			{
				Color: chat.Aqua, UnderLined: true,
				Text:       "https://github.com/layou233/ZBProxy",
				ClickEvent: chat.OpenURL("https://github.com/layou233/ZBProxy"),
			},
		},
	}
}
