package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-wordle/wordle"
)

var (
	flagAnswer      = flag.String("answer", "", "Pre-set answer if you don't want a random word")
	flagAttempts    = flag.String("attempts", "", "Words to attempt before prompting user")
	flagHelp        = flag.Bool("help", false, "Show this help-text?")
	flagMaxAttempts = flag.Int("max-attempts", 6, "Maximum attempts allowed")
	flagMaxLength   = flag.Int("max-length", 5, "Maximum word length")
	flagMinLength   = flag.Int("min-length", 5, "Minimum word length")

	colorHints              = text.Colors{text.FgHiBlack, text.Italic}
	colorsSpecial           = [3]text.Color{text.FgBlack, text.BgBlack, text.FgHiYellow}
	colorsUnknown           = [3]text.Color{text.FgHiBlack, text.BgHiBlack, text.FgHiWhite}
	colorsNotPresent        = [3]text.Color{text.FgBlack, text.BgBlack, text.FgHiBlack}
	colorsInWrongLocation   = [3]text.Color{text.FgHiYellow, text.BgHiYellow, text.FgBlack}
	colorsInCorrectLocation = [3]text.Color{text.FgHiGreen, text.BgHiGreen, text.FgBlack}
	linesRendered           = 0
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	if *flagHelp {
		printHelp()
	}
	if *flagMinLength > *flagMaxLength {
		logErrorAndExit("min-length [%d] > max-length [%d]", *flagMinLength, *flagMaxLength)
	}

	// over-ride keyboard handling
	if err := keyboard.Open(); err != nil {
		logErrorAndExit(err.Error())
	}
	defer func() {
		_ = keyboard.Close()
	}()

	// generate a new wordle
	w, err := wordle.New(
		wordle.WithAnswer(*flagAnswer),
		wordle.WithMaxAttempts(*flagMaxAttempts),
		wordle.WithWordFilters(
			wordle.WithLength(*flagMinLength, *flagMaxLength),
			wordle.WithNoRepeatingCharacters(),
		),
	)
	if err != nil {
		logErrorAndExit("failed to initiate new Wordle: %v", err)
	}
	// prompt for user inputs
	prompt(w)
}

func logErrorAndExit(msg string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "ERROR: "+strings.TrimSpace(msg)+"\n", a...)
	os.Exit(-1)
}

func printHelp() {
	fmt.Println(`go-wordle: A GoLang implementation of the Wordle game.

Flags
=====`)
	flag.PrintDefaults()
	os.Exit(0)
}

func prompt(w wordle.Wordle) {
	cliAttempts := strings.Split(*flagAttempts, ",")
	var currAttempt wordle.Attempt
	for {
		render(w, currAttempt)
		if w.Solved() {
			break
		}
		if len(w.Attempts()) == *flagMaxAttempts {
			fmt.Printf("Answer: '%v'\n", strings.ToUpper(w.Answer()))
			break
		}

		// if user provided words to attempt, do that first
		if len(cliAttempts) > 0 {
			_, _ = w.Attempt(cliAttempts[0])
			cliAttempts = cliAttempts[1:]
			continue
		}

		// prompt the user for input
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
			if len(currAttempt.Answer) > 0 {
				currAttempt.Answer = currAttempt.Answer[:len(currAttempt.Answer)-1]
			}
		case keyboard.KeyEnter:
			if len(currAttempt.Answer) == *flagMaxLength {
				_, _ = w.Attempt(currAttempt.Answer)
				currAttempt = wordle.Attempt{}
			}
		default:
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
				if len(currAttempt.Answer) < *flagMaxLength {
					currAttempt.Answer += strings.ToLower(string(char))
				}
			}
		}
	}
}

func render(w wordle.Wordle, currAttempt wordle.Attempt) {
	for linesRendered > 0 {
		fmt.Print(text.CursorUp.Sprint())
		fmt.Print(text.EraseLine.Sprint())
		linesRendered--
	}

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"░ ▒ ▓  W O R D L E  ▓ ▒ ░"})
	tw.AppendRow(table.Row{renderWordle(w, currAttempt)})
	tw.AppendFooter(table.Row{renderKeyboard(w)})
	tw.AppendSeparator()
	tw.AppendFooter(table.Row{colorHints.Sprint("[escape/ctrl+c to quit] [ctrl+r to restart]")})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter, AlignHeader: text.AlignCenter, AlignFooter: text.AlignCenter},
	})
	tw.SetStyle(table.StyleBold)
	tw.Style().Format.Footer = text.FormatDefault
	out := tw.Render()
	linesRendered = strings.Count(out, "\n") + 1
	fmt.Println(out)
}

func renderKeyboard(w wordle.Wordle) string {
	tw := table.NewWriter()
	alphabets := w.Alphabets()

	for _, legend := range []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"} {
		twRow := table.NewWriter()
		row := table.Row{}
		if legend == "zxcvbnm" {
			row = append(row, renderKey("ENTER", colorsSpecial))
		}
		for _, ch := range legend {
			char := string(ch)
			charStatus := alphabets[char]
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
		if legend == "zxcvbnm" {
			row = append(row, renderKey("BKSP", colorsSpecial))
		}
		twRow.AppendRow(row)
		twRow.Style().Options = table.OptionsNoBordersAndSeparators
		tw.AppendRow(table.Row{twRow.Render()})
	}
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
	for idx := 0; idx < *flagMaxLength; idx++ {
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
