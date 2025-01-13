package ngram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Adding_values_less_than_4_chars_contains(t *testing.T) {
	testFilter(t, []string{"foo"}, []string{"foo"}, []string{"bar", "asfs"})
}

func Test_Adding_values_2_chars_contains(t *testing.T) {
	t.Skip("Skipping as we don't support less than four char substring search")
	testFilter(t, []string{"foo fo"}, []string{"fo"}, []string{"bar", "asfs"})
}

func Test_Adding_values_with_4_chars_contains(t *testing.T) {
	testFilter(t, []string{"foos"}, []string{"foos"}, []string{"baras", "asfs", "foo", "oos"})
}

func Test_Adding_values_with_more_than_4_chars_contains(t *testing.T) {
	testFilter(t, []string{"foos bar ball"}, []string{"foos", "ball", "bar ", "foos bar ball"},
		[]string{"abars", "foo2", " foos bar ball "})
}

func testFilter(t *testing.T, input, contained, not_contained []string) {
	t.Helper()
	b := NewBuilder(len(input))
	for _, v := range input {
		b.Add("not-used", v)
	}
	f, err := b.Build()
	assert.NoError(t, err)

	t.Run("contains", func(t *testing.T) {
		for _, v := range contained {
			t.Run(v, func(t *testing.T) {
				assert.True(t, f.Contains(v))
			})
		}
	})

	t.Run("does not contain", func(t *testing.T) {
		for _, not_c := range not_contained {
			t.Run(not_c, func(t *testing.T) {
				assert.False(t, f.Contains(not_c))
			})
		}
	})
}

func Test__SubstringIterator(t *testing.T) {
	ttable := []struct {
		input  string
		output []string
	}{
		{
			input:  "",
			output: []string{""},
		},
		{
			input:  "fo",
			output: []string{"fo"},
		},

		{
			input:  "foo",
			output: []string{"foo"},
		},
		{
			input:  "foob",
			output: []string{"foob"},
		},
		{
			input:  "foobar",
			output: []string{"foob", "ooba", "obar"},
		},
	}

	for _, tt := range ttable {
		t.Run(tt.input, func(t *testing.T) {
			i := 0
			for v := range SubstringIterator(tt.input) {
				assert.Equal(t, tt.output[i], v)
				i += 1
			}
		})
	}
}

func Test__SubstringIteratorWithMarkers(t *testing.T) {
	ttable := []struct {
		input  string
		output []string
	}{
		{
			input:  "",
			output: []string{"^$"},
		},
		{
			input:  "fo",
			output: []string{"fo", "^fo$"},
		},

		{
			input:  "foo",
			output: []string{"^foo", "foo", "foo$"},
		},
		{
			input:  "foob",
			output: []string{"^foo", "foob", "oob$"},
		},
		{
			input:  "foobar",
			output: []string{"^foo", "foob", "ooba", "obar", "bar$"},
		},
	}

	for _, tt := range ttable {
		t.Run(tt.input, func(t *testing.T) {
			i := 0
			for v := range SubstringIteratorWithMarkers(tt.input) {
				assert.Equal(t, tt.output[i], v)
				i += 1
			}
		})
	}
}
