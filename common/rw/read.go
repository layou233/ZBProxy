package rw

import "io"

func ReadByte(reader io.Reader) (byte, error) {
	if br, isBr := reader.(io.ByteReader); isBr {
		return br.ReadByte()
	}
	b, err := ReadBytes(reader, 1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func ReadBytes(reader io.Reader, size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := reader.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
