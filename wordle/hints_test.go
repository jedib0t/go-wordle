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

func computeStatus(answer string, word string) []CharacterStatus {
	rsp := make([]CharacterStatus, len(answer))
	for idx := range word {
		if word[idx] == answer[idx] {
			rsp[idx] = CorrectLocation
		} else {
			foundLetter := false
			for wrongIdx := range answer {
				if word[idx] == answer[wrongIdx] {
					rsp[idx] = WrongLocation
					foundLetter = true
				}
			}
			if !foundLetter {
				rsp[idx] = NotPresent
			}
		}
	}
	return rsp
}

func recordAttempt(alphaStatusMap *map[string]CharacterStatus, attempts *[]Attempt, answer string, result []CharacterStatus) {
	for idx, charStatus := range result {
		charStr := string(answer[idx])
		switch charStatus {
		case NotPresent:
			delete(*alphaStatusMap, charStr)
		case WrongLocation:
			if (*alphaStatusMap)[charStr] != CorrectLocation {
				(*alphaStatusMap)[charStr] = WrongLocation
			}
		case CorrectLocation:
			(*alphaStatusMap)[charStr] = CorrectLocation
		}
	}
	*attempts = append(*attempts, Attempt{Answer: answer, Result: result})
}

func assertHintingEfficiency(t *testing.T, answer string, expectedAttempts int) {
	filters := Filters{
		WithLength(len(answer)),
	}
	dictionary := filters.Apply(&wordsEnglish)

	assert.Contains(t, dictionary, answer)
	attempts := make([]Attempt, 0)
	alphaStatusMap := make(map[string]CharacterStatus)
	for _, r := range englishAlphabets {
		alphaStatusMap[string(r)] = Unknown
	}

	t.Logf("%s: expecting result in/under %d attempts", answer, expectedAttempts)
	for attempt := 1; attempt <= expectedAttempts; attempt++ {
		hints := generateHints(dictionary, attempts, alphaStatusMap)
		if len(hints) == 0 {
			break
		}
		t.Logf("%s: attempt #%d: hints[0] == '%s'", answer, attempt, hints[0])
		if hints[0] == answer { // success
			return
		}
		recordAttempt(&alphaStatusMap, &attempts, hints[0], computeStatus(answer, hints[0]))
	}
	t.Errorf("%s: failed after %d attempts", answer, expectedAttempts)
}

func Test_generateHints(t *testing.T) {
	assertHintingEfficiency(t, "aroma", 3)
	assertHintingEfficiency(t, "crave", 4)
	assertHintingEfficiency(t, "cynic", 4)
	assertHintingEfficiency(t, "softy", 5)
	assertHintingEfficiency(t, "widdy", 5)
	assertHintingEfficiency(t, "wists", 6)
	assertHintingEfficiency(t, "aardvark", 3)
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
