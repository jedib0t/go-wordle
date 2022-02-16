package wordle

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testFilters = Filters{
		WithLength(5),
	}
	testDictionary = testFilters.Apply(&wordsEnglish)
)

func recordAttempt(alphaStatusMap *map[string]CharacterStatus, attempts *[]Attempt, answer string, result []CharacterStatus) {
	for idx, charStatus := range result {
		charStr := string(answer[idx])
		switch charStatus {
		case NotPresent:
			delete(*alphaStatusMap, charStr)
		case PresentInWrongLocation:
			if (*alphaStatusMap)[charStr] != PresentInCorrectLocation {
				(*alphaStatusMap)[charStr] = PresentInWrongLocation
			}
		case PresentInCorrectLocation:
			(*alphaStatusMap)[charStr] = PresentInCorrectLocation
		}
	}
	*attempts = append(*attempts, Attempt{Answer: answer, Result: result})
}

func Test_generateHints_aroma(t *testing.T) {
	assert.Contains(t, testDictionary, "aroma")
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range englishAlphabets {
		alphaStatusMap[string(r)] = Unknown
	}

	hints := generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "arose", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "arose", []CharacterStatus{3, 3, 3, 0, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "bight", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "bight", []CharacterStatus{0, 0, 0, 0, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, 1)
	assert.Equal(t, "aroma", hints[0])
}

func Test_generateHints_crave(t *testing.T) {
	assert.Contains(t, testDictionary, "crave")
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range englishAlphabets {
		alphaStatusMap[string(r)] = Unknown
	}

	hints := generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "arose", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "arose", []CharacterStatus{2, 3, 0, 0, 3})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "crate", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "crate", []CharacterStatus{3, 3, 3, 0, 3})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "knaps", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "knaps", []CharacterStatus{0, 0, 3, 0, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, 2)
	assert.Equal(t, "crave", hints[0])
	assert.Equal(t, "craze", hints[1])
}

func Test_generateHints_cynic(t *testing.T) {
	assert.Contains(t, testDictionary, "cynic")
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range englishAlphabets {
		alphaStatusMap[string(r)] = Unknown
	}

	hints := generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "arose", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "arose", []CharacterStatus{0, 0, 0, 0, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "unity", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "unity", []CharacterStatus{0, 2, 2, 0, 2})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "calve", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "calve", []CharacterStatus{3, 0, 0, 0, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, 1)
	assert.Equal(t, "cynic", hints[0])
}

func Test_generateHints_widdy(t *testing.T) {
	assert.Contains(t, testDictionary, "widdy")
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range englishAlphabets {
		alphaStatusMap[string(r)] = Unknown
	}

	hints := generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "arose", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "arose", []CharacterStatus{0, 0, 0, 0, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "unity", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "unity", []CharacterStatus{0, 0, 2, 0, 3})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "dimly", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "dimly", []CharacterStatus{2, 3, 0, 0, 3})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "bewig", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "bewig", []CharacterStatus{0, 0, 2, 2, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, 1)
	assert.Equal(t, "widdy", hints[0])
}

func Test_generateHints_wists(t *testing.T) {
	assert.Contains(t, testDictionary, "wists")
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range englishAlphabets {
		alphaStatusMap[string(r)] = Unknown
	}

	hints := generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "arose", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "arose", []CharacterStatus{0, 0, 0, 2, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "suint", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "suint", []CharacterStatus{2, 0, 2, 0, 2})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "kilts", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "kilts", []CharacterStatus{0, 3, 0, 3, 3})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "chimb", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "chimb", []CharacterStatus{0, 0, 2, 0, 0})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, maxHints)
	assert.Equal(t, "aglow", hints[0])

	recordAttempt(&alphaStatusMap, &attempts, "aglow", []CharacterStatus{0, 0, 0, 0, 2})
	hints = generateHints(testDictionary, attempts, alphaStatusMap)
	assert.Len(t, hints, 1)
	assert.Equal(t, "wists", hints[0])
}

func Test_buildCharacterFrequencyMap(t *testing.T) {
	freqMap := buildCharacterFrequencyMap([]string{
		"aroma",
		"crave",
		"wists",
	})
	assert.Len(t, freqMap, 11)
	assert.Equal(t, 3, freqMap["a"])
	assert.Equal(t, 1, freqMap["c"])
	assert.Equal(t, 1, freqMap["e"])
	assert.Equal(t, 1, freqMap["i"])
	assert.Equal(t, 1, freqMap["m"])
	assert.Equal(t, 1, freqMap["o"])
	assert.Equal(t, 2, freqMap["r"])
	assert.Equal(t, 2, freqMap["s"])
	assert.Equal(t, 1, freqMap["t"])
	assert.Equal(t, 1, freqMap["v"])
	assert.Equal(t, 1, freqMap["w"])
}

func Test_buildKnownCharacterLocationMap(t *testing.T) {
	correctLocationMap, incorrectLocationMap := buildKnownCharacterLocationMap([]Attempt{
		{Answer: "arose", Result: []CharacterStatus{0, 0, 0, 0, 0}},
		{Answer: "unity", Result: []CharacterStatus{0, 0, 2, 0, 3}},
		{Answer: "dimly", Result: []CharacterStatus{2, 3, 0, 0, 3}},
	})
	assert.Len(t, correctLocationMap, 2)
	assert.Equal(t, []int{4}, correctLocationMap["y"])
	assert.Equal(t, []int{1}, correctLocationMap["i"])
	assert.Len(t, incorrectLocationMap, 2)
	assert.Equal(t, []int{0}, incorrectLocationMap["d"])
	assert.Equal(t, []int{2}, incorrectLocationMap["i"])
}

