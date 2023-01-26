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
	num := uint32(i)
	numWrite := 0
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		vi[numWrite] = byte(b)
		numWrite++
		if num == 0 {
			break
		}
	}

	nn, err := w.Write(vi[:numWrite])
	return int64(nn), err
}

func (i VarInt) WriteToBuffer(buffer *buf.Buffer) {
	i32 := int32(i)
	WriteVarIntTo(buffer.Extend(VarIntLen(i32)), i32)
}

// WriteVarIntTo writes an encoded VarInt into bs.
// Note that it would panic when bs length not enough.
func WriteVarIntTo(bs []byte, n int32) (i int) {
	num := uint32(n)
	numWrite := 0
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		bs[numWrite] = byte(b)
		numWrite++
		if num == 0 {
			break
		}
	}

	return numWrite
}

func VarIntLen(n int32) int {
	switch {
	case n < 0:
		return 5
	case n < 128:
		return 1
	case n < 16384:
		return 2
	case n < 2097152:
		return 3
	case n < 268435456:
		return 4
	default:
		return 5
	}
}

func ReadVarIntFrom(r io.Reader) (i int32, n int64, err error) {
	var V uint32
	var num int64
	for sec := byte(0x80); sec&0x80 != 0; num++ {
		if num > MaxVarIntLen {
			return 0, n, ErrVarIntTooBig
		}

		sec, err = rw.ReadByte(r)
		if err != nil {
			return 0, n, err
		}
		n += 1

		V |= uint32(sec&0x7F) << uint32(7*num)
	}

	i = int32(V)
	return
}
