package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-wordle/wordle"
)

var (
	// colors
	colorHints        = text.Colors{text.FgHiBlack, text.Italic}
	colorsSpecialKeys = [3]text.Colors{{text.FgBlack}, {text.BgBlack}, {text.FgHiYellow}} // Enter/BkSp/etc.
	colorsStatusMap   = map[wordle.CharacterStatus][3]text.Colors{
		wordle.Unknown:         {{text.FgHiBlack}, {text.BgHiBlack}, {text.FgHiWhite}},
		wordle.NotPresent:      {{text.FgBlack}, {text.BgBlack}, {text.FgHiBlack}},
		wordle.WrongLocation:   {{text.FgHiYellow}, {text.BgHiYellow}, {text.FgBlack}},
		wordle.CorrectLocation: {{text.FgHiGreen}, {text.BgHiGreen}, {text.FgBlack}},
	}

	// keyboard
	keyboardRows = [][]string{
		{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
		{"a", "s", "d", "f", "g", "h", "j", "k", "l"},
		{"enter", "z", "x", "c", "v", "b", "n", "m", "bksp"},
	}

	// misc
	linesRendered = 0

	// controls
	renderEnabled = true
	renderMutex   = sync.Mutex{}
)

func render(wordles []wordle.Wordle, hints []string, currAttempts []wordle.Attempt) {
	renderMutex.Lock()
	defer renderMutex.Unlock()
	if !renderEnabled {
		return
	}

	for linesRendered > 0 {
		fmt.Print(text.CursorUp.Sprint())
		fmt.Print(text.EraseLine.Sprint())
		linesRendered--
	}

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{renderTitle()})
	tw.AppendRow(table.Row{renderWordle(wordles, currAttempts)})
	if *flagHints {
		tw.AppendFooter(table.Row{renderHints(wordles, hints)})
	}
	tw.AppendFooter(table.Row{renderKeyboard(wordles)})
	tw.AppendFooter(table.Row{renderKeyboardShortcuts()})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter, AlignHeader: text.AlignCenter, AlignFooter: text.AlignCenter},
	})
	tw.SetStyle(table.StyleBold)
	tw.Style().Format.Header = text.FormatDefault
	tw.Style().Format.Footer = text.FormatDefault
	tw.Style().Options.SeparateRows = true
	out := tw.Render()
	linesRendered = strings.Count(out, "\n") + 1
	fmt.Println(out)
}

func renderHints(wordles []wordle.Wordle, hints []string) string {
	if isGameOver(wordles) {
		return colorHints.Sprint("- game over -")
	}
	if len(hints) == 0 {
		return colorHints.Sprint("- no hints found -")
	}

	tw := table.NewWriter()
	twRow := table.Row{}
	for idx, word := range hints {
		if idx%5 == 0 && len(twRow) > 0 {
			tw.AppendRow(twRow)
			twRow = table.Row{}
		}
		twRow = append(twRow, colorHints.Sprintf(word))
	}
	if len(twRow) > 0 {
		tw.AppendRow(twRow)
	}
	tw.SetStyle(table.StyleLight)
	tw.Style().Options.DrawBorder = false
	tw.Style().Options.SeparateRows = true
	return tw.Render()
}

func renderKey(key string, colors [3]text.Colors) string {
	colorBg1 := colors[0]
	colorBg2 := colors[1]
	colorLetter := colors[2]
	key = strings.ToUpper(key)

	return fmt.Sprintf("%s\n%s\n%s",
		colorBg1.Sprint(strings.Repeat("▄", len(key)+2)),
		colorBg2.Sprintf(" %s ", colorLetter.Sprint(strings.ToUpper(key))),
		colorBg1.Sprint(strings.Repeat("▀", len(key)+2)),
	)
}

