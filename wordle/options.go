package wordle

import (
	"strings"

	"github.com/jedib0t/go-wordle/words"
)

// Option helps customize the Wordle game.
type Option func(w *wordle)

var (
	// wordsEnglish contains all available English words
	wordsEnglish = words.English()

	// defaultOpts are applied before user provided options
	defaultOpts = []Option{
		WithDictionary(&wordsEnglish),
		WithMaxAttempts(5),
		WithWordFilters(
			WithLength(5, 5),
		),
	}
)

// WithAnswer sets up the answer and prevents random selection.
func WithAnswer(answer string) Option {
	return func(w *wordle) {
		w.answer = answer
	}
}

// WithDictionary sets the list of words allowed for answers and attempts.
func WithDictionary(dict *[]string) Option {
	return func(w *wordle) {
		w.dictionary = dict
	}
}

// WithMaxAttempts sets up the maximum number of attempts allowed.
func WithMaxAttempts(n int) Option {
	return func(w *wordle) {
		w.maxAttempts = n
	}
}

// WithWordFilters sets up filters to filter the list of words used.
func WithWordFilters(filters ...Filter) Option {
	return func(w *wordle) {
		w.wordFilters = filters
	}
}

// WithUnknownAnswer sets up the puzzle with an unknown answer; can be used to
// solve puzzle externally by just leveraging the hinting system.
func WithUnknownAnswer(wordLen int) Option {
	return func(w *wordle) {
		w.answer = strings.Repeat(" ", wordLen)
		w.answerUnknown = true
	}
}
