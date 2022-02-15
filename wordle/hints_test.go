package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateHints(t *testing.T) {
	filters := Filters{
		WithLength(5, 5),
	}
	dictionary := filters.Apply(&wordsEnglish)
	assert.Contains(t, dictionary, "cynic")
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		alphaStatusMap[string(r)] = Unknown
	}

	hints := generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "arose", hints[0])

	attempts = append(attempts, Attempt{Answer: "arose", Result: []CharacterStatus{0, 0, 0, 0, 0}})
	delete(alphaStatusMap, "a")
	delete(alphaStatusMap, "r")
	delete(alphaStatusMap, "o")
	delete(alphaStatusMap, "s")
	delete(alphaStatusMap, "e")
	hints = generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "until", hints[0])

	attempts = append(attempts, Attempt{Answer: "until", Result: []CharacterStatus{0, 2, 0, 3, 0}})
	delete(alphaStatusMap, "u")
	delete(alphaStatusMap, "t")
	delete(alphaStatusMap, "l")
	alphaStatusMap["n"] = PresentInWrongLocation
	alphaStatusMap["i"] = PresentInWrongLocation
	hints = generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, 2)
	assert.Equal(t, "cynic", hints[0])
}
