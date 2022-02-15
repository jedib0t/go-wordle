package wordle

// Wordle defines methods to interact with a Wordle game.
type Wordle interface {
	Alphabets() map[string]CharacterStatus
	Answer() string
	Attempt(word string, result ...CharacterStatus) (*Attempt, error)
	Attempts() []Attempt
	Dictionary() []string
	DictionaryHas(word string) bool
	Hints() []string
	Reset() error
	Solved() bool
}
