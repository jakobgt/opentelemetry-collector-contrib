package ngram

import (
	"github.com/FastFilter/xorfilter"
	"github.com/twmb/murmur3"
)

type Filter struct {
	fuse *xorfilter.BinaryFuse[uint16]
}

func (f *Filter) Contains(value string) bool {
	for k := range SubstringIterator(value) {
		hk := murmur3.StringSum64(k)
		if !f.fuse.Contains(hk) {
			return false
		}
	}
	return true
}
