package wordle

import (
	"fmt"
	"math/rand"
)

type wordle struct {
	alphabets    map[string]CharacterStatus
	answer       string
	attempts     []Attempt
	dictionary   *[]string
	maxAttempts  int
	options      []Option
	solved       bool
	wordFilters  Filters
	wordsAllowed []string
}

// New generates a new Wordle game with a randomly selected word.
func New(opts ...Option) (Wordle, error) {
	w := &wordle{}
	w.options = append(defaultOpts, opts...)

	err := w.init()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *wordle) Alphabets() map[string]CharacterStatus {
	return w.alphabets
}

func (w *wordle) Answer() string {
	return w.answer
}

func (w *wordle) Attempt(word string) (*Attempt, error) {
	if err := w.validateAttempt(word); err != nil {
		return nil, err
	}

	// prep to record a new attempt
	attempt := Attempt{
		Answer: word,
		// default to NotPresent status for all characters
		Result: make([]CharacterStatus, len(word)),
	}
	// loop through the word one character at a time
	for idx := range word {
		if word[idx] == w.answer[idx] {
			attempt.Result[idx] = PresentInCorrectLocation
			w.alphabets[string(word[idx])] = PresentInCorrectLocation
		} else {
			charNotFound := true
			for answerIdx := range w.answer {
				if word[idx] == w.answer[answerIdx] {
					attempt.Result[idx] = PresentInWrongLocation
					if w.alphabets[string(word[idx])] < PresentInWrongLocation {
						w.alphabets[string(word[idx])] = PresentInWrongLocation
					}
					charNotFound = false
				}
			}
			if charNotFound {
				delete(w.alphabets, string(word[idx]))
			}
		}
	}
	// record the new attempt
	w.attempts = append(w.attempts, attempt)
	// mark as solved if so
	if attempt.Answer == w.answer {
		w.solved = true
	}
	return &attempt, nil
}

func (w wordle) Attempts() []Attempt {
	return w.attempts
}

func (w wordle) Hints() []string {
	if w.solved {
		return nil
	}
	return generateHints(w.wordsAllowed, w.attempts, w.alphabets)
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

	w.alphabets = make(map[string]CharacterStatus, 26)
	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		w.alphabets[string(r)] = Unknown
	}
	if w.answer == "" {
		w.answer = w.wordsAllowed[rand.Intn(len(w.wordsAllowed))]
	}
	w.attempts = make([]Attempt, 0, w.maxAttempts)
	w.solved = false

	return nil
}

func (w *wordle) validateAttempt(word string) error {
	if w.solved {
		return fmt.Errorf("the last attempt succeeded; no more attempts allowed")
	}
	if len(w.attempts) >= w.maxAttempts {
		return fmt.Errorf("attempted %d times; no more attempts allowed", w.maxAttempts)
	}
	if len(word) != len(w.answer) {
		return fmt.Errorf("word length [%d] does not match answer length [%d]", len(word), len(w.answer))
	}
	for _, attempt := range w.attempts {
		if word == attempt.Answer {
			return fmt.Errorf("word [%s] has been attempted already", word)
		}
	}
	notFound := true
	for _, dictWord := range *w.dictionary {
		if dictWord == word {
			notFound = false
			break
		}
	}
	if notFound {
		return fmt.Errorf("not a valid word: '%s'", word)
	}
	return nil
}
