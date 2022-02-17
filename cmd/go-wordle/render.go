package main

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-wordle/wordle"
)

var (
	// colors
	colorHints              = text.Colors{text.FgHiBlack, text.Italic}
	colorsSpecial           = [3]text.Color{text.FgBlack, text.BgBlack, text.FgHiYellow}
	colorsUnknown           = [3]text.Color{text.FgHiBlack, text.BgHiBlack, text.FgHiWhite}
	colorsNotPresent        = [3]text.Color{text.FgBlack, text.BgBlack, text.FgHiBlack}
	colorsInWrongLocation   = [3]text.Color{text.FgHiYellow, text.BgHiYellow, text.FgBlack}
	colorsInCorrectLocation = [3]text.Color{text.FgHiGreen, text.BgHiGreen, text.FgBlack}

	// misc state
	linesRendered = 0
)

func render(w wordle.Wordle, hints []string, currAttempt wordle.Attempt) {
	for linesRendered > 0 {
		fmt.Print(text.CursorUp.Sprint())
		fmt.Print(text.EraseLine.Sprint())
		linesRendered--
	}

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{renderTitle()})
	tw.AppendRow(table.Row{renderWordle(w, currAttempt)})
	if *flagHints && !w.Solved() {
		tw.AppendFooter(table.Row{renderHints(hints)})
	}
	tw.AppendFooter(table.Row{renderKeyboard(w)})
	tw.AppendFooter(table.Row{renderKeyboardShortcuts()})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter, AlignHeader: text.AlignCenter, AlignFooter: text.AlignCenter},
	})
	tw.SetStyle(table.StyleBold)
	tw.Style().Format.Footer = text.FormatDefault
	tw.Style().Options.SeparateRows = true
	out := tw.Render()
	linesRendered = strings.Count(out, "\n") + 1
	fmt.Println(out)
}

func renderHints(hints []string) string {
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

func renderKeyboard(w wordle.Wordle) string {
	tw := table.NewWriter()

	if *flagHelper && inputModeCharStatus {
		tw.AppendRow(table.Row{renderKeyboardCharacterStatus()})
	} else {
		tw.AppendRow(table.Row{renderKeyboardLegend(w)})
	}

	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderKeyboardCharacterStatus() string {
	tw := table.NewWriter()
	row := table.Row{}

	for _, legend := range []wordle.CharacterStatus{
		wordle.NotPresent,
		wordle.PresentInWrongLocation,
		wordle.PresentInCorrectLocation,
	} {
		colors := colorsUnknown
		switch legend {
		case wordle.PresentInCorrectLocation:
			colors = colorsInCorrectLocation
		case wordle.PresentInWrongLocation:
			colors = colorsInWrongLocation
		case wordle.NotPresent:
			colors = colorsNotPresent
		}

		row = append(row, renderKey(fmt.Sprint(int(legend)), colors))
	}

	tw.AppendRow(row)
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderKeyboardLegend(w wordle.Wordle) string {
	tw := table.NewWriter()
	alphabets := w.Alphabets()

	for _, legend := range []string{"QWERTYUIOP", "ASDFGHJKL", "ZXCVBNM"} {
		twRow := table.NewWriter()
		row := table.Row{}
		if legend == "ZXCVBNM" {
			row = append(row, renderKey("ENTER", colorsSpecial))
		}
		for _, ch := range legend {
			char := string(ch)
			charStatus := alphabets[strings.ToLower(char)]
			colors := colorsUnknown
			switch charStatus {
			case wordle.PresentInCorrectLocation:
				colors = colorsInCorrectLocation
			case wordle.PresentInWrongLocation:
				colors = colorsInWrongLocation
			case wordle.NotPresent:
				colors = colorsNotPresent
			}

			row = append(row, renderKey(char, colors))
		}
		if legend == "ZXCVBNM" {
			row = append(row, renderKey("BKSP", colorsSpecial))
		}
		twRow.AppendRow(row)
		twRow.Style().Options = table.OptionsNoBordersAndSeparators
		tw.AppendRow(table.Row{twRow.Render()})
	}

	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderKeyboardShortcuts() string {
	if *flagDemo {
		return colorHints.Sprint("escape/ctrl+c to quit")
	}
	return colorHints.Sprint("escape/ctrl+c to quit; ctrl+r to restart")
}

func renderTitle() string {
	colors := text.Colors{text.FgHiWhite}

	tw := table.NewWriter()
	tw.AppendRow(table.Row{colors.Sprint("       ▞ ▛▀▀▀▀▀▀▀▀▀▀▀▀▀▜ ▚       ")})
	tw.AppendRow(table.Row{colors.Sprint("░ ▒ ▓ █ ▌  W O R D L E  ▐ █ ▓ ▒ ░")})
	tw.AppendRow(table.Row{colors.Sprint("       ▚ ▙▄▄▄▄▄▄▄▄▄▄▄▄▄▟ ▞       ")})
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderWordle(w wordle.Wordle, currAttempt wordle.Attempt) string {
	tw := table.NewWriter()
	attempts := w.Attempts()
	for attemptIdx := 0; attemptIdx < *flagMaxAttempts; attemptIdx++ {
		var attempt wordle.Attempt
		if attemptIdx < len(attempts) {
			attempt = attempts[attemptIdx]
		} else if attemptIdx == len(attempts) {
			attempt = currAttempt
		}

		tw.AppendRow(table.Row{renderWordleAttempt(attempt)})
	}
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderWordleAttempt(attempt wordle.Attempt) string {
	tw := table.NewWriter()
	twAttemptRow := table.Row{}
	for idx := 0; idx < *flagWordLength; idx++ {
		var char string
		if idx < len(attempt.Answer) {
			char = strings.ToUpper(string(attempt.Answer[idx]))
		} else {
			char = " "
		}
		colors := colorsUnknown
		if idx < len(attempt.Result) {
			switch attempt.Result[idx] {
			case wordle.PresentInCorrectLocation:
				colors = colorsInCorrectLocation
			case wordle.PresentInWrongLocation:
				colors = colorsInWrongLocation
			case wordle.NotPresent:
				colors = colorsNotPresent
			}
		}
		twAttemptRow = append(twAttemptRow, renderKey(char, colors))
	}
	tw.AppendRow(twAttemptRow)
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderKey(key string, colors [3]text.Color) string {
	return fmt.Sprintf("%s\n%s\n%s",
		colors[0].Sprint(strings.Repeat("▄", len(key)+2)),
		colors[1].Sprintf(" %s ", colors[2].Sprint(key)),
		colors[0].Sprint(strings.Repeat("▀", len(key)+2)),
	)
}
