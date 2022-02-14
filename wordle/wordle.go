package wordle

import (
	"fmt"
	"math/rand"
)

type wordle struct {
	alphabets    map[rune]CharacterStatus
	answer       string
	attempts     []Attempt
	dictionary   *[]string
	maxAttempts  int
	options      []Option
	solved       bool
	wordFilters  Filters
	wordsAllowed []string
}

func New(opts ...Option) (Wordle, error) {
	w := &wordle{}
	w.options = append(defaultOpts, opts...)

	err := w.init()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *wordle) Alphabets() map[rune]CharacterStatus {
	return w.alphabets
}

func (w *wordle) Answer() string {
	return w.answer
}

func (w *wordle) Attempt(word string) (*Attempt, error) {
	if w.solved {
		return nil, fmt.Errorf("the last attempt succeeded; no more attempts allowed")
	}
	if len(w.attempts) >= w.maxAttempts {
		return nil, fmt.Errorf("attempted %d times; no more attempts allowed", w.maxAttempts)
	}
	if len(word) != len(w.answer) {
		return nil, fmt.Errorf("word length [%d] does not match answer length [%d]", len(word), len(w.answer))
	}
	notFound := true
	for _, dictWord := range *w.dictionary {
		if dictWord == word {
			notFound = false
			break
		}
	}
	if notFound {
		return nil, fmt.Errorf("not a valid word: '%s'", word)
	}

	attempt := Attempt{
		Answer: word,
		Result: make([]CharacterStatus, len(word)),
	}
	for idx := range word {
		if word[idx] == w.answer[idx] {
			attempt.Result[idx] = PresentInCorrectLocation
			w.alphabets[rune(word[idx])] = PresentInCorrectLocation
		} else {
			for answerIdx := range w.answer {
				if word[idx] == w.answer[answerIdx] {
					attempt.Result[idx] = PresentInWrongLocation
					if w.alphabets[rune(word[idx])] < PresentInWrongLocation {
						w.alphabets[rune(word[idx])] = PresentInWrongLocation
					}
				}
			}
		}
	}
	w.attempts = append(w.attempts, attempt)
	if attempt.Answer == w.answer {
		w.solved = true
	}
	return &attempt, nil
}

func (w wordle) Attempts() []Attempt {
	return w.attempts
}

func (w *wordle) Reset() error {
	return w.init()
}

func (w wordle) Solved() bool {
	return w.solved
}

func (w *wordle) init() error {
	for _, opt := range w.options {
		opt(w)
	}
	w.wordsAllowed = w.wordFilters.Apply(w.dictionary)
	if len(w.wordsAllowed) == 0 {
		return fmt.Errorf("found no words to choose from after applying all filters")
	}

	w.alphabets = make(map[rune]CharacterStatus, 26)
	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		w.alphabets[r] = Unknown
	}
	if w.answer == "" {
		w.answer = w.wordsAllowed[rand.Intn(len(w.wordsAllowed))]
	}
	w.attempts = make([]Attempt, 0, w.maxAttempts)
	w.solved = false

	return nil
}
