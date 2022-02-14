package wordle

type Wordle interface {
	Alphabets() map[rune]CharacterStatus
	Answer() string
	Attempt(word string) (*Attempt, error)
	Attempts() []Attempt
	Reset() error
	Solved() bool
}
