package socks

const (
	version4 byte = 4

	ReplyCode4Granted                     byte = 0x5A
	ReplyCode4RejectedOrFailed            byte = 0x5B
	ReplyCode4CannotConnectToIdentd       byte = 0x5C
	ReplyCode4IdentdReportDifferentUserID byte = 0x5D
)
