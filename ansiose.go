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
	FgMagenta           = lipgloss.NewStyle().Foreground(lipgloss.Color("54"))
	colHeader           = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
	headerTitleStyle    = lipgloss.NewStyle().Align(lipgloss.Center).Bold(true).Foreground(lipgloss.Color("54"))
	headerSubtitleStyle = lipgloss.NewStyle().Align(lipgloss.Center).Italic(true)
	// TODO: ASCII and Unicode default mode
	// CSI    = bold.Render("CSI") // TODO: allow setting via --csi-chars=
	// UIn = italic.Render("n") // TODO: allow setting via --user-input-chars=
	CSI        = "⍧"
	UIn        = "⎀"
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

// func GetSGRheaderTitle() string {
// 	if textSizing == true {
// 		// Needs another \n for every size increment for correct alignment
// 		var t strings.Builder
// 		// t.WriteString("\x1b]66;s=2:w=0;")
// 		t.WriteString(sgr.TextSize.Apply("s=2:w=0"))
// 		// t.WriteString("\x1b]66;s=3;S\x07\x1b]66;s=2:n=1:d=2:w=1:v=2;GR")
// 		t.WriteString(SGRheaderTitle)
// 		t.WriteString(sgr.TextSize.Off() + "\n")
// 		// t.WriteString("\x07\x1b[m\n")
// 		// t.WriteString(sgr.TextSize.Apply("2", "", "hey") + sgr.TextSize.Off())
// 		// t.WriteString(sgr.TextSize.Apply("2", "n=1:d=2:v=0", ""))
// 		return t.String()
// 	} else {
// 		return "SGR Header"
// 	}
// }

func (m model) View() string {

	// SGRheaderSubtitle := "Usage: " + CSI + UIn + "m"
	SGRheaderSubtitle := italic.Bold(true).Render("Usage: ") + center.Width(SGRheaderTitleWidth()-7).Render(CSI+UIn+"m")

	var SGRheader string

	if textSizing == true {
		var b strings.Builder
		b.WriteString("╭" + strings.Repeat("─", SGRheaderTitleWidth()+2) + "╮\n")
		b.WriteString(border.BorderBottom(false).BorderTop(false).Width(0).Padding(0, 1).Bold(true).
			Render(FgMagenta.Render(sgr.TextSize.Apply("s=2:w=0") + SGRheaderTitle + sgr.TextSize.Off() + "\n")))
		//b.WriteString(border.BorderBottom(false).BorderTop(false).Render(strings.Repeat(" ", SGRheaderTitleWidth()+2)))
		b.WriteString("\n")
		b.WriteString(border.Width(SGRheaderTitleWidth() + 2).BorderBottom(false).BorderTop(false).Align(lipgloss.Center).Render(SGRheaderSubtitle))
		b.WriteString("\n╰" + strings.Repeat("─", SGRheaderTitleWidth()+2) + "╯")
		SGRheader = b.String()
		// SGRheader = lipgloss.JoinVertical(lipgloss.Left, b.String(), SGRheaderSubtitle)
	} else {
		SGRheader = border.Render(lipgloss.JoinVertical(lipgloss.Center, SGRheaderTitle, SGRheaderSubtitle))
	}

	var view strings.Builder
	view.WriteString(SGRheader)
	view.WriteString("\n")

	// Prepare the two slices
	formatCodes := []sgr.AnsiCode{
		sgr.Bold.On, sgr.Dim.On, sgr.Italic.On, sgr.Underline.On, sgr.Blink.Slow,
		sgr.Blink.Rapid, sgr.Reverse.On, sgr.Hidden.On, sgr.Strike.On, sgr.DefaultFont,
	}

	fontCodes := []sgr.AnsiCode{
		sgr.AlternativeFont1, sgr.AlternativeFont2, sgr.AlternativeFont3, sgr.AlternativeFont4, sgr.AlternativeFont5,
		sgr.AlternativeFont6, sgr.AlternativeFont7, sgr.AlternativeFont8, sgr.AlternativeFont9, sgr.GothicFont,
	}

	// Build the combined rows
	var rows []table.Row
	maxLen := len(formatCodes)
	if len(fontCodes) > maxLen {
		maxLen = len(fontCodes)
	}
	for i := 0; i < maxLen; i++ {
		rowData := table.RowData{}
		if i < len(formatCodes) {
			code := formatCodes[i]
			rowData["AnsiCode1"] = fmt.Sprintf("%2d", i+1) // or code.Code if you prefer
			rowData["FormatCodes"] = code.Apply() + code.Name
		}
		if i < len(fontCodes) {
			code := fontCodes[i]
			rowData["AnsiCode2"] = fmt.Sprintf("%2d", i+1+len(formatCodes)) // or code.Code
			rowData["FontCodes"] = code.Apply() + code.Name
		}
		rows = append(rows, table.NewRow(rowData))
	}

	var cols []table.Column
	cols = append(cols, table.NewColumn("AnsiCode1", center.Width(3).Bold(true).Render(UIn), 3))
	cols = append(cols, table.NewColumn("FormatCodes", colHeader.Render("Formatting"), 14).WithStyle(center))
	cols = append(cols, table.NewColumn("AnsiCode2", center.Width(3).Bold(true).Render(UIn), 3))
	cols = append(cols, table.NewColumn("FontCodes", "Fonts", 20).WithStyle(center))

	simpleTable := table.New(cols).WithRows(rows)

	view.WriteString(simpleTable.View())
	return view.String()

}
func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
	fmt.Println("\n\n\n\n\n\n")

	// fmt.Print(border.MaxWidth(8).Width(4).Height(2).Render(lipgloss.NewStyle().Width(4).Render(sgr.TextSize.Apply("s=2:w=0") + " S" + sgr.TextSize.Off() + "")))
	fmt.Println(sgr.Fg.Color(34).Apply() + " " + sgr.Reset.Apply() + "fmt.Print(sgr.Bold.On.Apply() + \"Bold text\" + sgr.Bold.Off.Apply())")
	// fmt.Println("    " + sgr.Bold.On.Apply() + "Bold text" + sgr.Bold.Off.Apply())
	frame := lipgloss.NewStyle().Background(lipgloss.Color("00")).MarginLeft(4).MarginBottom(1).Padding(1)
	fmt.Println(frame.Bold(true).Render("Bold text"))
	fmt.Println(sgr.Fg.Color(34).Apply() + " " + sgr.Reset.Apply() + "fmt.Print(sgr.Fg.Color(34).Apply() + \"Blue text\" + sgr.Reset.Apply())")
	fmt.Println(frame.Foreground(lipgloss.Color("27")).Render("Blue text"))

	fmt.Println(sgr.Fg.Color(34).Apply() + " " + sgr.Reset.Apply() + "fmt.Print(sgr.TextSize.Render())\n")
	// fmt.Print("   ")
	fmt.Print(sgr.TextSize.Apply("s=4") + " " + sgr.TextSize.Off() + sgr.Bg.Color(00).Apply() + sgr.TextSize.Render() + sgr.Reset.Apply())
	fmt.Println("\n\n")
	// fmt.Println("\n\n\n")
	// fmt.Print(sgr.TextSize.Apply("n=1:d=2:v=0:w=1") + "01\n" + sgr.TextSize.Off())
	// fmt.Print(sgr.TextSize.Apply("n=1:d=2:v=1:w=1") + "02\n" + sgr.TextSize.Off())
	// fmt.Print(sgr.TextSize.Apply("n=1:d=2:v=2:w=1") + "02\n" + sgr.TextSize.Off())
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
