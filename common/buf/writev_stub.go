//go:build !windows

package buf

func newVectorizedWriter() (vectorizedWriter, bool) {
	return nil, false
}
