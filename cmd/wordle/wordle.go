package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-wordle/words"
)

var (
	flagSolve = flag.Bool("solve", false, "Solve?")
)

func main() {
	flag.Parse()

	if *flagSolve {
		solve()
	}
	for idx, word := range words.English() {
		fmt.Printf("%d: %v\n", idx, word)
		if idx > 100 {
			break
		}
	}
}

func logErrorAndExit(msg string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "ERRORL: "+strings.TrimSpace(msg)+"\n", a...)
	os.Exit(-1)
}

func solve() {

}
