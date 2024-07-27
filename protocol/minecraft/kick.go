package minecraft

import (
	"fmt"
	"time"

	"github.com/layou233/zbproxy/v3/common/mcprotocol"
	"github.com/layou233/zbproxy/v3/config"
)

// generateKickMessage 创建一个玩家被踢下线的消息
func generateKickMessage(s *config.Outbound, name string) mcprotocol.Message {
	return mcprotocol.Message{
		Color: mcprotocol.White,
		Extra: []mcprotocol.Message{
			{Bold: true, Color: mcprotocol.Yellow, Text: "ZedWAre"},
			{Bold: true, Text: "MC Reverse Proxy"},
			{Text: " - "},
			{Bold: true, Color: mcprotocol.Gold, Text: "Connection Rejected\n"},

			{Text: "ACCESS & CONNECT DENIED by Developer\n"},
			{Text: "Reason: "},
			{Color: mcprotocol.LightPurple, Text: "It is working for the BETA User currently BUT you don't have permission to use it.\n"},
			{Text: "Please contact the Developer for help.\n\n"},

			{Text: "DEBUG INFO:\n"},
			{
				Color: mcprotocol.Gray,
				Text: fmt.Sprintf("Timestamp: %d\n", time.Now().UnixMilli()),
			},
			{
				Color: mcprotocol.Gray,
				Text: "Player Name: " + name + "\n",
			},
			{
				Color: mcprotocol.Gray,
				Text: "Service: " + s.Name + "\n",
			},
			{Text: "Developed by "},
			{
				Text: "ZedWAre, CloudDaisy, Guttridge, BarceCinear, Ren\u00e9B\u00e5\u0192",
			},
		},
	}
}

// generatePlayerNumberLimitExceededMessage 创建一个玩家人数超过限制的消息
func generatePlayerNumberLimitExceededMessage(s *config.Outbound, name string) mcprotocol.Message {
	return mcprotocol.Message{
		Color: mcprotocol.White,
		Extra: []mcprotocol.Message{
			{Text: "FULL"},
			{
				Text: fmt.Sprintf("Service: %s\n", s.Name),
			},
		},
	}
}
