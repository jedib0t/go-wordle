package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// flags
	flagAnswer      = flag.String("answer", "", "Pre-set answer if you don't want a random word")
	flagAttempts    = flag.String("attempts", "", "Words to attempt before prompting user")
	flagDemo        = flag.Bool("demo", false, "Show an automated demo?")
	flagHelp        = flag.Bool("help", false, "Show this help-text?")
	flagHelper      = flag.Bool("helper", false, "Help solve Wordle puzzle from elsewhere?")
	flagHints       = flag.Bool("hints", false, "Show hints and help solve?")
	flagMaxAttempts = flag.Int("max-attempts", 6, "Maximum attempts allowed")
	flagWordLength  = flag.Int("word-length", 5, "Number of characters in the Word")
)

func initFlags() {
	flag.Parse()
	if *flagHelp {
		printHelp()
	}
	if *flagWordLength < 3 || *flagWordLength > 10 {
		logErrorAndExit("word-length [%d] has to be between 3 and 10", *flagWordLength)
	}
	if *flagDemo || *flagHelper {
		*flagHints = true
	}
	if *flagAnswer != "" {
		*flagWordLength = len(*flagAnswer)
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
