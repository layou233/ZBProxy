package buf

func Get(size int) []byte {
	if size <= 65536 {
		return DefaultAllocator.Get(size)
	}
	return make([]byte, size)
}

func Put(buf []byte) error {
	if cap(buf) > 65536 {
		return nil
	}
	return DefaultAllocator.Put(buf)
}

func PutMulti(buffers [][]byte) {
	for _, buffer := range buffers {
		Put(buffer) //nolint:errcheck
	}
}
