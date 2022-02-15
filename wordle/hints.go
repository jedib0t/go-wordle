package wordle

import (
	"sort"
	"strings"
)

var (
	maxHints = 5
)

// generateHints produces a list of possible words that will help solve the
// Wordle puzzle.
func generateHints(dictionary []string, attempts []Attempt, alphasStatusMap map[string]CharacterStatus) []string {
	alphasInCorrectLocation := make(map[string]bool)
	alphasInIncorrectLocation := make(map[string]bool)
	alphasNotPresent := make(map[string]bool)
	alphasPresent := make(map[string]bool)
	alphasUnknown := make(map[string]bool)
	for _, char := range "abcdefghijklmnopqrstuvwxyz" {
		charStr := string(char)
		switch alphasStatusMap[charStr] {
		case Unknown:
			alphasUnknown[charStr] = true
		case NotPresent:
			alphasNotPresent[charStr] = true
		case PresentInWrongLocation:
			alphasInIncorrectLocation[charStr] = true
			alphasPresent[charStr] = true
		case PresentInCorrectLocation:
			alphasInCorrectLocation[charStr] = true
			alphasPresent[charStr] = true
		}
	}

	// build some useful maps
	correctLocationMap, incorrectLocationMap := buildCharacterLocationMap(attempts)
	// remove words with letters known to be not present
	words := filterWordsWithLettersNotPresent(dictionary, alphasNotPresent)
	// remove words with characters in wrong places
	words = filterWordsWithLettersInWrongLocations(words, correctLocationMap, incorrectLocationMap)
	// remove words with characters missing
	words = filterWordsWithLettersMissing(words, alphasPresent)

	// build a frequency map and sort by it
	freqMap := buildCharacterFrequencyMap(words)
	sort.Slice(words, func(i, j int) bool {
		return calculateFrequencyValue(words[i], freqMap) > calculateFrequencyValue(words[j], freqMap)
	})

	// return the top list
	if len(words) > maxHints {
		return words[:maxHints]
	}
	return words
}

func buildCharacterFrequencyMap(words []string) map[string]int {
	rsp := make(map[string]int)
	for _, word := range words {
		for _, char := range word {
			rsp[string(char)]++
		}
	}
	return rsp
}

func buildCharacterLocationMap(attempts []Attempt) (map[string][]int, map[string][]int) {
	correctLocationMap, incorrectLocationMap := make(map[string][]int), make(map[string][]int)
	for _, attempt := range attempts {
		for idx, char := range attempt.Answer {
			charStr := string(char)
			if attempt.Result[idx] == PresentInCorrectLocation {
				if correctLocationMap[charStr] == nil {
					correctLocationMap[charStr] = make([]int, 0)
				}
				correctLocationMap[charStr] = append(correctLocationMap[charStr], idx)
			} else if attempt.Result[idx] == PresentInWrongLocation {
				if incorrectLocationMap[charStr] == nil {
					incorrectLocationMap[charStr] = make([]int, 0)
				}
				incorrectLocationMap[charStr] = append(incorrectLocationMap[charStr], idx)
			}
		}
	}
	return correctLocationMap, incorrectLocationMap
}

func calculateFrequencyValue(word string, freqMap map[string]int) int {
	val := 0
	for _, char := range word {
		val += freqMap[string(char)]
	}
	return val
}

func filterWordsWithLettersInWrongLocations(words []string, correctLocationMap map[string][]int, incorrectLocationMap map[string][]int) []string {
	hasCharacterInWrongLocation := func(word string) bool {
		for char, indices := range correctLocationMap {
			for _, idx := range indices {
				if string(word[idx]) != char {
					return true
				}
			}
		}
		for char, indices := range incorrectLocationMap {
			for _, idx := range indices {
				if string(word[idx]) == char {
					return true
				}
			}
		}
		return false
	}

	var rsp []string
	for _, word := range words {
		if !hasCharacterInWrongLocation(word) {
			rsp = append(rsp, word)
		}
	}
	return rsp
}

func filterWordsWithLettersMissing(words []string, lettersMap map[string]bool) []string {
	doesNotHaveAllLetters := func(word string) bool {
		for char := range lettersMap {
			if !strings.Contains(word, char) {
				return true
			}
		}
		return false
	}

	var rsp []string
	for _, word := range words {
		if !doesNotHaveAllLetters(word) {
			rsp = append(rsp, word)
		}
	}
	return rsp
}

func filterWordsWithLettersNotPresent(words []string, lettersMap map[string]bool) []string {
	hasLettersNotPresent := func(word string) bool {
		for _, char := range word {
			if lettersMap[string(char)] {
				return true
			}
		}
		return false
	}

	var rsp []string
	for _, word := range words {
		if !hasLettersNotPresent(word) {
			rsp = append(rsp, word)
		}
	}
	return rsp
}

func filterWordsWithMostVowels(words []string) []string {
	countVowels := func(word string) int {
		countMap := make(map[rune]bool)
		for _, char := range word {
			if char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u' {
				countMap[char] = true
			}
		}
		return len(countMap)
	}

	wordVowelCountMap := make(map[int][]string)
	for _, word := range words {
		count := countVowels(word)
		if wordVowelCountMap[count] == nil {
			wordVowelCountMap[count] = make([]string, 0)
		}
		wordVowelCountMap[count] = append(wordVowelCountMap[count], word)
	}
	var rsp []string
	for count := 5; count > 0; count-- {
		for _, word := range wordVowelCountMap[count] {
			rsp = append(rsp, word)
		}
		if len(rsp) > maxHints {
			break
		}
	}
	return rsp[:maxHints]
}
