package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	// exitHandlers contains all functions that need to be called during exit
	exitHandlers []func()
	// timeStart is used to render the game timer
	timeStart = time.Now()
	// version
	version = "dev"
)

func init() {
	// seed the RNG or the word auto-selected will remain the same all the time
	rand.Seed(time.Now().UnixNano())

	// cleanup on termination
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		cleanup()
	}()

	// init other things
	initFlags()
	initKeyboard()
	initHeaderAndFooter()
}

func cleanup() {
	renderMutex.Lock()
	defer renderMutex.Unlock() // unnecessary
	renderEnabled = false

	for _, exitHandler := range exitHandlers {
		exitHandler()
	}
}

func logErrorAndExit(msg string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "ERROR: "+strings.TrimSpace(msg)+"\n", a...)
	cleanup()
	os.Exit(-1)
}

func main() {
	defer cleanup()

	// instantiate
	wordles := generateWordles(*flagNumWordles)

	// play
	play(wordles)
}
