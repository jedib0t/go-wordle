package wordle

// Wordle defines methods to interact with a Wordle game.
type Wordle interface {
	Alphabets() map[string]CharacterStatus
	Answer() string
	AnswerUnknown() bool
	Attempt(word string, result ...CharacterStatus) (*Attempt, error)
	Attempts() []Attempt
	Dictionary() []string
	DictionaryHas(word string) bool
	GameOver() bool
	Hints() []string
	Reset() error
	Solved() bool
}
