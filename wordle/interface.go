package wordle

type Wordle interface {
	Alphabets() map[string]CharacterStatus
	Answer() string
	Attempt(word string) (*Attempt, error)
	Attempts() []Attempt
	Reset() error
	Solved() bool
}
