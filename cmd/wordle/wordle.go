package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-wordle/wordle"
	"github.com/jedib0t/go-wordle/words"
)

var (
	flagSolve = flag.Bool("solve", false, "Solve?")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	if *flagSolve {
		solve()
	}
	prompt()
}

func logErrorAndExit(msg string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "ERRORL: "+strings.TrimSpace(msg)+"\n", a...)
	os.Exit(-1)
}

func prompt() {
	w, err := wordle.New(wordle.WithAnswer("robin"))
	if err != nil {
		logErrorAndExit("failed to initiate new Wordle: %v", err)
	}
	fmt.Printf("\"%v\" -- %s\n", w.Answer(), words.EnglishMeaning(w.Answer()))

	for _, word := range []string{"adieu", "disco", "zombi", "robin"} {
		attempt, err := w.Attempt(word)
		fmt.Printf("%v, %v, %v\n", attempt, err, w.Solved())
	}
}

func render() {

}

func solve() {

}
