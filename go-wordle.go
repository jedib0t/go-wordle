package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-wordle/wordle"
)

var (
	flagAnswer        = flag.String("answer", "", "Pre-set answer if you don't want a random word")
	flagAttempts      = flag.String("attempts", "", "Words to attempt before prompting user")
	flagDemo          = flag.Bool("demo", false, "Show an automated demo?")
	flagHelp          = flag.Bool("help", false, "Show this help-text?")
	flagHints         = flag.Bool("hints", false, "Show hints and help solve?")
	flagMaxAttempts   = flag.Int("max-attempts", 6, "Maximum attempts allowed")
	flagMaxLength     = flag.Int("max-length", 5, "Maximum word length")
	flagMinLength     = flag.Int("min-length", 5, "Minimum word length")
	flagSolveExternal = flag.Bool("solve-external", false, "Solve Wordle puzzle from elsewhere?")

	colorHints              = text.Colors{text.FgHiBlack, text.Italic}
	colorsSpecial           = [3]text.Color{text.FgBlack, text.BgBlack, text.FgHiYellow}
	colorsUnknown           = [3]text.Color{text.FgHiBlack, text.BgHiBlack, text.FgHiWhite}
	colorsNotPresent        = [3]text.Color{text.FgBlack, text.BgBlack, text.FgHiBlack}
	colorsInWrongLocation   = [3]text.Color{text.FgHiYellow, text.BgHiYellow, text.FgBlack}
	colorsInCorrectLocation = [3]text.Color{text.FgHiGreen, text.BgHiGreen, text.FgBlack}
	demoWord                = ""
	demoWordSet             = false
	linesRendered           = 0
	inputModeCharStatus     = false
)

func main() {
	initFlagsAndKeyboard()
	defer exitHandler()

	filters := []wordle.Filter{
		wordle.WithLength(*flagMinLength, *flagMaxLength),
	}
	opts := []wordle.Option{
		wordle.WithAnswer(*flagAnswer),
		wordle.WithMaxAttempts(*flagMaxAttempts),
		wordle.WithWordFilters(filters...),
	}
	if *flagSolveExternal {
		opts = append(opts, wordle.WithUnknownAnswer(*flagMaxLength))
	}

	// generate a new wordle
	w, err := wordle.New(opts...)
	if err != nil {
		logErrorAndExit("failed to initiate new Wordle: %v", err)
	}
	if *flagDemo && *flagAnswer != "" && !w.DictionaryHas(*flagAnswer) {
		logErrorAndExit("demo cannot proceed as '%s' is not in dictionary", *flagAnswer)
	}
	// prompt for user inputs
	prompt(w)
}

func exitHandler() {
	_ = keyboard.Close()
}

func getUserInput(w wordle.Wordle, currAttempt wordle.Attempt, hints []string) (wordle.Attempt, []string) {
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
		if inputModeCharStatus {
			if len(currAttempt.Result) > 0 {
				currAttempt.Result = currAttempt.Result[:len(currAttempt.Result)-1]
			}
		} else {
			if len(currAttempt.Answer) > 0 {
				currAttempt.Answer = currAttempt.Answer[:len(currAttempt.Answer)-1]
			}
		}
	case keyboard.KeyEnter:
		if len(currAttempt.Answer) == *flagMaxLength {
			if *flagSolveExternal && len(currAttempt.Result) < len(currAttempt.Answer) {
				inputModeCharStatus = true
			} else {
				inputModeCharStatus = false
				_, _ = w.Attempt(currAttempt.Answer, currAttempt.Result...)
				hints = w.Hints()
				currAttempt = wordle.Attempt{}
			}
		}
	default:
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			if len(currAttempt.Answer) < *flagMaxLength {
				currAttempt.Answer += strings.ToLower(string(char))
			}
		} else if inputModeCharStatus {
			if char == '0' {
				currAttempt.Result = append(currAttempt.Result, wordle.NotPresent)
			} else if char == '2' {
				currAttempt.Result = append(currAttempt.Result, wordle.PresentInWrongLocation)
			} else if char == '3' {
				currAttempt.Result = append(currAttempt.Result, wordle.PresentInCorrectLocation)
			}
		}
	}
	return currAttempt, hints
}

func initFlagsAndKeyboard() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	if *flagHelp {
		printHelp()
	}
	if *flagMinLength > *flagMaxLength {
		logErrorAndExit("min-length [%d] > max-length [%d]", *flagMinLength, *flagMaxLength)
	}
	if *flagDemo || *flagSolveExternal {
		*flagHints = true
	}

	// over-ride keyboard handling
	if err := keyboard.Open(); err != nil {
		logErrorAndExit(err.Error())
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		exitHandler()
	}()

	// demo mode needs special Esp/Ctrl+C handling
	if *flagDemo {
		go func() {
			_, key, _ := keyboard.GetSingleKey()
			switch key {
			case keyboard.KeyEsc, keyboard.KeyCtrlC:
				os.Exit(0)
			}
		}()
	}
}

func logErrorAndExit(msg string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "ERROR: "+strings.TrimSpace(msg)+"\n", a...)
	exitHandler()
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
	currAttempt := wordle.Attempt{}
	demoWord = ""
	hints := w.Hints()
	for {
		render(w, hints, currAttempt)
		if w.Solved() {
			break
		}
		if len(w.Attempts()) == *flagMaxAttempts {
			if !*flagSolveExternal {
				fmt.Printf("Answer: '%v'\n", strings.ToUpper(w.Answer()))
			}
			break
		}

		// if user provided words to attempt, do that first
		if len(cliAttempts) > 0 {
			_, _ = w.Attempt(cliAttempts[0])
			cliAttempts = cliAttempts[1:]
			hints = w.Hints()
			continue
		}

		// prompt the user for input
		if *flagDemo {
			if len(hints) == 0 {
				logErrorAndExit("Uh oh... failed to solve the Wordle! Big Sad!")
			}
			currAttempt, hints = demoSolveWithHints(w, currAttempt, hints)
		} else {
			currAttempt, hints = getUserInput(w, currAttempt, hints)
		}
	}
}

func demoSolveWithHints(w wordle.Wordle, currAttempt wordle.Attempt, hints []string) (wordle.Attempt, []string) {
	if demoWordSet {
		time.Sleep(time.Second)
		// if the word is empty and moved over to the answer, attempt it
		if demoWord == "" {
			_, _ = w.Attempt(currAttempt.Answer)
			demoWordSet = false
			return wordle.Attempt{}, w.Hints()
		}
	}
	if demoWord == "" {
		demoWord = hints[0]
		demoWordSet = true
	}
	if len(demoWord) > 0 {
		// move one letter over to the answer
		currAttempt.Answer += string(demoWord[0])
		demoWord = demoWord[1:]
	}
	return currAttempt, hints
}

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

	if *flagSolveExternal && inputModeCharStatus {
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
