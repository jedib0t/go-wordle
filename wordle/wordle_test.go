package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	w, err := New(WithMaxAttempts(100))
	assert.NotNil(t, w)
	assert.Nil(t, err)
	assert.Len(t, w.Alphabets(), 26)
	assert.NotEmpty(t, w.Dictionary())
	assert.Len(t, w.Answer(), defaultWordLength)
	assert.Empty(t, w.Attempts())
	assert.False(t, w.GameOver())
	assert.False(t, w.Solved())

	wObj, ok := w.(*wordle)
	assert.NotNil(t, wObj)
	assert.True(t, ok)
	assert.Equal(t, 100, wObj.maxAttempts)
	for _, word := range wObj.wordsAllowed {
		assert.Len(t, word, defaultWordLength)
	}
}

func TestWordle_Alphabets(t *testing.T) {
	w := &wordle{}
	assert.Empty(t, w.alphabets)
	assert.Equal(t, NotPresent, w.Alphabets()["a"])

	w.alphabets = map[string]CharacterStatus{
		"a": CorrectLocation,
	}
	assert.NotEmpty(t, w.alphabets)
	assert.Equal(t, CorrectLocation, w.Alphabets()["a"])
}

func TestWordle_Answer(t *testing.T) {
	w := &wordle{}
	assert.Empty(t, w.answer)
	assert.Empty(t, w.Answer())

	w.answer = "foo"
	assert.Equal(t, "foo", w.Answer())
}

func TestWordle_AnswerUnknown(t *testing.T) {
	w := &wordle{}
	assert.False(t, w.answerUnknown)
	assert.False(t, w.AnswerUnknown())

	w.answerUnknown = true
	assert.True(t, w.answerUnknown)
	assert.True(t, w.AnswerUnknown())
}

func TestWordle_Attempt(t *testing.T) {
	w, err := New(WithAnswer("train"), WithMaxAttempts(2))
	assert.NotNil(t, w)
	assert.Nil(t, err)
	assert.Len(t, w.Attempts(), 0)

	attempt, err := w.Attempt("tribe")
	assert.NotNil(t, attempt)
	assert.Nil(t, err)
	assert.Len(t, w.Attempts(), 1)
	assert.Equal(t, "tribe", attempt.Answer)
	assert.Equal(t,
		[]CharacterStatus{CorrectLocation, CorrectLocation, WrongLocation, NotPresent, NotPresent},
		attempt.Result,
	)
	assert.Equal(t, CorrectLocation, w.Alphabets()["t"])
	assert.Equal(t, CorrectLocation, w.Alphabets()["r"])
	assert.Equal(t, WrongLocation, w.Alphabets()["i"])
	assert.Equal(t, NotPresent, w.Alphabets()["b"])
	assert.Equal(t, NotPresent, w.Alphabets()["e"])
	assert.False(t, w.GameOver())
	assert.False(t, w.Solved())

	attempt, err = w.Attempt("baron")
	assert.NotNil(t, attempt)
	assert.Nil(t, err)
	assert.Len(t, w.Attempts(), 2)
	assert.Equal(t, "baron", attempt.Answer)
	assert.Equal(t,
		[]CharacterStatus{NotPresent, WrongLocation, WrongLocation, NotPresent, CorrectLocation},
		attempt.Result,
	)
	assert.Equal(t, CorrectLocation, w.Alphabets()["t"])
	assert.Equal(t, CorrectLocation, w.Alphabets()["r"])
	assert.Equal(t, WrongLocation, w.Alphabets()["i"])
	assert.Equal(t, NotPresent, w.Alphabets()["b"])
	assert.Equal(t, NotPresent, w.Alphabets()["e"])
	assert.Equal(t, WrongLocation, w.Alphabets()["a"])
	assert.Equal(t, NotPresent, w.Alphabets()["o"])
	assert.Equal(t, CorrectLocation, w.Alphabets()["n"])
	assert.True(t, w.GameOver())
	assert.False(t, w.Solved())

	attemptIncorrect, err := w.Attempt("train")
	assert.NotNil(t, attempt)
	assert.Nil(t, err)
	assert.Equal(t, attempt, attemptIncorrect)
}

func TestWordle_Attempts(t *testing.T) {
	w := &wordle{}
	assert.Empty(t, w.attempts)
	assert.Empty(t, w.Attempts())

	w.attempts = []Attempt{
		{Answer: "foo", Result: []CharacterStatus{NotPresent, WrongLocation, CorrectLocation}},
	}
	assert.NotEmpty(t, w.attempts)
	assert.Len(t, w.Attempts(), len(w.attempts))
}

func TestWordle_Dictionary(t *testing.T) {
	w := &wordle{}
	assert.Nil(t, w.dictionary)
	assert.Nil(t, w.Dictionary())

	w.dictionary = &[]string{"foo", "bar"}
	assert.NotNil(t, w.dictionary)
	assert.NotNil(t, w.Dictionary())
	assert.Len(t, w.Dictionary(), len(*w.dictionary))
	assert.Len(t, w.Dictionary(), len(*w.dictionary))
}

func TestWordle_DictionaryHas(t *testing.T) {
	w := &wordle{}
	assert.False(t, w.DictionaryHas("foo"))

	w.dictionary = &[]string{"foo", "bar"}
	assert.True(t, w.DictionaryHas("foo"))
	assert.True(t, w.DictionaryHas("bar"))
	assert.False(t, w.DictionaryHas("baz"))
}

func TestWordle_GameOver(t *testing.T) {
	w := &wordle{maxAttempts: 1}
	assert.False(t, w.solved)
	assert.False(t, w.GameOver())

	w.solved = true
	assert.True(t, w.GameOver())

	w.solved = false
	w.attempts = []Attempt{{Answer: "foo"}, {Answer: "bar"}}
	w.maxAttempts = 1
	assert.False(t, w.GameOver())

	w.maxAttempts = 2
	assert.True(t, w.GameOver())
}

func TestWordle_Hints(t *testing.T) {
	w := &wordle{}
	w.alphabets = make(map[string]CharacterStatus, 26)
	for _, r := range englishAlphabets {
		w.alphabets[string(r)] = Unknown
	}
	assert.Empty(t, w.Hints())

	w.wordsAllowed = []string{"foo", "bar", "baz", "foof", "barb", "bazz"}
	assert.Len(t, w.Hints(), maxHints)
}

func TestWordle_Reset(t *testing.T) {
	w, err := New()
	assert.NotNil(t, w)
	assert.Nil(t, err)
	assert.Len(t, w.Alphabets(), 26)
	assert.NotEmpty(t, w.Dictionary())
	assert.Len(t, w.Answer(), defaultWordLength)
	assert.Empty(t, w.Attempts())
	assert.False(t, w.GameOver())
	assert.False(t, w.Solved())

	wObj, ok := w.(*wordle)
	assert.NotNil(t, wObj)
	assert.True(t, ok)
	for _, word := range wObj.wordsAllowed {
		assert.Len(t, word, defaultWordLength)
	}
}

func TestWordle_Solved(t *testing.T) {
	w := &wordle{}
	assert.False(t, w.solved)
	assert.False(t, w.Solved())

	w.solved = true
	assert.True(t, w.Solved())
}
