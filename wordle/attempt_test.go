package wordle

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_computeResult(t *testing.T) {
	compareAttempt := func(t *testing.T, answer string, word string, expectedResults []CharacterStatus) {
		attempt := Attempt{Answer: word}
		attempt.computeResult(answer)

		message := fmt.Sprintf("answer=%s, attempt=%s, expected=%v, result=%v",
			answer, word, expectedResults, attempt.Result)
		assert.Equal(t, expectedResults, attempt.Result, message)
	}

	compareAttempt(t, "antic", "aargh", []CharacterStatus{
		CorrectLocation, NotPresent, NotPresent, NotPresent, NotPresent,
	})
	compareAttempt(t, "antic", "valid", []CharacterStatus{
		NotPresent, WrongLocation, NotPresent, CorrectLocation, NotPresent,
	})
	compareAttempt(t, "aroma", "japan", []CharacterStatus{
		NotPresent, WrongLocation, NotPresent, WrongLocation, NotPresent,
	})
	compareAttempt(t, "aroma", "abate", []CharacterStatus{
		CorrectLocation, NotPresent, WrongLocation, NotPresent, NotPresent,
	})
	compareAttempt(t, "antic", "antic", []CharacterStatus{
		CorrectLocation, CorrectLocation, CorrectLocation, CorrectLocation, CorrectLocation,
	})
}
