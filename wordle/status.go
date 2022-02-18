package wordle

// CharacterStatus denotes the status of a character based on the location it
// was attempted at.
type CharacterStatus int

// Available CharacterStatus values.
const (
	NotPresent CharacterStatus = iota
	Unknown
	WrongLocation
	CorrectLocation
)
