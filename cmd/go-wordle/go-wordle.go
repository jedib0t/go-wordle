package main

import (
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/jedib0t/go-wordle/wordle"
)

var (
	// other variables
	inputCharStatus           = false
	inputCharStatusAttemptIdx = 0
	solveWord                 = ""
	solveWordSet              = false

	// version
	version = "dev"
)

func main() {
	defer cleanup()

	// instantiate
	wordles := getWordles(*flagNumWordles)

	// prompt and render
	prompt(wordles)
}

func getUserInput(wordles []wordle.Wordle, currAttempts []wordle.Attempt, hints []string) ([]wordle.Attempt, []string) {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		logErrorAndExit("failed to get input: %v", err)
	}

	switch key {
	case keyboard.KeyEsc, keyboard.KeyCtrlC:
		os.Exit(0)
	case keyboard.KeyCtrlR:
		for _, w := range wordles {
			_ = w.Reset()
		}
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		for attemptIdx := range currAttempts {
			if len(currAttempts[attemptIdx].Answer) > 0 {
				currAttempts[attemptIdx].Answer = currAttempts[attemptIdx].Answer[:len(currAttempts[attemptIdx].Answer)-1]
			}
		}
	case keyboard.KeyEnter:
		for attemptIdx := range currAttempts {
			if len(currAttempts[attemptIdx].Answer) == *flagWordLength {
				if *flagHelper && len(currAttempts[attemptIdx].Result) < len(currAttempts[attemptIdx].Answer) {
					inputCharStatus = true
				} else {
					for idx, w := range wordles {
						_, _ = w.Attempt(currAttempts[idx].Answer)
						if idx == len(wordles)-1 {
							currAttempts = make([]wordle.Attempt, len(wordles))
							hints = wordle.CombineHints(wordles...)
						}
					}
				}
			}
		}
	default:
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			for attemptIdx := range currAttempts {
				if len(currAttempts[attemptIdx].Answer) < *flagWordLength {
					currAttempts[attemptIdx].Answer += strings.ToLower(string(char))
				}
			}
		}
	}
	return currAttempts, hints
}

func getUserInputCharStatus(wordles []wordle.Wordle, currAttempts []wordle.Attempt, hints []string) ([]wordle.Attempt, []string) {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		logErrorAndExit("failed to get input: %v", err)
	}

	// get the current attempt object being modified
	getAttempt := func(direction int) (int, wordle.Attempt) {
		switch direction {
		case -1:
			for idx := inputCharStatusAttemptIdx - 1; idx >= 0; idx-- {
				if !wordles[idx].GameOver() {
					inputCharStatusAttemptIdx = idx
					break
				}
			}
		case 0:
			if wordles[inputCharStatusAttemptIdx].GameOver() {
				for idx := inputCharStatusAttemptIdx + 1; idx < len(wordles); idx++ {
					if !wordles[idx].GameOver() {
						inputCharStatusAttemptIdx = idx
						break
					}
				}
			}
		case 1:
			for idx := inputCharStatusAttemptIdx + 1; idx < len(wordles); idx++ {
				if !wordles[idx].GameOver() {
					inputCharStatusAttemptIdx = idx
					break
				}
			}
		}
		return inputCharStatusAttemptIdx, currAttempts[inputCharStatusAttemptIdx]
	}
	isAtLastUnsolvedWordle := func() bool {
		for idx := inputCharStatusAttemptIdx + 1; idx < len(wordles); idx++ {
			if !wordles[idx].GameOver() {
				return false
			}
		}
		return true
	}
	attemptIdx, currAttempt := getAttempt(0)

	switch key {
	case keyboard.KeyEsc, keyboard.KeyCtrlC:
		os.Exit(0)
	case keyboard.KeyCtrlR:
		for _, w := range wordles {
			_ = w.Reset()
		}
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		if attemptIdx > 0 && len(currAttempt.Result) == 0 {
			attemptIdx, currAttempt = getAttempt(-1)
		}
		if len(currAttempt.Result) > 0 {
			currAttempts[attemptIdx].Result = currAttempt.Result[:len(currAttempt.Result)-1]
		}
	case keyboard.KeyEnter:
		if isAtLastUnsolvedWordle() && len(currAttempt.Result) == len(currAttempt.Answer) {
			inputCharStatus = false
			inputCharStatusAttemptIdx = 0
			for idx, w := range wordles {
				_, _ = w.Attempt(currAttempts[idx].Answer, currAttempts[idx].Result...)
				if idx == len(wordles)-1 {
					currAttempts = make([]wordle.Attempt, len(wordles))
					hints = wordle.CombineHints(wordles...)
				}
			}
		}
	default:
		if char == '0' || char == '2' || char == '3' {
			if len(currAttempt.Result) == len(currAttempt.Answer) {
				if attemptIdx < len(wordles)-1 {
					attemptIdx, currAttempt = getAttempt(+1)
				}
			}
			if len(currAttempt.Result) < len(currAttempt.Answer) {
				switch char {
				case '0':
					currAttempts[attemptIdx].Result = append(currAttempt.Result, wordle.NotPresent)
				case '2':
					currAttempts[attemptIdx].Result = append(currAttempt.Result, wordle.PresentInWrongLocation)
				case '3':
					currAttempts[attemptIdx].Result = append(currAttempt.Result, wordle.PresentInCorrectLocation)
				}
			}
		}
	}
	return currAttempts, hints
}

