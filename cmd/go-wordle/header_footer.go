package main

import (
	"fmt"
	"log"
	"net"
	"os/user"
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	// Official Unicode Consortium code chart (PDF)
	//          0	1	2	3	4	5	6	7	8	9	A	B	C	D	E	F
	// U+258x	▀	▁	▂	▃	▄	▅	▆	▇	█	▉	▊	▋	▌	▍	▎	▏
	// U+259x	▐	░	▒	▓	▔	▕	▖	▗	▘	▙	▚	▛	▜	▝	▞	▟

	titleText1 = `
       ▞ ▛▀▀▀▀▀▀▀▀▀▀▀▀▀▜ ▚       
░ ▒ ▓ █ ▌  W O R D L E  ▐ █ ▓ ▒ ░
       ▚ ▙▄▄▄▄▄▄▄▄▄▄▄▄▄▟ ▞       
`
	titleTextBig = `
░░     ░░  ░░░░░░  ░░░░░░  ░░░░░░  ░░      ░░░░░░░ 
▒▒     ▒▒ ▒▒    ▒▒ ▒▒   ▒▒ ▒▒   ▒▒ ▒▒      ▒▒      
▒▒  ▒  ▒▒ ▒▒    ▒▒ ▒▒▒▒▒▒  ▒▒   ▒▒ ▒▒      ▒▒▒▒▒   
▓▓ ▓▓▓ ▓▓ ▓▓    ▓▓ ▓▓   ▓▓ ▓▓   ▓▓ ▓▓      ▓▓      
 ███ ███   ██████  ██   ██ ██████  ███████ ███████ 
`
	titleTextScary = `
▄▄▌ ▐ ▄▌      ▄▄▄  ·▄▄▄▄  ▄▄▌  ▄▄▄ .
██· █▌▐█ ▄█▀▄ ▀▄ █·██· ██ ██•  ▀▄.▀·
██▪▐█▐▐▌▐█▌.▐▌▐▀▀▄ ▐█▪ ▐█▌██ ▪ ▐▀▀▪▄
▐█▌██▐█▌▐█▌.▐▌▐█•█▌██. ██ ▐█▌ ▄▐█▄▄▌
 ▀▀▀▀ ▀▪ ▀█▄▀▪.▀  ▀▀▀▀▀▀• .▀▀▀  ▀▀▀ 
`
	titleText = titleTextScary
)

var (
	colorAddress = text.Colors{text.Faint}
	colorTimer   = text.Colors{text.FgHiCyan}
	colorTitle   = []text.Color{
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
	colorUser          = text.Colors{text.FgHiBlue, text.Bold}

	localAddress  *net.UDPAddr
	localUsername = "unknown"
	titleOnce     sync.Once
)

func initHeaderAndFooter() {
	titleOnce.Do(func() {
		// localUsername
		if u, err := user.Current(); err == nil {
			localUsername = u.Username
		}

		// localAddress
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = conn.Close()
		}()
		localAddress = conn.LocalAddr().(*net.UDPAddr)

		// title Colors
		if colorTitleAnimated {
			colorTitleOnce.Do(func() {
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
	})
}

func renderFooter() string {
	timeGame := time.Now().Sub(timeStart)
	timeGameStr := fmt.Sprintf("%02d:%02d:%02d",
		int(timeGame.Hours()), int(timeGame.Minutes()), int(timeGame.Seconds()))

	tw := table.NewWriter()
	tw.AppendRow(table.Row{
		colorUser.Sprint(localUsername),
		colorTimer.Sprint(timeGameStr),
	})
	tw.SetStyle(table.StyleLight)
	tw.Style().Options.DrawBorder = false
	return tw.Render()
}

func renderTitle() string {
	colors := text.Colors{colorTitle[colorTitleIdx], text.Bold}

	tw := table.NewWriter()
	for _, line := range strings.Split(titleText, "\n") {
		if line != "" {
			tw.AppendRow(table.Row{colors.Sprint(line)})
		}
	}
	tw.Style().Options = table.OptionsNoBordersAndSeparators
	return tw.Render()
}

func renderUserDetails() string {
	timeGame := time.Now().Sub(timeStart)
	timeGameStr := fmt.Sprintf("%02d:%02d:%02d",
		int(timeGame.Hours()), int(timeGame.Minutes()), int(timeGame.Seconds()))

	tw := table.NewWriter()
	tw.AppendRow(table.Row{colorUser.Sprint(localUsername)})
	tw.AppendSeparator()
	tw.AppendRow(table.Row{colorTimer.Sprint(timeGameStr)})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter},
	})
	tw.SetStyle(table.StyleLight)
	return tw.Render()
}
