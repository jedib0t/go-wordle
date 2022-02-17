package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/jedib0t/go-wordle/wordle"
)

var (
	// other variables
	demoWord            = ""
	demoWordSet         = false
	inputModeCharStatus = false

	// version
	version = "dev"
)

func main() {
	defer cleanup()

	// prepare the Wordle options
	filters := []wordle.Filter{
		wordle.WithLength(*flagWordLength),
	}
	opts := []wordle.Option{
		wordle.WithAnswer(*flagAnswer),
		wordle.WithMaxAttempts(*flagMaxAttempts),
		wordle.WithWordFilters(filters...),
	}
	if *flagHelper {
		opts = append(opts, wordle.WithUnknownAnswer(*flagWordLength))
	}
	// instantiate
	w, err := wordle.New(opts...)
	if err != nil {
		logErrorAndExit("failed to initiate new Wordle: %v", err)
	}

	// do some sanity checks
	if *flagDemo && *flagAnswer != "" && !w.DictionaryHas(*flagAnswer) {
		logErrorAndExit("demo cannot proceed as '%s' is not in dictionary", *flagAnswer)
	}

	// prompt for user inputs
	prompt(w)
}

func getUserInput(w wordle.Wordle, currAttempt wordle.Attempt, hints []string) (wordle.Attempt, []string) {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		logErrorAndExit("failed to get input: %v", err)
	}

	switch key {
	case keyboard.KeyEsc, keyboard.KeyCtrlC:
		os.Exit(0)
	case keyboard.KeyCtrlR:
		_ = w.Reset()
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		if inputModeCharStatus {
			if len(currAttempt.Result) > 0 {
				currAttempt.Result = currAttempt.Result[:len(currAttempt.Result)-1]
			}
		} else {
			if len(currAttempt.Answer) > 0 {
				currAttempt.Answer = currAttempt.Answer[:len(currAttempt.Answer)-1]
			}
		}
	case keyboard.KeyEnter:
		if len(currAttempt.Answer) == *flagWordLength {
			if *flagHelper && len(currAttempt.Result) < len(currAttempt.Answer) {
				inputModeCharStatus = true
			} else {
				inputModeCharStatus = false
				_, _ = w.Attempt(currAttempt.Answer, currAttempt.Result...)
				hints = w.Hints()
				currAttempt = wordle.Attempt{}
			}
		}
	default:
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			if len(currAttempt.Answer) < *flagWordLength {
				currAttempt.Answer += strings.ToLower(string(char))
			}
		} else if inputModeCharStatus {
			if char == '0' {
				currAttempt.Result = append(currAttempt.Result, wordle.NotPresent)
			} else if char == '2' {
				currAttempt.Result = append(currAttempt.Result, wordle.PresentInWrongLocation)
			} else if char == '3' {
				currAttempt.Result = append(currAttempt.Result, wordle.PresentInCorrectLocation)
			}
		}
	}
	return currAttempt, hints
}

func prompt(w wordle.Wordle) {
	cliAttempts := strings.Split(*flagAttempts, ",")
	currAttempt := wordle.Attempt{}
	demoWord = ""
	hints := w.Hints()
	for {
		render(w, hints, currAttempt)
		if w.Solved() {
			break
		}
		if len(w.Attempts()) == *flagMaxAttempts {
			if !*flagHelper {
				fmt.Printf("Answer: '%v'\n", strings.ToUpper(w.Answer()))
			}
			break
		}

		// if user provided words to attempt, do that first
		if len(cliAttempts) > 0 {
			_, _ = w.Attempt(cliAttempts[0])
			cliAttempts = cliAttempts[1:]
			hints = w.Hints()
			continue
		}

		// prompt the user for input
		if *flagDemo {
			currAttempt, hints = demoSolveWithHints(w, currAttempt, hints)
		} else {
			currAttempt, hints = getUserInput(w, currAttempt, hints)
		}
	}
}

func demoSolveWithHints(w wordle.Wordle, currAttempt wordle.Attempt, hints []string) (wordle.Attempt, []string) {
	if demoWordSet {
		time.Sleep(time.Second / 2)
		// if the word is empty and moved over to the answer, attempt it
		if demoWord == "" {
			_, _ = w.Attempt(currAttempt.Answer)
			demoWordSet = false
			return wordle.Attempt{}, w.Hints()
		}
	}
	if demoWord == "" {
		if len(hints) == 0 {
			logErrorAndExit("Uh oh... failed to solve the Wordle! Big Sad!")
		}
		demoWord = hints[0]
		demoWordSet = true
	}
	if len(demoWord) > 0 {
		// move one letter over to the answer
		currAttempt.Answer += string(demoWord[0])
		demoWord = demoWord[1:]
	}
	return currAttempt, hints
}
