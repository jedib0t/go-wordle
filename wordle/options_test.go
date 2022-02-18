package wordle

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithAnswer(t *testing.T) {
	w := &wordle{}
	assert.Empty(t, w.answer)
	WithAnswer("foo")(w)
	assert.Equal(t, "foo", w.answer)
}

func TestWithDictionary(t *testing.T) {
	dictionary := []string{
		"aardvark",
	}
	w := &wordle{}
	assert.Nil(t, w.dictionary)
	WithDictionary(&dictionary)(w)
	assert.NotNil(t, w.dictionary)
	assert.Equal(t, w.dictionary, &dictionary)
}

func TestWithMaxAttempts(t *testing.T) {
	w := &wordle{}
	assert.Zero(t, w.maxAttempts)
	WithMaxAttempts(5)(w)
	assert.Equal(t, 5, w.maxAttempts)
}

func TestWithUnknownAnswer(t *testing.T) {
	w := &wordle{}
	assert.Empty(t, w.answer)
	assert.False(t, w.answerUnknown)
	WithUnknownAnswer(5)(w)
	assert.Equal(t, strings.Repeat(" ", 5), w.answer)
	assert.True(t, w.answerUnknown)
}
