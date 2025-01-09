package bufio

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/buf"
)

type CachedConn struct {
	net.Conn
	cache *buf.Buffer
}

var (
	_ common.WrappedReader = (*CachedConn)(nil)
	_ common.WrappedWriter = (*CachedConn)(nil)
)

func NewCachedConn(c net.Conn) *CachedConn {
	if cachedConn, isCachedConn := c.(*CachedConn); isCachedConn {
		return cachedConn
	}
	return &CachedConn{
		Conn: c,
	}
}

func (c *CachedConn) Cache() *buf.Buffer {
	return c.cache
}

func (c *CachedConn) Read(p []byte) (n int, err error) {
	if c.cache != nil && !c.cache.IsEmpty() {
		return c.cache.Read(p)
	}
	n, err = c.Conn.Read(p)
	if n > 0 {
		if c.cache == nil {
			// allocate 4 KiB here since the Linux page size
			// is typically 4096 bytes for x86(-64) processors.
			// we hope to get some performance benefits from here.
			// however, since the buffer size can't increase by itself,
			// we wouldn't support sniffing any protocol that has
			// more than 4096 bytes in handshake packet WITHOUT modifying this.
			c.cache = buf.NewSize(4096)
		}
		_n, _ := c.cache.Write(p[:n])
		c.cache.Advance(_n)
		if _n != n {
			return 0, io.ErrShortBuffer
		}
	}
	return
}

func (c *CachedConn) Peek(n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("peek %d bytes: %w", n, buf.ErrNegativeRead)
	}
	if c.cache == nil {
		// see above
		c.cache = buf.NewSize(4096)
	}
	if need := n - c.cache.Len(); need > 0 {
		_, err := c.cache.ReadAtLeastFrom(c.Conn, need)
		if err != nil {
			return nil, err
		}
	}
	return c.cache.Peek(n)
}

func (c *CachedConn) PeekUntil(end ...[]byte) ([]byte, []byte, error) {
	if c.cache == nil {
		// see above
		c.cache = buf.NewSize(4096)
	}
	for {
		buffer := c.cache.Bytes()
		for _, e := range end {
			if index := bytes.Index(buffer, e); index >= 0 {
				c.cache.Advance(index + len(e))
				return buffer[:index], e, nil
			}
		}
		_, err := c.cache.ReadAtLeastFrom(c.Conn, 1)
		if err != nil {
			return nil, nil, err
		}
	}
}

func (c *CachedConn) Rewind(position int) {
	if c.cache != nil {
		c.cache.Rewind(position)
	}
}

func (c *CachedConn) Release() {
	c.cache.Release()
	c.cache = nil
}

func (c *CachedConn) CurrentPosition() int {
	if c.cache == nil {
		return -1
	}
	return c.cache.CurrentPosition()
}

func (c *CachedConn) Close() error {
	c.Release()
	return c.Conn.Close()
}

func (c *CachedConn) UpstreamReader() io.Reader {
	if c.cache == nil || c.cache.IsEmpty() {
		return c.Conn
	}
	return nil
}

func (c *CachedConn) UpstreamWriter() io.Writer {
	return c.Conn
}
