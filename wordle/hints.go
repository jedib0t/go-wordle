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

	var words []string
	// if there are too many unknowns, try words with unknown letters to gather
	// more data for next round
	if len(alphasUnknown) > 18 {
		words = filterWordsWithLettersUnknown(dictionary, alphasUnknown)
	}

	// if it couldn't find any usable words, proceed with normal logic
	if len(words) == 0 {
		// build some useful maps
		correctLocationMap, incorrectLocationMap := buildKnownCharacterLocationMap(attempts)
		// remove words with letters known to be not present
		words = filterWordsWithLettersNotPresent(dictionary, alphasNotPresent)
		// remove words with characters in wrong places
		words = filterWordsWithLettersInWrongLocations(words, correctLocationMap, incorrectLocationMap)
		// remove words with characters missing
		words = filterWordsWithLettersMissing(words, alphasPresent)
	}

	// if most letters are in right position, but there are still a lot more
	// words to choose from, try to make words using the missing letters
	// (ex.: cra_e; options=crake|crane|crate|crave|craze; find words with k,n,t,v,z)
	maxWordLength := calculateMaximumWordLength(dictionary)
	if len(alphasInCorrectLocation) >= (maxWordLength*75/100) && len(words) >= maxWordLength-1 {
		words = findWordsWithMostMissingLetters(dictionary, findMissingLetters(words, alphasInCorrectLocation))
	} else {
		// build a frequency map and sort by it
		freqMap := buildCharacterFrequencyMap(words)
		// sort in descending order of frequency
		sort.SliceStable(words, func(i, j int) bool {
			return calculateFrequencyValue(words[i], freqMap) > calculateFrequencyValue(words[j], freqMap)
		})
	}

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

func buildKnownCharacterLocationMap(attempts []Attempt) (map[string][]int, map[string][]int) {
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
	charSeen := make(map[string]bool)
	for _, char := range word {
		charStr := string(char)
		if !charSeen[charStr] {
			val += freqMap[charStr]
			charSeen[charStr] = true
		}
	}
	return val
}

func calculateMaximumWordLength(words []string) int {
	maxWordLength := 0
	for _, word := range words {
		if len(word) > maxWordLength {
			maxWordLength = len(word)
		}
	}
	return maxWordLength
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

func filterWordsWithLettersUnknown(words []string, lettersMap map[string]bool) []string {
	hasAllUnknownLetters := func(word string) bool {
		for _, char := range word {
			if !lettersMap[string(char)] {
				return false
			}
		}
		return true
	}

	var rsp []string
	for _, word := range words {
		if hasAllUnknownLetters(word) {
			rsp = append(rsp, word)
		}
	}
	return rsp
}

func findMissingLetters(words []string, letterMap map[string]bool) map[string]bool {
	rspMap := make(map[string]bool)
	for _, word := range words {
		for _, char := range word {
			if !letterMap[string(char)] {
				rspMap[string(char)] = true
			}
		}
	}
	return rspMap
}

func findWordsWithMostMissingLetters(words []string, lettersMap map[string]bool) []string {
	missingLettersScore := func(word string) int {
		score := 0
		for _, char := range word {
			if lettersMap[string(char)] {
				score++
			}
		}
		return score
	}

	// sort in descending order of score
	sort.SliceStable(words, func(i, j int) bool {
		return missingLettersScore(words[i]) > missingLettersScore(words[j])
	})
	return words
}
