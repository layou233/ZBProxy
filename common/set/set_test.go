package set

import (
	"testing"

	"github.com/zhangyunhao116/fastrand"
	//"golang.org/x/exp/slices"
)

func BenchmarkStringSet_Has(b *testing.B) {
	const ElementNum = 64
	listSlice := make([]string, 0, ElementNum)
	for i := 0; i < ElementNum; i++ {
		var nameBytes [16]byte
		fastrand.Read(nameBytes[:])
		listSlice = append(listSlice, string(nameBytes[:]))
	}
	s := NewStringSetFromSlice(listSlice)
	target := listSlice[fastrand.Int31n(ElementNum)]
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !s.Has(target) {
			b.Fatal("target element not found")
		}
	}
}

/*
func BenchmarkStrings_Contain(b *testing.B) {
	const ElementNum = 64
	listSlice := make([]string, 0, ElementNum)
	for i := 0; i < ElementNum; i++ {
		var nameBytes [16]byte
		rand.Read(nameBytes[:])
		listSlice = append(listSlice, string(nameBytes[:]))
	}
	target := listSlice[mrand.Int31n(ElementNum)]
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !slices.Contains(listSlice, target) {
			b.Fatal("target element not found")
		}
	}
}
*/
