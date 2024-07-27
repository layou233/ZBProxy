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
