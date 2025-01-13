package ngram

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func _test_random_string(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func Test__Mar(t *testing.T) {
	b := NewBuilder(2)
	for i := 0; i < 500; i++ {
		assert.NoError(t, b.Add("key1", _test_random_string(200)))
	}
	// GIVEN: A Binary fuse filter built
	expected_filter, err := b.Build()
	require.NoError(t, err)

	// GIVEN: Marshal the filter
	bytz, err := expected_filter.Marshal()
	require.NoError(t, err)
	assert.NotEmpty(t, bytz)

	// GIVEN: Unmarshal the byte stream
	after_marshal, err := Unmarshal(bytz)
	require.NoError(t, err)

	// EXPECT: The marshalled + unmarshalled version match.
	assert.Equal(t, expected_filter, after_marshal)
}

// TODO: Test the concurrency nature of Unmarshal and Marshal

// TODO: Fuzzy testing the Unmarshal function.
