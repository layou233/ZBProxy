package minecraft

import (
	"errors"
	"io"
	"time"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/buf"
	"github.com/layou233/zbproxy/v3/common/bufio"
	"github.com/layou233/zbproxy/v3/common/mcprotocol"
)

var ErrBadPacket = errors.New("bad Minecraft handshake packet")

func SniffClientHandshake(conn bufio.PeekConn, metadata *adapter.Metadata) error {
	if metadata.Minecraft == nil {
		metadata.Minecraft = &adapter.MinecraftMetadata{
			NextState: -1,
		}
	}
	defer conn.SetReadDeadline(time.Time{}) // clear deadline

	// handshake packet
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	packetSize, _, err := mcprotocol.ReadVarIntFrom(conn)
	if err != nil {
		return common.Cause("read packet size: ", err)
	}
	if packetSize > 264 { // maximum possible size of this kind of packet
		return ErrBadPacket
	}
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	packetContent, err := conn.Peek(int(packetSize))
	if err != nil {
		return common.Cause("read handshake packet: ", err)
	}
	buffer := buf.As(packetContent)

	var packetID byte
	packetID, err = buffer.ReadByte()
	if err != nil {
		return common.Cause("read packet ID: ", err)
	}
	if packetID != 0 { // Server bound : Handshake
		return ErrBadPacket
	}
	protocolVersion, _, err := mcprotocol.ReadVarIntFrom(buffer)
	if err != nil {
		return common.Cause("read protocol version: ", err)
	}
	if protocolVersion <= 0 {
		return ErrBadPacket
	}
	metadata.Minecraft.ProtocolVersion = uint(protocolVersion)

	metadata.Minecraft.OriginDestination, err = mcprotocol.ReadString(buffer)
	if err != nil {
		return common.Cause("read destination: ", err)
	}
	if metadata.Minecraft.OriginDestination == "" {
		return ErrBadPacket
	}

	metadata.Minecraft.OriginPort, err = mcprotocol.ReadUint16(buffer)
	if err != nil {
		return common.Cause("read port: ", err)
	}
	if metadata.Minecraft.OriginPort == 0 {
		return ErrBadPacket
	}

	intent, err := buffer.ReadByte()
	if err != nil {
		return common.Cause("read next state: ", err)
	}
	switch intent {
	case mcprotocol.IntentLogin,
		mcprotocol.IntentStatus,
		mcprotocol.IntentTransfer:
	default:
		return ErrBadPacket
	}
	metadata.Minecraft.NextState = int8(intent)

	metadata.Minecraft.SniffPosition = conn.CurrentPosition()
	if intent == mcprotocol.IntentStatus {
		// status packet
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, err = conn.Peek(2)
		if err != nil {
			return common.Cause("read status request: ", err)
		}
	} else {
		// login packet or transfer packet
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		packetSize, _, err = mcprotocol.ReadVarIntFrom(conn)
		if err != nil {
			return common.Cause("read packet size: ", err)
		}
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		packetContent, err = conn.Peek(int(packetSize))
		if err != nil {
			return common.Cause("read login packet: ", err)
		}
		buffer = buf.As(packetContent)

		packetID, err = buffer.ReadByte()
		if err != nil {
			return common.Cause("read packet ID: ", err)
		}
		if packetID != 0 { // Server bound : Login Start
			return ErrBadPacket
		}
		metadata.Minecraft.PlayerName, err = mcprotocol.ReadLimitedString(buffer, 16)
		if err != nil {
			return common.Cause("read player name: ", err)
		}
		if metadata.Minecraft.ProtocolVersion >= 764 { // 1.20.2
			if buffer.Len() == 16 { // UUID exists
				copy(metadata.Minecraft.UUID[:], buffer.Bytes())
			}
		} else if metadata.Minecraft.ProtocolVersion >= 761 { // 1.19.3
			var hasUUID byte
			hasUUID, err = buffer.ReadByte()
			if err != nil {
				return common.Cause("read has UUID: ", err)
			}
			if hasUUID == mcprotocol.BooleanTrue {
				copy(metadata.Minecraft.UUID[:], buffer.Bytes())
			}
		} else if metadata.Minecraft.ProtocolVersion >= 759 { // 1.19
			var hasSigData byte
			hasSigData, err = buffer.ReadByte()
			if err != nil {
				return common.Cause("read has sig data: ", err)
			}
			if hasSigData == mcprotocol.BooleanTrue {
				// skip timestamp
				buffer.Advance(8) // size of Long
				var length int32
				// skip public key
				length, _, err = mcprotocol.ReadVarIntFrom(buffer)
				if err != nil {
					return common.Cause("read public key length: ", err)
				}
				buffer.Advance(int(length))
				// skip signature
				length, _, err = mcprotocol.ReadVarIntFrom(buffer)
				if err != nil {
					return common.Cause("read signature length: ", err)
				}
				buffer.Advance(int(length))
			}
			var hasUUID byte
			hasUUID, err = buffer.ReadByte()
			if err != nil && err != io.EOF {
				return common.Cause("read has UUID: ", err)
			}
			if hasUUID == mcprotocol.BooleanTrue {
				copy(metadata.Minecraft.UUID[:], buffer.Bytes())
			}
		}
	}

	return nil
}
