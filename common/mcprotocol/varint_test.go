package mcprotocol

import (
	"bytes"
	"io"
	"testing"

	"github.com/layou233/ZBProxy/common/buf"
)

// samples from https://wiki.vg/Protocol#VarInt_and_VarLong

func checkWrite(t *testing.T, n int32, result [MaxVarIntLen]byte) {
	buffer := buf.NewSize(MaxVarIntLen + 1)
	defer buffer.Release()
	_, err := VarInt(n).WriteTo(buffer)
	if err != nil {
		return
	}
	t.Log("VarInt", n, "WriteTo", buffer.Bytes())
	buffer.Truncate(MaxVarIntLen)
	if !bytes.Equal(buffer.Bytes(), result[:]) {
		t.Fatalf("VarInt WriteTo error: got %v, expect %v", buffer.Bytes(), result)
	}
}

func checkRead(t *testing.T, n int32, result [MaxVarIntLen]byte) {
	buffer := buf.As(result[:])
	defer buffer.Release()
	vi, _, err := ReadVarIntFrom(buffer)
	if err != nil {
		return
	}
	t.Log("VarInt", vi, "ReadFrom", result)
	if n != vi {
		t.Fatalf("VarInt ReadFrom error: got %v, expect %v", vi, n)
	}
}

func TestVarInt_WriteTo(t *testing.T) {
	checkWrite(t, 0, [MaxVarIntLen]byte{0})
	checkWrite(t, 1, [MaxVarIntLen]byte{1})
	checkWrite(t, 2, [MaxVarIntLen]byte{2})
	checkWrite(t, 127, [MaxVarIntLen]byte{127})
	checkWrite(t, 128, [MaxVarIntLen]byte{128, 1})
	checkWrite(t, 255, [MaxVarIntLen]byte{255, 1})
	checkWrite(t, 25565, [MaxVarIntLen]byte{221, 199, 1})
	checkWrite(t, 2097151, [MaxVarIntLen]byte{255, 255, 127})
	checkWrite(t, 2147483647, [MaxVarIntLen]byte{255, 255, 255, 255, 7})
	checkWrite(t, -1, [MaxVarIntLen]byte{255, 255, 255, 255, 15})
	checkWrite(t, -2147483648, [MaxVarIntLen]byte{128, 128, 128, 128, 8})
}

func TestReadFrom(t *testing.T) {
	checkRead(t, 0, [MaxVarIntLen]byte{0})
	checkRead(t, 1, [MaxVarIntLen]byte{1})
	checkRead(t, 2, [MaxVarIntLen]byte{2})
	checkRead(t, 127, [MaxVarIntLen]byte{127})
	checkRead(t, 128, [MaxVarIntLen]byte{128, 1})
	checkRead(t, 255, [MaxVarIntLen]byte{255, 1})
	checkRead(t, 25565, [MaxVarIntLen]byte{221, 199, 1})
	checkRead(t, 2097151, [MaxVarIntLen]byte{255, 255, 127})
	checkRead(t, 2147483647, [MaxVarIntLen]byte{255, 255, 255, 255, 7})
	checkRead(t, -1, [MaxVarIntLen]byte{255, 255, 255, 255, 15})
	checkRead(t, -2147483648, [MaxVarIntLen]byte{128, 128, 128, 128, 8})
}

func BenchmarkVarInt_WriteTo(b *testing.B) {
	b.ReportAllocs()
	vi := VarInt(25565)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vi.WriteTo(io.Discard)
	}
}
