package wordle

// Filter returns true if the word passes through the filter. Returns false for
// words that do not qualify.
type Filter func(word string) bool

// WithLength helps filter the words by min/max lengths.
func WithLength(min, max int) Filter {
	return func(word string) bool {
		if len(word) < min || len(word) > max {
			return false
		}
		return true
	}
}

// WithNoRepeatingCharacters prevents words from having repeated characters.
func WithNoRepeatingCharacters() Filter {
	return func(word string) bool {
		runeCount := make(map[rune]int, len(word))
		for _, r := range word {
			runeCount[r]++
			if runeCount[r] > 1 {
				return false
			}
		}
		return true
	}
}

// Filters is a list of filters.
type Filters []Filter

// Apply returns the list of words that are allowed by all filters
func (f Filters) Apply(words *[]string) []string {
	var rsp []string
	for _, word := range *words {
		dnq := false
		for _, filter := range f {
			if !filter(word) {
				dnq = true
				break
			}
		}
		if !dnq {
			rsp = append(rsp, word)
		}
	}
	return rsp
}
