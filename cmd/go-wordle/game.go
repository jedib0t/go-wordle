package main

import (
	"os"
	"strings"
	"sync"
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
)

func getUserInput(wordles []wordle.Wordle, currAttempts []wordle.Attempt, hints []string) ([]wordle.Attempt, []string) {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		logErrorAndExit("failed to get input: %v", err)
	}

	switch key {
	case keyboard.KeyEsc, keyboard.KeyCtrlC:
		handleShortcutQuit()
	case keyboard.KeyCtrlD:
		handleShortcutDecrementAttempts(wordles)
	case keyboard.KeyCtrlI:
		handleShortcutIncrementAttempts(wordles)
	case keyboard.KeyCtrlR:
		handleShortcutReset(wordles)
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		for attemptIdx := range currAttempts {
			if len(currAttempts[attemptIdx].Answer) > 0 {
				renderMutex.Lock()
				currAttempts[attemptIdx].Answer = currAttempts[attemptIdx].Answer[:len(currAttempts[attemptIdx].Answer)-1]
				renderMutex.Unlock()
			}
		}
	case keyboard.KeyEnter:
		for attemptIdx := range currAttempts {
			if len(currAttempts[attemptIdx].Answer) == *flagWordLength {
				if *flagHelper && len(currAttempts[attemptIdx].Result) < len(currAttempts[attemptIdx].Answer) {
					inputCharStatus = true
				} else {
					renderMutex.Lock()
					for idx, w := range wordles {
						_, _ = w.Attempt(currAttempts[idx].Answer)
						if idx == len(wordles)-1 {
							currAttempts = make([]wordle.Attempt, len(wordles))
							hints = wordle.CombineHints(wordles...)
						}
					}
					renderMutex.Unlock()
				}
			}
		}
	default:
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			for attemptIdx := range currAttempts {
				if len(currAttempts[attemptIdx].Answer) < *flagWordLength {
					renderMutex.Lock()
					currAttempts[attemptIdx].Answer += strings.ToLower(string(char))
					renderMutex.Unlock()
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

	attemptIdx, currAttempt := getAttempt(wordles, currAttempts, 0)
	switch key {
	case keyboard.KeyEsc, keyboard.KeyCtrlC:
		handleShortcutQuit()
	case keyboard.KeyCtrlD:
		handleShortcutDecrementAttempts(wordles)
	case keyboard.KeyCtrlI:
		handleShortcutIncrementAttempts(wordles)
	case keyboard.KeyCtrlR:
		handleShortcutReset(wordles)
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		if attemptIdx > 0 && len(currAttempt.Result) == 0 {
			attemptIdx, currAttempt = getAttempt(wordles, currAttempts, -1)
		}
		if len(currAttempt.Result) > 0 {
			renderMutex.Lock()
			currAttempts[attemptIdx].Result = currAttempt.Result[:len(currAttempt.Result)-1]
			renderMutex.Unlock()
		}
	case keyboard.KeyEnter:
		if isAtLastUnsolvedWordle(wordles) && len(currAttempt.Result) == len(currAttempt.Answer) {
			inputCharStatus = false
			inputCharStatusAttemptIdx = 0
			renderMutex.Lock()
			for idx, w := range wordles {
				_, _ = w.Attempt(currAttempts[idx].Answer, currAttempts[idx].Result...)
				if idx == len(wordles)-1 {
					currAttempts = make([]wordle.Attempt, len(wordles))
					hints = wordle.CombineHints(wordles...)
				}
			}
			renderMutex.Unlock()
		}
	default:
		if char == '0' || char == '2' || char == '3' {
			if len(currAttempt.Result) == len(currAttempt.Answer) {
				if attemptIdx < len(wordles)-1 {
					attemptIdx, currAttempt = getAttempt(wordles, currAttempts, +1)
				}
			}
			if len(currAttempt.Result) < len(currAttempt.Answer) {
				renderMutex.Lock()
				switch char {
				case '0':
					currAttempts[attemptIdx].Result = append(currAttempt.Result, wordle.NotPresent)
				case '2':
					currAttempts[attemptIdx].Result = append(currAttempt.Result, wordle.WrongLocation)
				case '3':
					currAttempts[attemptIdx].Result = append(currAttempt.Result, wordle.CorrectLocation)
				}
				renderMutex.Unlock()
			}
		}
	}
	return currAttempts, hints
}

func handleShortcutDecrementAttempts(wordles []wordle.Wordle) {
	canDecrement := *flagMaxAttempts > 1
	for _, w := range wordles {
		if len(w.Attempts()) >= *flagMaxAttempts-1 {
			canDecrement = false
		}
	}
	if canDecrement {
		renderMutex.Lock()
		*flagMaxAttempts--
		for idx, w := range wordles {
			if !w.DecrementMaxAttempts() {
				logErrorAndExit("failed to decrement max attempts for Wordle[%d]", idx)
			}
		}
		renderMutex.Unlock()
	}
}

func handleShortcutQuit() {
	cleanup()
	os.Exit(0)
}

func handleShortcutReset(wordles []wordle.Wordle) {
	renderMutex.Lock()
	for _, w := range wordles {
		_ = w.Reset()
	}
	renderMutex.Unlock()
}

func handleShortcutIncrementAttempts(wordles []wordle.Wordle) {
	renderMutex.Lock()
	*flagMaxAttempts++
	for _, w := range wordles {
		w.IncrementMaxAttempts()
	}
	renderMutex.Unlock()
}

func play(wordles []wordle.Wordle) {
	cliAttempts := strings.Split(*flagAttempts, ",")
	currAttempts := make([]wordle.Attempt, len(wordles))
	hints := wordle.CombineHints(wordles...)

	// render forever in a separate routine
	chStop := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		timer := time.Tick(time.Second / 10)
		for {
			select {
			case <-chStop: // render one final time and return
				render(wordles, hints, currAttempts)
				return
			case <-timer: // render as part of regular cycle
				render(wordles, hints, currAttempts)
			}
		}
	}()

	// loop until game is over
	for {
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

	chStop <- true
	wg.Wait()
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
