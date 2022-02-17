# go-wordle

[![Build Status](https://github.com/jedib0t/go-wordle/workflows/CI/badge.svg?branch=main)](https://github.com/jedib0t/go-wordle/actions?query=workflow%3ACI+event%3Apush+branch%3Amain)

A golang implementation of the popular New York Times
game [Wordle](https://www.nytimes.com/games/wordle/index.html).

<img src="go-wordle.gif"/>

## Install

Pre-built binaries for your Operating System can be found at
the [latest release](https://github.com/jedib0t/go-wordle/releases/latest)
page.

If you want to build you own using GoLang:

* `go get -u github.com/jedib0t/go-wordle/cmd/go-wordle`
* `go-wordle`

If you want to run from source, after `git clone`:

* `go run ./cmd/go-wordle` or `make run`

## Features

* Hinting mode with `-hints`
    * Shows recommended/possible answers above the keyboard
* Solve mode with `-solve`
    * Automated solving with or without a pre-set answer
    * Uses the built-in hinting system to choose answers
* Helper mode to solve external Wordle puzzles with `-helper`
    * If you are using this tool as a cheat and solve Wordle puzzles elsewhere ;)
    * Show hints and helps record the results to get you to the answer
* _Mega_ Wordles mode with `-num-wordles`
    * Run more than 1 Wordle puzzle at the same time
    * Needs a wide-enough terminal to fit the contents