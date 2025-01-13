package ngram

import (
	"encoding/json"

	"github.com/FastFilter/xorfilter"
	"github.com/klauspost/compress/zstd"
	"github.com/twmb/murmur3"
)

// Create a writer that caches compressors.
// For this operation type we supply a nil Reader.
var compressor, _ = zstd.NewWriter(nil)
var decompressor, _ = zstd.NewReader(nil)

type Filter struct {
	fuse *xorfilter.BinaryFuse[uint16]
}

// Contains return true if
func (f *Filter) Contains(value string) bool {
	for k := range SubstringIterator(value) {
		hk := murmur3.StringSum64(k)
		if !f.fuse.Contains(hk) {
			return false
		}
	}
	return true
}

// Marshal returns a marshalled and compressed verison of the Binary fuse filter
func (f *Filter) Marshal() ([]byte, error) {
	// TODO: Consider whether to serialize with something else?
	bytz, err := json.Marshal(f.fuse)
	if err != nil {
		return nil, err
	}

	// TODO: Allow to give the byte arrayslice, to allow for manual memory management.
	dst := []byte{}
	return compressor.EncodeAll(bytz, dst), nil
}

// Unmarshal returns the parsed Filter from the byte stream, or an error if one
// is encountered.
func Unmarshal(src []byte) (Filter, error) {
	dst := []byte{}
	dst, err := decompressor.DecodeAll(src, dst)
	if err != nil {
		return Filter{}, nil
	}
	f := &xorfilter.BinaryFuse[uint16]{}
	err = json.Unmarshal(dst, f)
	if err != nil {
		return Filter{}, err
	}
	return Filter{
		fuse: f,
	}, nil
}
