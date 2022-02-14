package wordle

// Attempt help record a user-attempt of a word, and the results per character.
type Attempt struct {
	Answer string
	Result []CharacterStatus
}
