package wordle

import (
	"sort"
	"strings"
)

const (
	englishAlphabets              = "abcdefghijklmnopqrstuvwxyz"
	englishAlphabetsWithoutVowels = "bcdfghjklmnpqrstvwxyz"
)

var (
	maxHints = 5
)

// generateHints produces a list of possible words that will help solve the
// Wordle puzzle.
func generateHints(dictionary []string, attempts []Attempt, alphasStatusMap map[string]CharacterStatus) []string {
	alphasInCorrectLocation := make(map[string]bool)
	alphasInWrongLocation := make(map[string]bool)
	alphasNotPresent := make(map[string]bool)
	alphasPresent := make(map[string]bool)
	alphasUnknown := make(map[string]bool)
	for _, char := range englishAlphabets {
		charStr := string(char)
		switch alphasStatusMap[charStr] {
		case Unknown:
			alphasUnknown[charStr] = true
		case NotPresent:
			alphasNotPresent[charStr] = true
		case WrongLocation:
			alphasInWrongLocation[charStr] = true
			alphasPresent[charStr] = true
		case CorrectLocation:
			alphasInCorrectLocation[charStr] = true
			alphasPresent[charStr] = true
		}
	}

	var words []string
	// if there are too many known and unknown letters, try words with unknown
	// letters to gather more data for next round
	maxWordLength := calculateMaximumWordLength(dictionary)
	maxWordLength75Percent := maxWordLength * 75 / 100
	if len(alphasUnknown) > 20 && len(alphasInCorrectLocation) >= maxWordLength75Percent {
		words = findWordsWithMostUnknownLetters(dictionary, alphasUnknown)
	}

	// if it couldn't find any usable words, proceed with normal logic
	if len(words) == 0 {
		// remove words with letters known to be not present
		words = filterWordsWithLettersNotPresent(dictionary, alphasNotPresent)
		// remove words with characters in wrong places
		words = filterWordsWithLettersInWrongLocations(words, attempts)
		// remove words without letters known to be present
		words = filterWordsWithoutLetters(words, alphasPresent)
	}

	// if there is more than one option, try narrowing it down
	if len(words) > 1 {
		if len(alphasInCorrectLocation) >= maxWordLength75Percent && len(words) >= maxHints-1 {
			// if most letters are in right position, but there are still a lot
			// more words to choose from, try to make words using the missing
			// letters (ex.: cra_e; options=crake|crane|crape|crave|craze;
			// find words with k,n,p,t,v,z like knaps|knave)
			missingLetters := findMissingLetters(words, alphasInCorrectLocation)
			words = findWordsWithMostMissingLetters(dictionary, missingLetters)
		} else if len(alphasInCorrectLocation) < maxWordLength75Percent && len(words) <= maxHints {
			// if few letters are in right position, and there are only a few
			// options to choose from, try to make words using the unique
			// letters in all the words (ex.: _iddy; options=biddy|giddy|kiddy|widdy;
			// find words with b,g,k,w)
			differingLetters := findDifferingLetters(words)
			words = findWordsWithMostMissingLetters(dictionary, differingLetters)
		} else {
			// build a frequency map and sort by it to eliminate most common
			// letters for the next attempt
			freqMap := buildCharacterFrequencyMap(words)
			// sort in descending order of frequency
			sort.SliceStable(words, func(i, j int) bool {
				iFreq := calculateFrequencyValue(words[i], freqMap)
				jFreq := calculateFrequencyValue(words[j], freqMap)
				if iFreq == jFreq {
					return words[i] < words[j] // sort alphabetically
				}
				return iFreq > jFreq
			})
		}
	}

	// filter words that have been attempted already
	words = filterWordsAlreadyAttempted(words, attempts)

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
	hasValue := func(values []int, val int) bool {
		for _, v := range values {
			if val == v {
				return true
			}
		}
		return false
	}
	for _, attempt := range attempts {
		for idx, char := range attempt.Answer {
			charStr := string(char)
			if attempt.Result[idx] == CorrectLocation {
				if correctLocationMap[charStr] == nil {
					correctLocationMap[charStr] = make([]int, 0)
				}
				if !hasValue(correctLocationMap[charStr], idx) {
					correctLocationMap[charStr] = append(correctLocationMap[charStr], idx)
				}
			} else if attempt.Result[idx] == WrongLocation {
				if incorrectLocationMap[charStr] == nil {
					incorrectLocationMap[charStr] = make([]int, 0)
				}
				if !hasValue(incorrectLocationMap[charStr], idx) {
					incorrectLocationMap[charStr] = append(incorrectLocationMap[charStr], idx)
				}
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

func countUniqueLetters(word string) int {
	uniqueLettersMap := make(map[string]bool)
	for _, char := range word {
		uniqueLettersMap[string(char)] = true
	}
	return len(uniqueLettersMap)
}

func filterWordsAlreadyAttempted(words []string, attempts []Attempt) []string {
	var rsp []string
	for _, word := range words {
		found := false
		for _, attempt := range attempts {
			if word == attempt.Answer {
				found = true
				break
			}
		}
		if !found {
			rsp = append(rsp, word)
		}
	}
	return rsp
}

func filterWordsWithLettersInWrongLocations(words []string, attempts []Attempt) []string {
	// build some useful maps
	correctLocationMap, incorrectLocationMap := buildKnownCharacterLocationMap(attempts)
	// function to parse a word based on the maps
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

func filterWordsWithoutLetters(words []string, lettersMap map[string]bool) []string {
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

func findDifferingLetters(words []string) map[string]bool {
	letterCountMap := make(map[string]int)
	for _, word := range words {
		letterFoundMap := make(map[string]bool)
		for _, char := range word {
			charStr := string(char)
			// count each letter only once per word
			if !letterFoundMap[charStr] {
				letterCountMap[charStr]++
				letterFoundMap[charStr] = true
			}
		}
	}

	rspMap := make(map[string]bool)
	for char, count := range letterCountMap {
		if count < 2 {
			rspMap[char] = true
		}
	}
	return rspMap
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
		letterSeen := make(map[string]bool)
		for _, char := range word {
			charStr := string(char)
			// count each letter only once per word
			if !letterSeen[charStr] {
				if lettersMap[charStr] {
					score++
				}
				letterSeen[charStr] = true
			}
		}
		return score
	}

	// sort in descending order of score
	sort.SliceStable(words, func(i, j int) bool {
		iScore := missingLettersScore(words[i])
		jScore := missingLettersScore(words[j])
		if iScore == jScore {
			return words[i] < words[j] // sort alphabetically
		}
		return iScore > jScore
	})
	return words
}

func findWordsWithMostUnknownLetters(words []string, lettersMap map[string]bool) []string {
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
	// sort and move up the words with the maximum number of unique letters
	sort.SliceStable(rsp, func(i, j int) bool {
		iCount := countUniqueLetters(rsp[i])
		jCount := countUniqueLetters(rsp[j])
		if iCount == jCount { // sort alphabetically
			return rsp[i] < rsp[j]
		}
		return iCount > jCount
	})
	return rsp
}
