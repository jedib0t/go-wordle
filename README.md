# go-wordle

[![Build Status](https://github.com/jedib0t/go-wordle/workflows/CI/badge.svg?branch=main)](https://github.com/jedib0t/go-wordle/actions?query=workflow%3ACI+event%3Apush+branch%3Amain)

Play the popular New York Times
game [Wordle](https://www.nytimes.com/games/wordle/index.html)
in the terminal.

<img src="images/solver.gif" alt="go-wordle -solve"/>

More samples can be found in the [/images](images) folder.

## Install

Pre-built binaries for your Operating System can be found at
the [latest release](https://github.com/jedib0t/go-wordle/releases/latest)
page.

If you want to build your own using GoLang:

* `go get -u github.com/jedib0t/go-wordle/cmd/go-wordle`
* `go-wordle`

If you want to run from source:

* `git clone git@github.com/jedib0t/go-wordle`
* `cd go-wordle`
* `make run`

## Features

* Hinting mode with `-hints`
    * Shows recommended/possible answers above the keyboard
* Solve mode with `-solve`
    * Automated solving with or without a pre-set answer
    * Uses the built-in hinting system to choose answers
* _Mega_ Wordles mode with `-num-wordles`
    * Run more than 1 Wordle puzzle at the same time
    * Needs a wide-enough terminal to fit the contents
* Helper mode to solve external Wordle puzzles with `-helper`
    * Use this tool as a cheat and solve Wordle puzzles elsewhere ;)
    * Show hints and helps record the results to get you closer to the answer
