package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func cleanupKeyboard() {
	cursorShow()
	_ = keyboard.Close()
}

func cursorHide() {
	fmt.Printf("\x1b[?25l")
}

func cursorShow() {
	fmt.Printf("\x1b[?25h")
}

func initKeyboard() {
	// over-ride keyboard handling
	if err := keyboard.Open(); err != nil {
		logErrorAndExit(err.Error())
	}
	cursorHide()

	// ensure cleanupKeyboard gets called on exit
	exitHandlers = append(exitHandlers, cleanupKeyboard)

	// solve mode needs special Esc/Ctrl+C handling
	if *flagSolve {
		go func() {
			_, key, _ := keyboard.GetSingleKey()
			switch key {
			case keyboard.KeyEsc, keyboard.KeyCtrlC:
				handleShortcutQuit()
			}
		}()
	}
}
