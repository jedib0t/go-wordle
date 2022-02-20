package main

import (
	"strings"

	"github.com/jedib0t/go-wordle/wordle"
)

// generateWordles returns a list of wordle.Wordle games with customizations
// defined by flags.
func generateWordles(numWordles int) []wordle.Wordle {
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
		opts := generateWordlesOptions()
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

// generateWordlesOptions returns a list of wordle.Option to construct a Wordle
// game.
func generateWordlesOptions() []wordle.Option {
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

// getAttempt returns the current attempt object being modified. This is more
// specifically needed to move between multiple wordles in mega-wordles mode.
func getAttempt(wordles []wordle.Wordle, currAttempts []wordle.Attempt, direction int) (int, wordle.Attempt) {
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

// isAtLastUnsolvedWordle returns true if the inputCharStatusAttemptIdx is at
// the last unsolved Wordle in the list of Wordles.
func isAtLastUnsolvedWordle(wordles []wordle.Wordle) bool {
	for idx := inputCharStatusAttemptIdx + 1; idx < len(wordles); idx++ {
		if !wordles[idx].GameOver() {
			return false
		}
	}
	return true
}
