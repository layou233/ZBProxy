package minecraft

import (
	"fmt"
	"time"

	"github.com/layou233/ZBProxy/common/mcprotocol"
	"github.com/layou233/ZBProxy/config"
)

func generateKickMessage(s *config.ConfigProxyService, name string) mcprotocol.Message {
	return mcprotocol.Message{
		Color: mcprotocol.White,
		Extra: []mcprotocol.Message{
			{Bold: true, Color: mcprotocol.Red, Text: "ZB"},
			{Bold: true, Text: "Proxy"},
			{Text: " - "},
			{Bold: true, Color: mcprotocol.Gold, Text: "Connection Rejected\n"},

			{Text: "Your connection request is refused by ZBProxy.\n"},
			{Text: "Reason: "},
			{Color: mcprotocol.LightPurple, Text: "You don't have permission to access this service.\n"},
			{Text: "Please contact the Administrators for help.\n\n"},

			{
				Color: mcprotocol.Gray,
				Text: fmt.Sprintf("Timestamp: %d | Player Name: %s | Service: %s\n",
					time.Now().UnixMilli(), name, s.Name),
			},
			{Text: "GitHub: "},
			{
				Color: mcprotocol.Aqua, UnderLined: true,
				Text: "https://github.com/layou233/ZBProxy",
				//ClickEvent: chat.OpenURL("https://github.com/layou233/ZBProxy"),
			},
		},
	}
}

func generatePlayerNumberLimitExceededMessage(s *config.ConfigProxyService, name string) mcprotocol.Message {
	return mcprotocol.Message{
		Color: mcprotocol.White,
		Extra: []mcprotocol.Message{
			{Bold: true, Color: mcprotocol.Red, Text: "ZB"},
			{Bold: true, Text: "Proxy"},
			{Text: " - "},
			{Bold: true, Color: mcprotocol.Gold, Text: "Connection Rejected\n"},

			{Text: "Your connection request is refused by ZBProxy.\n"},
			{Text: "Reason: "},
			{Color: mcprotocol.LightPurple, Text: "Service online player number limitation exceeded.\n"},
			{Text: "Please contact the Administrators for help.\n\n"},

			{
				Color: mcprotocol.Gray,
				Text: fmt.Sprintf("Timestamp: %d | Player Name: %s | Service: %s\n",
					time.Now().UnixMilli(), name, s.Name),
			},
			{Text: "GitHub: "},
			{
				Color: mcprotocol.Aqua, UnderLined: true,
				Text: "https://github.com/layou233/ZBProxy",
				//ClickEvent: chat.OpenURL("https://github.com/layou233/ZBProxy"),
			},
		},
	}
}