func Test_calculateFrequencyValue(t *testing.T) {
	freqMap := buildCharacterFrequencyMap([]string{
		"aroma",
		"crave",
		"wists",
	})
	assert.Len(t, freqMap, 11)

	assert.Equal(t, 7, calculateFrequencyValue("aroma", freqMap))
	assert.Equal(t, 8, calculateFrequencyValue("crave", freqMap))
	assert.Equal(t, 5, calculateFrequencyValue("wists", freqMap))
	assert.Equal(t, 0, calculateFrequencyValue("zzzzz", freqMap))
}

func Test_calculateMaximumWordLength(t *testing.T) {
	maxWordLen := calculateMaximumWordLength([]string{
		"aroma",
		"brazen",
		"redemption",
		"zzyxx",
	})
	assert.Equal(t, 10, maxWordLen)
}

func Test_countUniqueLetters(t *testing.T) {
	assert.Equal(t, 4, countUniqueLetters("aroma"))
	assert.Equal(t, 5, countUniqueLetters("crave"))
	assert.Equal(t, 4, countUniqueLetters("drama"))
	assert.Equal(t, 3, countUniqueLetters("momma"))
}

func Test_filterWordsWithLettersInWrongLocations(t *testing.T) {
	words := filterWordsWithLettersInWrongLocations(testDictionary, []Attempt{
		{Answer: "arose", Result: []CharacterStatus{0, 0, 0, 0, 0}},
		{Answer: "unity", Result: []CharacterStatus{0, 0, 2, 0, 3}},
		{Answer: "dimly", Result: []CharacterStatus{2, 3, 0, 0, 3}},
	})
	assert.True(t, len(words) > 1)
	for _, word := range words {
		assert.Equal(t, "i", string(word[1]), word)
		assert.Equal(t, "y", string(word[4]), word)
		assert.NotEqual(t, "d", string(word[0]), word)
		assert.NotEqual(t, "i", string(word[2]), word)
	}
}

func Test_filterWordsWithLettersNotPresent(t *testing.T) {
	words := filterWordsWithLettersNotPresent(testDictionary,
		map[string]bool{"a": true, "e": true, "i": true, "o": true, "u": true},
	)
	assert.True(t, len(words) > 1)
	for _, word := range words {
		assert.NotContains(t, word, "a")
		assert.NotContains(t, word, "e")
		assert.NotContains(t, word, "i")
		assert.NotContains(t, word, "o")
		assert.NotContains(t, word, "u")
	}
}

func Test_filterWordsWithoutLetters(t *testing.T) {
	words := filterWordsWithoutLetters(testDictionary,
		map[string]bool{"g": true, "m": true, "p": true, "y": true},
	)
	assert.True(t, len(words) > 0)
	for _, word := range words {
		assert.True(t, strings.Contains(word, "g"), word)
		assert.True(t, strings.Contains(word, "m"), word)
		assert.True(t, strings.Contains(word, "p"), word)
		assert.True(t, strings.Contains(word, "y"), word)
	}
}

func Test_findDifferingLetters(t *testing.T) {
	lettersMap := findDifferingLetters([]string{
		"biddy",
		"giddy",
		"kiddy",
		"widdy",
	})
	assert.Len(t, lettersMap, 4)
	assert.Contains(t, lettersMap, "b")
	assert.Contains(t, lettersMap, "g")
	assert.Contains(t, lettersMap, "k")
	assert.Contains(t, lettersMap, "w")
}

func Test_findMissingLetters(t *testing.T) {
	lettersMap := findMissingLetters(
		[]string{"crake", "crane", "crape", "crave", "craze"},
		map[string]bool{"c": true, "r": true, "a": true, "e": true},
	)
	assert.Len(t, lettersMap, 5)
	assert.Contains(t, lettersMap, "k")
	assert.Contains(t, lettersMap, "n")
	assert.Contains(t, lettersMap, "p")
	assert.Contains(t, lettersMap, "v")
	assert.Contains(t, lettersMap, "z")
}

func Test_findWordsWithMostMissingLetters(t *testing.T) {
	words := findWordsWithMostMissingLetters(testDictionary,
		map[string]bool{"k": true, "n": true, "p": true, "v": true, "z": true},
	)
	assert.True(t, len(words) >= 3)
	assert.Equal(t, "knaps", words[0])
	assert.Equal(t, "knave", words[1])
	assert.Equal(t, "knops", words[2])
}

func Test_findWordsWithMostUnknownLetters(t *testing.T) {
	alphasUnknown := make(map[string]bool)
	for _, char := range englishAlphabetsWithoutVowels { // no vowels
		alphasUnknown[string(char)] = true
	}
	words := findWordsWithMostUnknownLetters(testDictionary, alphasUnknown)
	assert.True(t, len(words) > 1)
	for _, word := range words {
		assert.NotContains(t, word, "a")
		assert.NotContains(t, word, "e")
		assert.NotContains(t, word, "i")
		assert.NotContains(t, word, "o")
		assert.NotContains(t, word, "u")
	}
}