func renderKeyboard(wordles []wordle.Wordle) string {
	tw := table.NewWriter()

	if *flagHelper && inputCharStatus {
		tw.AppendRow(table.Row{renderKeyboardCharacterStatus()})
	} else {
		tw.AppendRow(table.Row{renderKeyboardLegend(wordles)})
	}

	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderKeyboardCharacterStatus() string {
	tw := table.NewWriter()
	row := table.Row{}

	for _, legend := range []wordle.CharacterStatus{
		wordle.NotPresent,
		wordle.WrongLocation,
		wordle.CorrectLocation,
	} {
		colors := colorsStatusMap[legend]
		row = append(row, renderKey(fmt.Sprint(int(legend)), colors))
	}

	tw.AppendRow(row)
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderKeyboardLegend(wordles []wordle.Wordle) string {
	tw := table.NewWriter()

	// alphabet coloring applies only to single Wordle instance or for cases
	// where some alphabets have same status across all Wordles
	alphabets := make(map[string]wordle.CharacterStatus, 26)
	if len(wordles) == 1 {
		alphabets = wordles[0].Alphabets()
	} else {
		// init map with Unknown status
		for _, legend := range keyboardRows {
			for _, char := range legend {
				alphabets[char] = wordle.Unknown
			}
		}
		// if any keys have same status across all Wordles, use it
		for _, legend := range keyboardRows {
			for _, char := range legend {
				charStatuses := make(map[wordle.CharacterStatus]bool)
				for _, w := range wordles {
					charStatuses[w.Alphabets()[char]] = true
				}
				if len(charStatuses) == 1 {
					for k := range charStatuses {
						alphabets[char] = k
					}
				}
			}
		}
	}

	for _, legend := range keyboardRows {
		twRow := table.NewWriter()
		row := table.Row{}
		for _, char := range legend {
			charStatus := alphabets[char]
			colors := colorsStatusMap[charStatus]
			if len(char) > 1 {
				colors = colorsSpecialKeys
			}

			row = append(row, renderKey(char, colors))
		}
		twRow.AppendRow(row)
		twRow.Style().Options = table.OptionsNoBordersAndSeparators
		tw.AppendRow(table.Row{twRow.Render()})
	}

	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderKeyboardShortcuts() string {
	shortcuts := "escape/ctrl+c to quit; ctrl+r to restart"
	if *flagSolve {
		shortcuts = "escape/ctrl+c to quit"
	}
	shortcuts = text.AlignCenter.Apply(shortcuts, 56)
	return colorHints.Sprint(shortcuts)
}

func renderWordle(wordles []wordle.Wordle, currAttempts []wordle.Attempt) string {
	tw := table.NewWriter()
	twRow := table.Row{}
	for idx, w := range wordles {
		twWordle := table.NewWriter()
		attempts := w.Attempts()
		for attemptIdx := 0; attemptIdx < *flagMaxAttempts; attemptIdx++ {
			var attempt wordle.Attempt
			if attemptIdx < len(attempts) {
				attempt = attempts[attemptIdx]
			} else if attemptIdx == len(attempts) && !w.GameOver() {
				attempt = currAttempts[idx]
			} else if w.GameOver() {
				attempt = wordle.Attempt{
					Answer: strings.Repeat(" ", *flagWordLength),
					Result: make([]wordle.CharacterStatus, *flagWordLength),
				}
			}

			twWordle.AppendRow(table.Row{renderWordleAttempt(w, attempt, false)})
		}

		// append the answer in a new row and hide it until game is over
		if !w.AnswerUnknown() {
			twWordle.AppendSeparator()
			answerAttempt := wordle.Attempt{Answer: w.Answer(), Result: []wordle.CharacterStatus{3, 3, 3, 3, 3}}
			twWordle.AppendRow(table.Row{renderWordleAttempt(w, answerAttempt, true)})
		}

		twWordle.SetStyle(table.StyleLight)
		twWordle.Style().Options = table.OptionsNoBordersAndSeparators
		twRow = append(twRow, twWordle.Render())
	}
	tw.AppendRow(twRow)
	tw.SetStyle(table.StyleLight)
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	tw.Style().Options.SeparateColumns = true
	return tw.Render()
}

func renderWordleAttempt(w wordle.Wordle, attempt wordle.Attempt, isAnswer bool) string {
	tw := table.NewWriter()
	twAttemptRow := table.Row{}
	for idx := 0; idx < *flagWordLength; idx++ {
		var letter string
		if idx < len(attempt.Answer) {
			letter = strings.ToUpper(string(attempt.Answer[idx]))
		} else {
			letter = " "
		}
		colors := colorsStatusMap[wordle.Unknown]
		if idx < len(attempt.Result) {
			colors = colorsStatusMap[attempt.Result[idx]]
		}
		if isAnswer {
			letter, colors = getAnswerLetterAndColors(w, attempt.Answer, idx)
		}
		twAttemptRow = append(twAttemptRow, renderKey(letter, colors))
	}
	tw.AppendRow(twAttemptRow)
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}
