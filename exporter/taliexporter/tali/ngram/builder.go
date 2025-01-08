package ngram

import (
	"fmt"
	"iter"

	"github.com/FastFilter/xorfilter"
	"github.com/twmb/murmur3"
	_ "github.com/twmb/murmur3"
)

const (
	ngramLen    = 4
	startMarker = "^"
	endMarker   = "$"
)

type Builder struct {
	//filter xorfilter.BinaryFuse[uint8]
	keys map[uint64]struct{}
}

func NewBuilder(size int) Builder {
	return Builder{
		keys: map[uint64]struct{}{},
	}
}

func (b *Builder) Add(key, value string) error {
	for k := range SubstringIteratorWithMarkers(value) {
		hash := murmur3.StringSum64(k)
		_, ok := b.keys[hash]
		if ok {
			return nil
		}
		b.keys[hash] = struct{}{}
	}

	return nil
}

func (b *Builder) Build() (Filter, error) {
	keys := make([]uint64, 0, len(b.keys))
	for k := range b.keys {
		keys = append(keys, k)
	}

	// Technically, the error condition should not happen, as it is when
	// the fuse filter do too many iterations.
	f, err := xorfilter.NewBinaryFuse[uint16](keys)
	if err != nil {
		return Filter{}, err
	}
	return Filter{
		fuse: f,
	}, nil
}

// SubstringIterator returns subslices of length 4 that
func SubstringIterator(value string) iter.Seq[string] {
	// Should we skip spaces and other stuff?
	return func(yield func(string) bool) {
		if len(value) < ngramLen {
			yield(value)
			return
		}
		// We then loop over 4-length subslices
		pos := 0
		for ; pos+ngramLen <= len(value); pos += 1 {
			// once we've finished looping
			if !yield(value[pos : pos+ngramLen]) {
				// we then return and finish our iterations
				return
			}
		}
	}
}

// SubstringIterator returns subslices of length 4 that
func SubstringIteratorWithMarkers(value string) iter.Seq[string] {
	// Should we skip spaces and other stuff?
	return func(yield func(string) bool) {
		switch len(value) {
		case 0:
			yield(startMarker + endMarker)
			return
		case 1, 2:
			if !yield(value) {
				return
			}
			toAdd := fmt.Sprintf("%s%s%s", startMarker, value, endMarker)
			yield(toAdd)
			return
		case 3:
			// We hardcode this scenario
			prefix := fmt.Sprintf("%s%s", startMarker, value[0:3])
			if !yield(prefix) {
				return
			}
			if !yield(value) {
				return
			}
			suffix := fmt.Sprintf("%s%s", value, endMarker)
			yield(suffix)
			return
		default:
		}
		// First yield is the start marker
		prefix := fmt.Sprintf("%s%s", startMarker, value[0:3])
		if !yield(prefix) {
			return
		}

		for k := range SubstringIterator(value) {
			if !yield(k) {
				return
			}
		}
		// Lastly we handled the end-marker
		endPos := len(value) - 3
		suffix := fmt.Sprintf("%s%s", value[endPos:], endMarker)
		yield(suffix)
	}
}
