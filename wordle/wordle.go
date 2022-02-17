package wordle

import (
	"fmt"
	"math/rand"
)

type wordle struct {
	alphabets     map[string]CharacterStatus
	answer        string
	answerUnknown bool
	attempts      []Attempt
	dictionary    *[]string
	maxAttempts   int
	options       []Option
	solved        bool
	wordFilters   Filters
	wordsAllowed  []string
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

func (w wordle) AnswerUnknown() bool {
	return w.answerUnknown
}

func (w *wordle) Attempt(word string, result ...CharacterStatus) (*Attempt, error) {
	// do not attempt if game is over
	if w.GameOver() {
		return &w.attempts[len(w.attempts)-1], nil
	}

	// attempt through either available paths
	var attempt *Attempt
	var err error
	if w.answerUnknown {
		attempt, err = w.attemptUnknown(word, result)
	} else {
		attempt, err = w.attemptKnown(word)
	}
	if err != nil {
		return nil, err
	}

	// record the new attempt
	w.attempts = append(w.attempts, *attempt)

	return attempt, nil
}

func (w wordle) Attempts() []Attempt {
	return w.attempts
}

func (w wordle) Dictionary() []string {
	return *w.dictionary
}

func (w wordle) DictionaryHas(word string) bool {
	for _, dictWord := range *w.dictionary {
		if word == dictWord {
			return true
		}
	}
	return false
}

func (w *wordle) GameOver() bool {
	return w.Solved() || len(w.attempts) == w.maxAttempts
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

func (w *wordle) attemptKnown(word string) (*Attempt, error) {
	if err := w.validateAttempt(word); err != nil {
		return nil, err
	}

	// prep to record a new attempt
	attempt := &Attempt{
		Answer: word,
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
					if w.alphabets[string(word[idx])] == Unknown {
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

	// mark as solved if the answer matches
	if attempt.Answer == w.answer {
		w.solved = true
	}

	return attempt, nil
}

func (w *wordle) attemptUnknown(word string, result []CharacterStatus) (*Attempt, error) {
	if err := w.validateAttemptUnknown(word); err != nil {
		return nil, err
	}

	// prep to record a new attempt
	attempt := &Attempt{
		Answer: word,
		Result: result,
	}

	// compute alphabet map from the result
	for idx, status := range attempt.Result {
		charStr := string(word[idx])
		if status != NotPresent {
			w.alphabets[charStr] = status
		}
	}
	// 2nd pass to remove characters not present (2nd pass to support repeated chars)
	for idx, status := range attempt.Result {
		charStr := string(word[idx])
		if status == NotPresent {
			if w.alphabets[charStr] != PresentInWrongLocation && w.alphabets[charStr] != PresentInCorrectLocation {
				delete(w.alphabets, charStr)
			}
		}
	}

	// mark as solved if all characters are in right location
	numCorrect := 0
	for _, status := range attempt.Result {
		if status == PresentInCorrectLocation {
			numCorrect++
		}
	}
	if numCorrect == len(w.answer) {
		w.solved = true
	}

	return attempt, nil
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
	if w.answer == "" && !w.answerUnknown {
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

func (w *wordle) validateAttemptUnknown(word string) error {
	if len(word) != len(w.answer) {
		return fmt.Errorf("word length [%d] does not match answer length [%d]", len(word), len(w.answer))
	}
	for _, attempt := range w.attempts {
		if word == attempt.Answer {
			return fmt.Errorf("word [%s] has been attempted already", word)
		}
	}
	return nil
}
