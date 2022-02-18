package main

import "github.com/jedib0t/go-wordle/wordle"

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