func getWordles(numWordles int) []wordle.Wordle {
	var answers []string
	if *flagAnswer != "" {
		answers = strings.Split(*flagAnswer, ",")
		if numWordles > 1 && len(answers) > 1 && len(answers) < numWordles {
			logErrorAndExit("game has %d Wordles but only %d answers provided", len(answers), numWordles)
		}
	}

	// instantiate
	var rsp []wordle.Wordle
	for idx := 1; idx <= numWordles; idx++ {
		answer := ""
		opts := getWordlesOptions()
		if len(answers) > 0 {
			answer = answers[idx-1]
			opts = append(opts, wordle.WithAnswer(answer))
		}

		w, err := wordle.New(opts...)
		if err != nil {
			logErrorAndExit("failed to initiate new Wordle: %v, %v", err, opts)
		}
		if *flagSolve && answer != "" && !w.DictionaryHas(answer) {
			logErrorAndExit("solve will fail as '%s' is not in dictionary and will never be found", answer)
		}
		rsp = append(rsp, w)
	}
	return rsp
}

func getWordlesOptions() []wordle.Option {
	opts := []wordle.Option{
		wordle.WithMaxAttempts(*flagMaxAttempts),
		wordle.WithWordFilters(
			wordle.WithLength(*flagWordLength),
		),
	}
	if *flagHelper {
		opts = append(opts, wordle.WithUnknownAnswer(*flagWordLength))
	}
	return opts
}

func prompt(wordles []wordle.Wordle) {
	cliAttempts := strings.Split(*flagAttempts, ",")
	currAttempts := make([]wordle.Attempt, len(wordles))
	solveWord = ""
	hints := wordle.CombineHints(wordles...)
	for {
		// render the wordle(s)
		render(wordles, hints, currAttempts)

		// process attempt
		numGameOver := 0
		for idx, w := range wordles {
			if w.GameOver() {
				numGameOver++
				continue
			}

			// if user provided words to attempt, do that first
			if len(cliAttempts) > 0 {
				_, _ = w.Attempt(cliAttempts[0])
				if idx == len(wordles)-1 { // last wordle
					cliAttempts = cliAttempts[1:]
					hints = wordle.CombineHints(wordles...)
					render(wordles, hints, currAttempts)
				}
				continue
			}
		}
		if numGameOver == len(wordles) {
			break
		}
		if len(cliAttempts) > 0 {
			continue
		}

		// prompt the user for input
		if *flagSolve {
			currAttempts, hints = solveWithHints(wordles, currAttempts, hints)
		} else if *flagHelper && inputCharStatus {
			currAttempts, hints = getUserInputCharStatus(wordles, currAttempts, hints)
		} else {
			currAttempts, hints = getUserInput(wordles, currAttempts, hints)
		}
	}
}

func solveWithHints(wordles []wordle.Wordle, currAttempts []wordle.Attempt, hints []string) ([]wordle.Attempt, []string) {
	if solveWordSet {
		time.Sleep(time.Second / time.Duration(*flagSolveSpeed))
		// if the word is empty and moved over to the answer, attempt it
		if solveWord == "" {
			solveWordSet = false
			for idx, w := range wordles {
				_, _ = w.Attempt(currAttempts[idx].Answer)
			}
			return make([]wordle.Attempt, len(wordles)), wordle.CombineHints(wordles...)
		}
	}
	if solveWord == "" {
		if len(hints) == 0 {
			logErrorAndExit("Uh oh... failed to solve the Wordle! Big Sad!")
		}
		solveWord = hints[0]
		solveWordSet = true
	}
	if len(solveWord) > 0 {
		// move one letter over to the answer
		for idx := range currAttempts {
			currAttempts[idx].Answer += string(solveWord[0])
		}
		solveWord = solveWord[1:]
	}
	return currAttempts, hints
}
