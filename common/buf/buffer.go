package buf

import (
	"io"
	"strconv"
)

type Buffer struct {
	data    []byte
	start   int
	end     int
	managed bool
	closed  bool
}

func NewSize(size int) *Buffer {
	if size > 65535 {
		return &Buffer{
			data: make([]byte, size),
		}
	}
	return &Buffer{
		data:    Get(size),
		managed: true,
	}
}

func As(data []byte) *Buffer {
	return &Buffer{
		data: data,
		end:  len(data),
	}
}

func With(data []byte) *Buffer {
	return &Buffer{
		data: data,
	}
}

func (b *Buffer) Closed() bool {
	return b.closed
}

func (b *Buffer) Byte(index int) byte {
	return b.data[b.start+index]
}

func (b *Buffer) SetByte(index int, value byte) {
	b.data[b.start+index] = value
}

func (b *Buffer) Extend(n int) []byte {
	end := b.end + n
	if end > cap(b.data) {
		panic("buffer overflow: cap " + strconv.Itoa(cap(b.data)) + ",end " + strconv.Itoa(b.end) + ", need " + strconv.Itoa(n))
	}
	ext := b.data[b.end:end]
	b.end = end
	return ext
}

func (b *Buffer) Advance(from int) {
	b.start += from
}

func (b *Buffer) Truncate(to int) {
	b.end = b.start + to
}

func (b *Buffer) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return
	}
	if b.IsFull() {
		return 0, io.ErrShortBuffer
	}
	n = copy(b.data[b.end:], data)
	b.end += n
	return
}

func (b *Buffer) ExtendHeader(n int) []byte {
	if b.start < n {
		panic("buffer overflow: cap " + strconv.Itoa(cap(b.data)) + ",start " + strconv.Itoa(b.start) + ", need " + strconv.Itoa(n))
	}
	b.start -= n
	return b.data[b.start : b.start+n]
}

func (b *Buffer) WriteByte(d byte) error {
	if b.IsFull() {
		return io.ErrShortBuffer
	}
	b.data[b.end] = d
	b.end++
	return nil
}

func (b *Buffer) ReadOnceFrom(r io.Reader) (int64, error) {
	if b.IsFull() {
		return 0, io.ErrShortBuffer
	}
	n, err := r.Read(b.FreeBytes())
	b.end += n
	return int64(n), err
}

func (b *Buffer) ReadAtLeastFrom(r io.Reader, min int) (int64, error) {
	if min <= 0 {
		return b.ReadOnceFrom(r)
	}
	if b.IsFull() {
		return 0, io.ErrShortBuffer
	}
	n, err := io.ReadAtLeast(r, b.FreeBytes(), min)
	b.end += n
	return int64(n), err
}

func (b *Buffer) ReadFullFrom(r io.Reader, size int) (n int, err error) {
	if b.end+size > b.Cap() {
		return 0, io.ErrShortBuffer
	}
	n, err = io.ReadFull(r, b.data[b.end:b.end+size])
	b.end += n
	return
}

func (b *Buffer) ReadFrom(reader io.Reader) (n int64, err error) {
	for {
		if b.IsFull() {
			return 0, io.ErrShortBuffer
		}
		var readN int
		readN, err = reader.Read(b.FreeBytes())
		b.end += readN
		n += int64(readN)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
}

func (b *Buffer) WriteRune(s rune) (int, error) {
	return b.Write([]byte{byte(s)})
}

func (b *Buffer) WriteString(s string) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}
	if b.IsFull() {
		return 0, io.ErrShortBuffer
	}
	n := copy(b.data[b.end:], s)
	b.end += n
	return n, nil
}

func (b *Buffer) WriteZero() error {
	if b.IsFull() {
		return io.ErrShortBuffer
	}
	b.data[b.end] = 0
	b.end++
	return nil
}

func (b *Buffer) WriteZeroN(n int) error {
	if b.end+n > b.Cap() {
		return io.ErrShortBuffer
	}
	for i := b.end; i <= b.end+n; i++ {
		b.data[i] = 0
	}
	b.end += n
	return nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.IsEmpty() {
		return 0, io.EOF
	}

	nb := b.data[b.start]
	b.start++
	return nb, nil
}

func (b *Buffer) Read(data []byte) (n int, err error) {
	if b.IsEmpty() {
		return 0, io.EOF
	}
	n = copy(data, b.data[b.start:b.end])
	b.start += n
	return
}

func (b *Buffer) Peek(n int) (bytes []byte, err error) {
	if b.start+n > b.end {
		return nil, io.ErrShortBuffer
	}
	bytes = b.To(n)
	b.start += n
	return
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.Bytes())
	return int64(n), err
}

func (b *Buffer) Resize(start, end int) {
	b.start = start
	b.end = b.start + end
}

func (b *Buffer) Reset(pos int) {
	b.start = pos
	b.end = pos
}

func (b *Buffer) FullReset() {
	b.start = 0
	b.end = 0
}

func (b *Buffer) Release() {
	if b == nil || b.closed || !b.managed {
		return
	}
	Put(b.data) //nolint:errcheck
	*b = Buffer{closed: true}
}

func (b *Buffer) Cut(start int, end int) *Buffer {
	b.start += start
	b.end = len(b.data) - end
	return &Buffer{
		data: b.data[b.start:b.end],
	}
}

func (b *Buffer) Start() int {
	return b.start
}

func (b *Buffer) Len() int {
	return b.end - b.start
}

func (b *Buffer) Cap() int {
	return len(b.data)
}

func (b *Buffer) Bytes() []byte {
	return b.data[b.start:b.end]
}

func (b *Buffer) Slice() []byte {
	return b.data
}

func (b *Buffer) From(n int) []byte {
	return b.data[b.start+n : b.end]
}

func (b *Buffer) To(n int) []byte {
	return b.data[b.start : b.start+n]
}

func (b *Buffer) Range(start, end int) []byte {
	return b.data[b.start+start : b.start+end]
}

func (b *Buffer) Index(start int) []byte {
	return b.data[b.start+start : b.start+start]
}

func (b *Buffer) FreeLen() int {
	return b.Cap() - b.end
}

func (b *Buffer) FreeBytes() []byte {
	return b.data[b.end:b.Cap()]
}

func (b *Buffer) Rewind(start int) {
	b.start = start
}

func (b *Buffer) IsEmpty() bool {
	return b.end-b.start == 0
}

func (b *Buffer) IsFull() bool {
	return b.end == b.Cap()
}

func (b *Buffer) ToOwned() *Buffer {
	n := NewSize(len(b.data))
	copy(n.data[b.start:b.end], b.data[b.start:b.end])
	n.start = b.start
	n.end = b.end
	return n
}
