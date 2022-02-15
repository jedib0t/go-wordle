package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateHints_aroma(t *testing.T) {
	filters := Filters{
		WithLength(5, 5),
	}
	dictionary := filters.Apply(&wordsEnglish)
	assert.Contains(t, dictionary, "aroma")
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		alphaStatusMap[string(r)] = Unknown
	}

	hints := generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "arose", hints[0])

	attempts = append(attempts, Attempt{Answer: "arose", Result: []CharacterStatus{3, 3, 3, 0, 0}})
	delete(alphaStatusMap, "s")
	delete(alphaStatusMap, "e")
	hints = generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, 2)
	assert.Equal(t, "aroid", hints[0])
	assert.Equal(t, "aroma", hints[1])
}

func Test_generateHints_crave(t *testing.T) {
	filters := Filters{
		WithLength(5, 5),
	}
	dictionary := filters.Apply(&wordsEnglish)
	assert.Contains(t, dictionary, "crave")
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		alphaStatusMap[string(r)] = Unknown
	}

	hints := generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "arose", hints[0])

	attempts = append(attempts, Attempt{Answer: "arose", Result: []CharacterStatus{2, 3, 0, 0, 3}})
	delete(alphaStatusMap, "o")
	delete(alphaStatusMap, "s")
	alphaStatusMap["a"] = PresentInWrongLocation
	alphaStatusMap["r"] = PresentInCorrectLocation
	alphaStatusMap["e"] = PresentInCorrectLocation
	hints = generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "crate", hints[0])

	attempts = append(attempts, Attempt{Answer: "crate", Result: []CharacterStatus{3, 3, 3, 0, 3}})
	delete(alphaStatusMap, "t")
	alphaStatusMap["c"] = PresentInCorrectLocation
	alphaStatusMap["a"] = PresentInCorrectLocation
	hints = generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "kanzu", hints[0])

	attempts = append(attempts, Attempt{Answer: "kanzu", Result: []CharacterStatus{0, 2, 0, 0, 0}})
	delete(alphaStatusMap, "k")
	delete(alphaStatusMap, "n")
	delete(alphaStatusMap, "z")
	delete(alphaStatusMap, "u")
	hints = generateHints(dictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, 2)
	assert.Equal(t, "crape", hints[0])
	assert.Equal(t, "crave", hints[1])
}

func Test_generateHints_cynic(t *testing.T) {
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
