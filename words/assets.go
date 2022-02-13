package words

import (
	_ "embed"
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

// English returns all known English words.
func English() []string {
	englishOnce.Do(func() {
		englishWords = strings.Split(englishTxtRaw, "\n")
		sort.Strings(englishWords)
	})
	return englishWords
}
