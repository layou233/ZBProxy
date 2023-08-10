package outbound

type SocketOptions struct {
	Mark          int    `json:",omitempty"`
	Interface     string `json:",omitempty"`
	TCPFastOpen   bool   `json:",omitempty"`
	TCPCongestion string `json:",omitempty"`
	MultiPathTCP  bool   `json:",omitempty"`
}
