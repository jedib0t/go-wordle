package wordle

// Wordle defines methods to interact with a Wordle game.
type Wordle interface {
	Alphabets() map[string]CharacterStatus
	Answer() string
	Attempt(word string) (*Attempt, error)
	Attempts() []Attempt
	Reset() error
	Solved() bool
}
