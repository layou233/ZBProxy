package mcprotocol

import (
	"errors"
	"io"

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

func EncodeVarInt(n int32) ([MaxVarIntLen]byte, int) {
	var vi [MaxVarIntLen]byte
	num := uint32(n)
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

	return vi, numWrite
}
