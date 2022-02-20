package main

import (
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	titleText1 = `       ▞ ▛▀▀▀▀▀▀▀▀▀▀▀▀▀▜ ▚       
░ ▒ ▓ █ ▌  W O R D L E  ▐ █ ▓ ▒ ░
       ▚ ▙▄▄▄▄▄▄▄▄▄▄▄▄▄▟ ▞       `
	titleText2 = `░░     ░░  ░░░░░░  ░░░░░░  ░░░░░░  ░░      ░░░░░░░ 
▒▒     ▒▒ ▒▒    ▒▒ ▒▒   ▒▒ ▒▒   ▒▒ ▒▒      ▒▒      
▒▒  ▒  ▒▒ ▒▒    ▒▒ ▒▒▒▒▒▒  ▒▒   ▒▒ ▒▒      ▒▒▒▒▒   
▓▓ ▓▓▓ ▓▓ ▓▓    ▓▓ ▓▓   ▓▓ ▓▓   ▓▓ ▓▓      ▓▓      
 ███ ███   ██████  ██   ██ ██████  ███████ ███████ `
	titleText3 = `▄▄▌ ▐ ▄▌      ▄▄▄  ·▄▄▄▄  ▄▄▌  ▄▄▄ .
██· █▌▐█ ▄█▀▄ ▀▄ █·██· ██ ██•  ▀▄.▀·
██▪▐█▐▐▌▐█▌.▐▌▐▀▀▄ ▐█▪ ▐█▌██ ▪ ▐▀▀▪▄
▐█▌██▐█▌▐█▌.▐▌▐█•█▌██. ██ ▐█▌ ▄▐█▄▄▌
 ▀▀▀▀ ▀▪ ▀█▄▀▪.▀  ▀▀▀▀▀▀• .▀▀▀  ▀▀▀ `
	titleText = titleText2
)

var (
	colorTitle = []text.Color{
		text.FgHiWhite,
		text.FgHiYellow,
		text.FgHiCyan,
		text.FgHiBlue,
		text.FgHiGreen,
		text.FgHiMagenta,
		text.FgHiRed,
	}
	colorTitleAnimated = false
	colorTitleIdx      = 0
	colorTitleOnce     = sync.Once{}
)

func switchTitleColorAsync() {
	colorTitleOnce.Do(func() {
		// async switch colors every few milliseconds for the title
		go func() {
			for {
				time.Sleep(time.Second / 2)
				if colorTitleIdx < len(colorTitle)-1 {
					colorTitleIdx++
				} else {
					colorTitleIdx = 0
				}
			}
		}()
	})
}

func getTitleColors() text.Colors {
	if colorTitleAnimated {
		switchTitleColorAsync()
	}

	return text.Colors{colorTitle[colorTitleIdx], text.Bold}
}

func renderTitle() string {
	colors := getTitleColors()

	tw := table.NewWriter()
	for _, line := range strings.Split(titleText, "\n") {
		tw.AppendRow(table.Row{colors.Sprint(line)})
	}
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}
