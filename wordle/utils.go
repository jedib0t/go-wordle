package wordle

import "sort"

// CombineHints generates hints from all provided Wordles and combines them in
// descending order of character frequency.
func CombineHints(wordles ...Wordle) []string {
	// if only one wordle is available, just return Hints from it
	if len(wordles) == 1 {
		return wordles[0].Hints()
	}

	// function to tell if a word has been attempted already
	hasBeenAttempted := func(word string) bool {
		for _, w := range wordles {
			for _, attempt := range w.Attempts() {
				if attempt.Answer == word {
					return true
				}
			}
		}
		return false
	}

	// else get all hints and return it after sorting by character frequency
	var hints []string
	for _, w := range wordles {
		for _, hint := range w.Hints() {
			if !hasBeenAttempted(hint) {
				hints = append(hints, hint)
			}
		}
	}
	freqMap := buildCharacterFrequencyMap(hints)
	sort.SliceStable(hints, func(i, j int) bool {
		iFreq := calculateFrequencyValue(hints[i], freqMap)
		jFreq := calculateFrequencyValue(hints[j], freqMap)
		if iFreq == jFreq {
			return hints[i] < hints[j] // sort alphabetically
		}
		return iFreq > jFreq
	})

	if len(hints) > maxHints {
		return hints[:maxHints]
	}
	return hints
}
