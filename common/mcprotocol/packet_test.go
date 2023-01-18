package mcprotocol

/*
import (
	"io"
	"testing"

	"github.com/layou233/ZBProxy/common/buf"

	mcnet "github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/packet"
)

func BenchmarkZBPacketWrite(b *testing.B) {
	buffer := buf.NewSize(512)
	buffer.Reset(MaxVarIntLen)
	mcConn := StreamConn(nil)
	mcConn.Writer = io.Discard
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WriteToPacket(buffer,
			byte(0x00),
			VarInt(47),
			"mc.hypixel.net",
			uint16(25565),
			byte(0))
		mcConn.WritePacket(buffer)
	}
}

func BenchmarkGoMCPacketWrite(b *testing.B) {
	mcConn := mcnet.WrapConn(nil)
	mcConn.Writer = io.Discard
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mcConn.WritePacket(packet.Marshal(0x00,
			packet.VarInt(47),
			packet.String("mc.hypixel.net"),
			packet.UnsignedShort(25565),
			packet.Byte(0)))
	}
}

func BenchmarkPacket(b *testing.B) {
	b.Run("BenchmarkZBPacketWrite", BenchmarkZBPacketWrite)
	b.Run("BenchmarkGoMCPacketWrite", BenchmarkGoMCPacketWrite)
}
*/
