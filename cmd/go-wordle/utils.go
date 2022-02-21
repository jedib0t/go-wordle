package main

import (
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-wordle/wordle"
)

var (
	colorsAnswerFailed  = [3]text.Colors{{text.FgRed}, {text.BgRed}, {text.FgHiWhite, text.Bold}}
	colorsAnswerHidden  = [3]text.Colors{{text.FgWhite}, {text.BgWhite}, {text.FgBlack, text.Bold}}
	colorsAnswerSuccess = [3]text.Colors{{text.FgGreen}, {text.BgGreen}, {text.FgHiWhite, text.Bold}}
)

// getAnswerLetterAndColors returns the character and colors to be used to
// render a single letter.
func getAnswerLetterAndColors(w wordle.Wordle, answer string, idx int) (string, [3]text.Colors) {
	status := "unsolved"

	// if the letter was identified successfully even once, open it up
	for _, attempt := range w.Attempts() {
		if attempt.Answer[idx] == answer[idx] {
			status = "solved"
		}
	}
	if w.Solved() {
		// if the game is solved, everything can be marked "successful"
		status = "solved"
	} else if w.GameOver() {
		// uh oh, not solved but game over ==> failed
		status = "failed"
	}

	letter := "?"
	colors := colorsAnswerHidden
	switch status {
	case "solved":
		letter = string(answer[idx])
		colors = colorsAnswerSuccess
	case "failed":
		letter = string(answer[idx])
		colors = colorsAnswerFailed
	}
	return letter, colors
}

// isGameOver returns true if all Wordles are done.
func isGameOver(wordles []wordle.Wordle) bool {
	for _, w := range wordles {
		if !w.GameOver() {
			return false
		}
	}
	return true
}
