package game

import (
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

	// game state
	wordles      []wordle.Wordle
	currAttempts []wordle.Attempt
	hints        []string
)

// Play starts the game.
func Play() {
	defer cleanup()
	generateWordles(*flagNumWordles)
	cliAttempts := strings.Split(*flagAttempts, ",")
	currAttempts = make([]wordle.Attempt, len(wordles))
	hints = wordle.CombineHints(wordles...)

	// render forever in a separate routine
	chStop := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go renderAsync(chStop, &wg)

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
				renderMutex.Lock()
				_, _ = w.Attempt(cliAttempts[0])
				if idx == len(wordles)-1 { // last wordle
					cliAttempts = cliAttempts[1:]
					hints = wordle.CombineHints(wordles...)
				}
				renderMutex.Unlock()
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
			solve()
		} else {
			getUserInput()
		}
	}

	chStop <- true
	wg.Wait()
}

func getUserInput() {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		logErrorAndExit("failed to get input: %v", err)
	}

	switch key {
	case keyboard.KeyEsc, keyboard.KeyCtrlC:
		handleActionQuit()
	case keyboard.KeyCtrlD:
		handleActionDecrementAttempts()
	case keyboard.KeyCtrlI:
		handleActionIncrementAttempts()
	case keyboard.KeyCtrlR:
		handleActionReset()
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		handleActionBackSpace()
	case keyboard.KeyEnter:
		if inputCharStatus {
			handleActionAttemptStatus()
		} else {
			handleActionAttempt()
		}
	default:
		if inputCharStatus {
			handleActionInputStatus(char)
		} else {
			handleActionInput(char)
		}
	}
}

func solve() {
	if solveWordSet {
		// wait and honor solve speed set in flags
		time.Sleep(time.Second / time.Duration(*flagSolveSpeed))
		// if the word is empty and moved over to the answer, attempt it
		if solveWord == "" {
			renderMutex.Lock()
			for idx, w := range wordles {
				_, _ = w.Attempt(currAttempts[idx].Answer)
			}
			currAttempts = make([]wordle.Attempt, len(wordles))
			hints = wordle.CombineHints(wordles...)
			renderMutex.Unlock()
			solveWordSet = false
			return
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
		renderMutex.Lock()
		for idx := range currAttempts {
			currAttempts[idx].Answer += string(solveWord[0])
		}
		solveWord = solveWord[1:]
		renderMutex.Unlock()
	}
}
