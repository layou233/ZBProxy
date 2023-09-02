package mcprotocol

import (
	"errors"
	"io"

	"github.com/layou233/ZBProxy/common/buf"
	"github.com/layou233/ZBProxy/common/rw"
)

const MaxVarIntLen = 5

var ErrVarIntTooBig = errors.New("VarInt is too big")

type VarInt int32

func (i VarInt) Value() int {
	return int(i)
}

func (i VarInt) Value32() int32 {
	return int32(i)
}

func (i VarInt) WriteTo(w io.Writer) (n int64, err error) {
	var vi [MaxVarIntLen]byte
	numWrite := PutVarInt(vi[:], int32(i))
	nn, err := w.Write(vi[:numWrite])
	return int64(nn), err
}

func (i VarInt) WriteToBuffer(buffer *buf.Buffer) {
	i32 := int32(i)
	PutVarInt(buffer.Extend(VarIntLen(i32)), i32)
}

// PutVarInt encodes a Minecraft variable-length format int32 into bs and returns the number of bytes written.
// If the buffer is too small, PutVarInt will panic.
func PutVarInt(bs []byte, n int32) (numWrite int) {
	num := uint32(n)
	for num != 0 {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		bs[numWrite] = byte(b)
		numWrite++
	}
	return
}

func VarIntLen(n int32) int {
	switch {
	case n < 0:
		return 5
	case n < 1<<(7*1):
		return 1
	case n < 1<<(7*2):
		return 2
	case n < 1<<(7*3):
		return 3
	case n < 1<<(7*4):
		return 4
	default:
		return 5
	}
}

func ReadVarIntFrom(r io.Reader) (i int32, n int64, err error) {
	var v uint32
	br := rw.CreateByteReader(r)
	for sec := byte(0x80); sec&0x80 != 0; n++ {
		if n > MaxVarIntLen {
			return 0, n, ErrVarIntTooBig
		}

		sec, err = br.ReadByte()
		if err != nil {
			return 0, n, err
		}

		v |= uint32(sec&0x7F) << uint32(7*n)
	}

	i = int32(v)
	return
}
