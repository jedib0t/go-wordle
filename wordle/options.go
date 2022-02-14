package wordle

import "github.com/jedib0t/go-wordle/words"

type Option func(w *wordle)

var (
	wordsEnglish = words.English()

	defaultOpts = []Option{
		WithDictionary(&wordsEnglish),
		WithMaxAttempts(5),
		WithWordFilters([]Filter{
			WithLength(5, 5),
			WithNoRepeatingCharacters(),
		}),
	}
)

func WithAnswer(answer string) Option {
	return func(w *wordle) {
		w.answer = answer
	}
}

func WithAttempts(attempts []Attempt) Option {
	return func(w *wordle) {
		w.attempts = attempts
	}
}

func WithDictionary(dict *[]string) Option {
	return func(w *wordle) {
		w.dictionary = dict
	}
}

func WithMaxAttempts(n int) Option {
	return func(w *wordle) {
		w.maxAttempts = n
	}
}

func WithWordFilters(filters []Filter) Option {
	return func(w *wordle) {
		w.wordFilters = filters
	}
}
