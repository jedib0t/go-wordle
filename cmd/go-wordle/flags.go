package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	// defaults
	defaultSolveSpeed  = 4
	defaultMaxAttempts = 6
	defaultNumWordles  = 1
	defaultWordLength  = 5

	// min-max
	solveSpeedMin = 1
	solveSpeedMax = 10

	// flags
	flagAnswer      = flag.String("answer", "", "Pre-set answer if you don't want a random word")
	flagAttempts    = flag.String("attempts", "", "Words to attempt before prompting user")
	flagHelp        = flag.Bool("help", false, "Show this help-text?")
	flagHelper      = flag.Bool("helper", false, "Help solve Wordle puzzle from elsewhere?")
	flagHints       = flag.Bool("hints", false, "Show hints and help solve?")
	flagMaxAttempts = flag.Int("max-attempts", defaultMaxAttempts, "Maximum attempts allowed")
	flagNumWordles  = flag.Int("num-wordles", defaultNumWordles, "Number of Wordle Puzzles")
	flagSolve       = flag.Bool("solve", false, "Solve the puzzle?")
	flagSolveSpeed  = flag.Int("solve-speed", defaultSolveSpeed, "Speed of 'solve' (1-10, 10 being fastest)")
	flagWordLength  = flag.Int("word-length", defaultWordLength, "Number of characters in the Word")
)

func initFlags() {
	flag.Parse()
	if *flagHelp {
		printHelp()
	}
	if *flagHelper && *flagSolve {
		logErrorAndExit("flags -helper and -solve cannot be used together")
	}
	if *flagHelper || *flagSolve {
		*flagHints = true
	}
	if *flagSolveSpeed < solveSpeedMin || *flagSolveSpeed > solveSpeedMax {
		*flagSolveSpeed = defaultSolveSpeed
	}
	if *flagAnswer != "" {
		answers := strings.Split(*flagAnswer, ",")
		for idx := 1; idx < len(answers); idx++ {
			if len(answers[idx]) != len(answers[idx-1]) {
				logErrorAndExit("answers provided %v have differing lengths", answers)
			}
		}
		*flagWordLength = len(answers[0])
	}
	if *flagAttempts != "" {
		for idx, attempt := range strings.Split(*flagAttempts, ",") {
			if len(attempt) != *flagWordLength {
				logErrorAndExit("attempts[%d]='%s' does not meet word-length %d", idx, attempt, *flagWordLength)
			}
		}
	}
	if *flagWordLength < 3 || *flagWordLength > 10 {
		logErrorAndExit("word-length [%d] has to be between 3 and 10", *flagWordLength)
	}
}

func printHelp() {
	fmt.Println(`go-wordle: A GoLang implementation of the Wordle game.

Version: ` + version + `

Flags
=====`)
	flag.PrintDefaults()
	os.Exit(0)
}
