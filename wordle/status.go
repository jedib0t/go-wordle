package wordle

type CharacterStatus int

const (
	NotPresent CharacterStatus = iota
	Unknown
	PresentInWrongLocation
	PresentInCorrectLocation
)
