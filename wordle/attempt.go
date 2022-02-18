package wordle

import (
	"strings"
)

// Attempt help record a user-attempt of a word, and the results per character.
type Attempt struct {
	Answer string
	Result []CharacterStatus
}

func (a *Attempt) computeResult(answer string) {
	a.Result = make([]CharacterStatus, len(a.Answer))

	// 1st phase: find letters not present and present in correct location
	for idx, char := range a.Answer {
		charStr := string(char)
		if !strings.Contains(answer, charStr) {
			a.Result[idx] = NotPresent
			continue
		}
		if a.Answer[idx] == answer[idx] {
			a.Result[idx] = CorrectLocation
			continue
		}
	}
	// 2nd phase: find letters in wrong locations
	for idx, char := range a.Answer {
		charStr := string(char)
		if a.numberOfFinds(charStr) < strings.Count(answer, charStr) {
			if a.Result[idx] != CorrectLocation {
				a.Result[idx] = WrongLocation
			}
		}
	}
}

func (a *Attempt) numberOfFinds(letter string) int {
	count := 0
	for idx := range a.Answer {
		if string(a.Answer[idx]) == letter {
			if a.Result[idx] == CorrectLocation || a.Result[idx] == WrongLocation {
				count++
			}
		}
	}
	return count
}
