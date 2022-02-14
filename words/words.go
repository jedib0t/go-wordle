package words

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"sync"
)

//go:embed english.txt
var englishTxtRaw string

var (
	englishWords []string
	englishOnce  sync.Once
)

func init() {
	englishOnce.Do(func() {
		englishWords = strings.Split(englishTxtRaw, "\n")
		sort.Strings(englishWords)
	})
}

// English returns all known English words.
func English() []string {
	return englishWords
}

// EnglishMeaning returns a URL to the meaning/definition of the English word.
func EnglishMeaning(word string) string {
	return fmt.Sprintf("https://www.merriam-webster.com/dictionary/%s", word)
}
