package mcprotocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/layou233/ZBProxy/common/buf"
)

// Modified from https://github.com/Tnze/go-mc/blob/master/chat/message.go
// License (only for this file): MIT

const (
	Chat = iota
	System
	GameInfo
	SayCommand
	MsgCommand
	TeamMsgCommand
	EmoteCommand
	TellrawCommand
)

// Colors
const (
	Black       = "black"
	DarkBlue    = "dark_blue"
	DarkGreen   = "dark_green"
	DarkAqua    = "dark_aqua"
	DarkRed     = "dark_red"
	DarkPurple  = "dark_purple"
	Gold        = "gold"
	Gray        = "gray"
	DarkGray    = "dark_gray"
	Blue        = "blue"
	Green       = "green"
	Aqua        = "aqua"
	Red         = "red"
	LightPurple = "light_purple"
	Yellow      = "yellow"
	White       = "white"
)

// Message is a message sent by other
type Message struct {
	Text string `json:"text"`

	Bold          bool `json:"bold,omitempty"`          // 粗体
	Italic        bool `json:"italic,omitempty"`        // 斜体
	UnderLined    bool `json:"underlined,omitempty"`    // 下划线
	StrikeThrough bool `json:"strikethrough,omitempty"` // 删除线
	Obfuscated    bool `json:"obfuscated,omitempty"`    // 随机
	// Font of the message, could be one of minecraft:uniform, minecraft:alt or minecraft:default
	// This option is only valid on 1.16+, otherwise the property is ignored.
	Font  string `json:"font,omitempty"`  // 字体
	Color string `json:"color,omitempty"` // 颜色

	// Insertion contains text to insert. Only used for messages in chat.
	// When shift is held, clicking the component inserts the given text
	// into the chat box at the cursor (potentially replacing selected text).
	Insertion string `json:"insertion,omitempty"`

	Translate string    `json:"translate,omitempty"`
	With      []Message `json:"with,omitempty"`
	Extra     []Message `json:"extra,omitempty"`
}

// Same as Message, but "Text" is omitempty
type translateMsg struct {
	Text string `json:"text,omitempty"`

	Bold          bool `json:"bold,omitempty"`
	Italic        bool `json:"italic,omitempty"`
	UnderLined    bool `json:"underlined,omitempty"`
	StrikeThrough bool `json:"strikethrough,omitempty"`
	Obfuscated    bool `json:"obfuscated,omitempty"`

	Font  string `json:"font,omitempty"`
	Color string `json:"color,omitempty"`

	Insertion string `json:"insertion,omitempty"`

	Translate string    `json:"translate"`
	With      []Message `json:"with,omitempty"`
	Extra     []Message `json:"extra,omitempty"`
}

type jsonMsg Message

func (m Message) MarshalJSON() ([]byte, error) {
	if m.Translate != "" {
		return json.Marshal(translateMsg(m))
	} else {
		return json.Marshal(jsonMsg(m))
	}
}

// UnmarshalJSON decode json to Message
func (m *Message) UnmarshalJSON(raw []byte) (err error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return io.EOF
	}
	// The right way to distinguish JSON String and Object
	// is to look up the first character.
	switch raw[0] {
	case '"':
		return json.Unmarshal(raw, &m.Text) // Unmarshal as jsonString
	case '{':
		return json.Unmarshal(raw, (*jsonMsg)(m)) // Unmarshal as jsonMsg
	case '[':
		return json.Unmarshal(raw, &m.Extra) // Unmarshal as []Message
	default:
		return errors.New("unknown chat message type: '" + string(raw[0]) + "'")
	}
}

func (m *Message) ReadMessage(buffer *buf.Buffer) error {
	length, _, err := ReadVarIntFrom(buffer)
	if err != nil {
		return err
	}
	code, err := buffer.Peek(int(length))
	if err != nil {
		return err
	}
	err = json.Unmarshal(code, m)
	return err
}

// WriteTo encode Message into a ChatMsg packet
func (m Message) WriteTo(w io.Writer) (int64, error) {
	code, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	_, err = VarInt(int32(len(code))).WriteTo(w)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(code)
	return int64(n), err
}
