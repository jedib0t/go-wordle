package main

import (
	"os"

	"github.com/eiannone/keyboard"
)

func cleanupKeyboard() {
	_ = keyboard.Close()
}

func initKeyboard() {
	// over-ride keyboard handling
	if err := keyboard.Open(); err != nil {
		logErrorAndExit(err.Error())
	}

	// ensure cleanupKeyboard gets called on exit
	exitHandlers = append(exitHandlers, cleanupKeyboard)

	// demo mode needs special Esc/Ctrl+C handling
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
