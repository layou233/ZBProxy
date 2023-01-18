package rw

import "io"

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
	if _, err := reader.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
