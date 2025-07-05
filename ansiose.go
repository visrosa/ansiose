package main

import (
	"fmt"
	"os"
	"strings"
	"visrosa/sgr"

	// "golang.org/x/term"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

var (
	italic              = lipgloss.NewStyle().Italic(true)
	bold                = lipgloss.NewStyle().Bold(true)
	border              = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	left                = lipgloss.NewStyle().Align(lipgloss.Left)
	center              = lipgloss.NewStyle().Align(lipgloss.Center)
	BgMagenta           = lipgloss.NewStyle().Background(lipgloss.Color("54"))
	headerTitleStyle    = lipgloss.NewStyle().Align(lipgloss.Center).Bold(true)
	headerSubtitleStyle = lipgloss.NewStyle().Align(lipgloss.Center).Italic(true)
	// TODO: ASCII and Unicode default mode
	// CSI    = bold.Render("CSI") // TODO: allow setting via --csi-chars=
	// UIn = italic.Render("n") // TODO: allow setting via --user-input-chars=
	CSI = "⍧"
	UIn = "⎀"

	textSizing bool // True if using Kitty's text sizing protocol
)

type model struct {
	SGRheader   string
	simpleTable table.Model
}

func (m model) Init() tea.Cmd {
	if os.Getenv("KITTY_WINDOW_ID") != "" { // TODO: Check support correctly with Ansi Codes, don't hardcode Kitty
		textSizing = true
	} else {
		textSizing = false
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

var SGRheaderTitle = "SGR: Select Graphical Rendition"

func SGRheaderTitleWidth() int {
	if textSizing == true {
		return lipgloss.Width(SGRheaderTitle) * 2
	} else {
		return lipgloss.Width(SGRheaderTitle)
	}
}

func GetSGRheaderTitle() string {
	if textSizing == true {
		// Needs another \n for every size increment for correct alignment
		var t strings.Builder
		// t.WriteString("\x1b]66;s=2:w=0;")
		t.WriteString("\x1b]66;s=3;S\x07\x1b]66;s=2:n=1:d=2:w=1:v=2;GR")
		// t.WriteString(SGRheaderTitle)
		t.WriteString("\x07\x1b[m GR\n")
		return t.String()
	} else {
		return headerTitleStyle.Render(SGRheaderTitle)
	}
}

func (m model) View() string {
	// codes := []sgr.AnsiCode{

	// SGRheaderSubtitle := "Usage: " + CSI + UIn + "m"
	SGRheaderSubtitle := "Usage: " + center.Width(SGRheaderTitleWidth()-7).
		Render(CSI+UIn+"m"+fmt.Sprintf("%v", SGRheaderTitleWidth()))

	SGRheader := border.Render(lipgloss.
		JoinVertical(lipgloss.Left, lipgloss.Place(SGRheaderTitleWidth(), 2, lipgloss.Center, lipgloss.Top,
			GetSGRheaderTitle()),
			SGRheaderSubtitle))
	// view.WriteString(sgr.Underline.On.Apply() + "\x1b[58;5;124mThis is a test of the 256 color mode\x1b[0m")
	if textSizing == true {
		var b strings.Builder
		b.WriteString("╭" + strings.Repeat("─", SGRheaderTitleWidth()+2) + "╮\n")
		b.WriteString(border.BorderBottom(false).BorderTop(false).Width(0).Padding(0, 1).Render(GetSGRheaderTitle()))
		b.WriteString("\n")
		b.WriteString(border.Width(SGRheaderTitleWidth() + 2).BorderBottom(false).BorderTop(false).Align(lipgloss.Center).Render(SGRheaderSubtitle))
		b.WriteString("\n╰" + strings.Repeat("─", SGRheaderTitleWidth()+2) + "╯")
		SGRheader = b.String()
		// SGRheader = lipgloss.JoinVertical(lipgloss.Left, b.String(), SGRheaderSubtitle)

	}

	var view strings.Builder
	view.WriteString(SGRheader)
	view.WriteString("\n")

	formatCodes := []sgr.AnsiCode{
		sgr.Bold.On, sgr.Dim.On, sgr.Italic.On, sgr.Underline.On, sgr.Blink.Slow,
		sgr.Blink.Rapid, sgr.Reverse.On, sgr.Hidden.On, sgr.Strike.On, sgr.DefaultFont, sgr.AlternativeFont1,
	}

	// Build the list of rows
	var rows []table.Row
	for _, code := range formatCodes {
		rows = append(rows, table.NewRow(table.RowData{
			"Name":     code.Apply() + code.Name,
			"AnsiCode": code.Code,
		}))
	}

	var cols []table.Column
	cols = append(cols, table.NewColumn("AnsiCode", center.Width(3).Bold(true).Render(UIn), 3))
	cols = append(cols, table.NewColumn("Name", "Name", lipgloss.Width(SGRheader)-6).WithStyle(center))

	simpleTable := table.New(cols).WithRows(rows)

	view.WriteString(simpleTable.View())
	return view.String()

}
func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
	// codes := []sgr.AnsiCode{
	// 	sgr.Bold.On, sgr.Dim.On, sgr.Italic.On, sgr.Underline.On, sgr.Blink.Slow, sgr.Blink.Rapid, sgr.Reverse.On, sgr.Hidden.On, sgr.Strike.On,
	// 	sgr.FgBlack, sgr.FgRed, sgr.FgGreen, sgr.FgYellow, sgr.FgBlue, sgr.FgMagenta, sgr.FgCyan, sgr.FgWhite, sgr.FgDefault,
	// 	sgr.BgBlack, sgr.BgRed, sgr.BgGreen, sgr.BgYellow, sgr.BgBlue, sgr.BgMagenta, sgr.BgCyan, sgr.BgWhite, sgr.BgDefault,
	// 	sgr.FgBrightBlack, sgr.FgBrightRed, sgr.FgBrightGreen, sgr.FgBrightYellow, sgr.FgBrightBlue, sgr.FgBrightMagenta, sgr.FgBrightCyan, sgr.FgBrightWhite,
	// 	sgr.BgBrightBlack, sgr.BgBrightRed, sgr.BgBrightGreen, sgr.BgBrightYellow, sgr.BgBrightBlue, sgr.BgBrightMagenta, sgr.BgBrightCyan, sgr.BgBrightWhite,
	// }

	// for _, code := range codes {
	// 	fmt.Println(code.Render())
	// }

	// width, height, err := term.GetSize(int(os.Stdout.Fd()))
	// if err != nil {
	// 	fmt.Println("Could not get terminal size:", err)
	// 	return
	// }

	// fmt.Printf("Terminal size: %d columns × %d rows\n", width, height)
	// fmt.Print(sgr.Bold.On.Apply() + "CSI " + sgr.SetForeground.Code + " " + sgr.SetForeground.Name + " (256 colors):\t" + sgr.Reset.Apply())
	// fmt.Println("⎆ " + sgr.Dim.On.Apply() + "38;5;" + sgr.Reset.Apply() + sgr.Italic.On.Apply() + "n" + sgr.Reset.Apply() + sgr.Dim.On.Apply() + "m" + sgr.Reset.Apply())
	// fmt.Println(sgr.Reset.Apply())
	// fmt.Println("Standard colors         High-intensity colors")

	// for i := 0; i <= 7; i++ {
	// 	fmt.Print(sgr.CSI(sgr.SetForeground.Code, fmt.Sprint(i)) + fmt.Sprintf("%02d ", i) + sgr.Reset.Apply())
	// }

	// for i := 8; i <= 15; i++ {
	// 	fmt.Print(sgr.CSI(sgr.SetForeground.Code, fmt.Sprint(i)) + fmt.Sprintf("%02d ", i) + sgr.Reset.Apply())
	// }

	// fmt.Println()
	// for i := 0; i <= 7; i++ {
	// 	fmt.Print(sgr.CSI(sgr.SetBackground.Code, fmt.Sprint(i)) + fmt.Sprintf("%02d ", i) + sgr.Reset.Apply())
	// }

	// for i := 8; i <= 15; i++ {
	// 	fmt.Print(sgr.CSI(sgr.SetBackground.Code, fmt.Sprint(i)) + fmt.Sprintf("%02d ", i) + sgr.Reset.Apply())
	// }
	// fmt.Println()
	// iterationsPerLine := width / 26
	// fmt.Println("Block size:", iterationsPerLine)

	// count := 0
	// for i := 16; i <= 231; i++ {
	// 	fmt.Print(sgr.CSI(sgr.SetBackground.Code, fmt.Sprint(i)) + fmt.Sprintf("%03d ", i) + sgr.Reset.Apply())
	// 	if (i-15)%6 == 0 {
	// 		fmt.Print(" ")
	// 	}
	// 	count++
	// 	if count == iterationsPerLine*6 {
	// 		fmt.Println()
	// 		count = 0
	// 	}
	// }
	// fmt.Println("\nGrayscale colors:")
	// for i := 232; i <= 255; i++ {
	// 	fmt.Print(sgr.CSI(sgr.SetBackground.Code, fmt.Sprint(i)) + fmt.Sprintf("%03d ", i) + sgr.Reset.Apply())
	// }
	// fmt.Println()
	// count = 0
	// for i := 16; i <= 231; i++ {
	// 	fmt.Print(sgr.CSI(sgr.SetForeground.Code, fmt.Sprint(i)) + fmt.Sprintf("%03d ", i) + sgr.Reset.Apply())
	// 	if (i-15)%6 == 0 {
	// 		fmt.Print(" ")
	// 	}
	// 	count++
	// 	if count == iterationsPerLine*6 {
	// 		fmt.Println()
	// 		count = 0
	// 	}
	// }
	// fmt.Println("Grayscale colors:")
	// for i := 232; i <= 255; i++ {
	// 	fmt.Print(sgr.CSI(sgr.SetForeground.Code, fmt.Sprint(i)) + fmt.Sprintf("%03d ", i) + sgr.Reset.Apply())
	// }
	// fmt.Println()

}
