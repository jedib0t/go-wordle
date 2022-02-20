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
	exitHandlers []func()
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
