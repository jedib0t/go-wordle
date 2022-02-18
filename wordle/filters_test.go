package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLength(t *testing.T) {
	assert.True(t, WithLength(5)("brand"))
	assert.True(t, WithLength(4)("barn"))
	assert.True(t, WithLength(3)("bar"))
	assert.False(t, WithLength(5)("brands"))
	assert.False(t, WithLength(4)("barns"))
	assert.False(t, WithLength(3)("bars"))
}

func TestWithNoRepeatingCharacters(t *testing.T) {
	assert.True(t, WithNoRepeatingCharacters()("brand"))
	assert.True(t, WithNoRepeatingCharacters()("barn"))
	assert.True(t, WithNoRepeatingCharacters()("bar"))
	assert.False(t, WithNoRepeatingCharacters()("brandd"))
	assert.False(t, WithNoRepeatingCharacters()("barnn"))
	assert.False(t, WithNoRepeatingCharacters()("barr"))
}

func TestFilters_Apply(t *testing.T) {
	words := []string{
		"brain",
		"drain",
		"grain",
		"rain",
		"strain",
		"train",
		"trait",
	}

	filters := Filters{
		WithLength(5),
		WithNoRepeatingCharacters(),
	}
	assert.Nil(t, filters.Apply(nil))
	assert.Nil(t, filters.Apply(&[]string{}))

	filteredWords := filters.Apply(&words)
	assert.Len(t, filteredWords, 4)
	assert.Contains(t, filteredWords, "brain")
	assert.Contains(t, filteredWords, "drain")
	assert.Contains(t, filteredWords, "grain")
	assert.Contains(t, filteredWords, "train")
	assert.NotContains(t, filteredWords, "rain")
	assert.NotContains(t, filteredWords, "strain")
	assert.NotContains(t, filteredWords, "trait")
}
