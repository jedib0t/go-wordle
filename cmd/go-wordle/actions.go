package main

import (
	"os"
	"strings"

	"github.com/jedib0t/go-wordle/wordle"
)

func handleActionAttempt() {
	renderMutex.Lock()
	defer renderMutex.Unlock()

	for attemptIdx := range currAttempts {
		if len(currAttempts[attemptIdx].Answer) == *flagWordLength {
			if *flagHelper && len(currAttempts[attemptIdx].Result) < len(currAttempts[attemptIdx].Answer) {
				inputCharStatus = true
				continue
			}

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

func handleActionAttemptStatus() {
	renderMutex.Lock()
	defer renderMutex.Unlock()

	if inputCharStatus {
		_, currAttempt := getAttempt(0)
		if isAtLastUnsolvedWordle(wordles) && len(currAttempt.Result) == len(currAttempt.Answer) {
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
	}
}

func handleActionBackSpace() {
	renderMutex.Lock()
	defer renderMutex.Unlock()

	if inputCharStatus {
		attemptIdx, currAttempt := getAttempt(0)
		if attemptIdx > 0 && len(currAttempt.Result) == 0 {
			attemptIdx, currAttempt = getAttempt(-1)
		}
		if len(currAttempt.Result) > 0 {
			currAttempts[attemptIdx].Result = currAttempt.Result[:len(currAttempt.Result)-1]
		}
	} else {
		for attemptIdx := range currAttempts {
			if len(currAttempts[attemptIdx].Answer) > 0 {
				currAttempts[attemptIdx].Answer = currAttempts[attemptIdx].Answer[:len(currAttempts[attemptIdx].Answer)-1]
			}
		}
	}
}

func handleActionDecrementAttempts() {
	renderMutex.Lock()
	defer renderMutex.Unlock()

	canDecrement := *flagMaxAttempts > 1
	for _, w := range wordles {
		if len(w.Attempts()) >= *flagMaxAttempts-1 {
			canDecrement = false
		}
	}
	if canDecrement {
		*flagMaxAttempts--
		for idx, w := range wordles {
			if !w.DecrementMaxAttempts() {
				logErrorAndExit("failed to decrement max attempts for Wordle[%d]", idx)
			}
		}
	}
}

func handleActionIncrementAttempts() {
	renderMutex.Lock()
	defer renderMutex.Unlock()

	*flagMaxAttempts++
	for _, w := range wordles {
		w.IncrementMaxAttempts()
	}
}

func handleActionQuit() {
	cleanup()
	os.Exit(0)
}

func handleActionReset() {
	renderMutex.Lock()
	defer renderMutex.Unlock()

	for _, w := range wordles {
		_ = w.Reset()
	}
}

func handleActionInput(char rune) {
	renderMutex.Lock()
	defer renderMutex.Unlock()

	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
		for attemptIdx := range currAttempts {
			if len(currAttempts[attemptIdx].Answer) < *flagWordLength {
				currAttempts[attemptIdx].Answer += strings.ToLower(string(char))
			}
		}
	}
}

func handleActionInputStatus(char rune) {
	renderMutex.Lock()
	defer renderMutex.Unlock()

	attemptIdx, currAttempt := getAttempt(0)
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
				currAttempts[attemptIdx].Result = append(currAttempt.Result, wordle.WrongLocation)
			case '3':
				currAttempts[attemptIdx].Result = append(currAttempt.Result, wordle.CorrectLocation)
			}
		}
	}
}
