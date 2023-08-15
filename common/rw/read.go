package rw

import "io"

type ByteReaderWrapper struct {
	io.Reader
}

func (r ByteReaderWrapper) ReadByte() (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(r.Reader, b[:])
	return b[0], err
}

func CreateByteReader(reader io.Reader) io.ByteReader {
	if br, ok := reader.(io.ByteReader); ok {
		return br
	}
	return ByteReaderWrapper{reader}
}

func ReadByte(reader io.Reader) (byte, error) {
	if br, ok := reader.(io.ByteReader); ok {
		return br.ReadByte()
	}
	var b [1]byte
	_, err := io.ReadFull(reader, b[:])
	return b[0], err
}

func ReadBytes(reader io.Reader, size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := io.ReadFull(reader, b); err != nil {
		return nil, err
	}
	return b, nil
}
