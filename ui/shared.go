package ui

import (
	"strings"
	"math/rand"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
)

var (
	secondaryColor    = "240"
	primaryColor      = strconv.Itoa(rand.Intn(230))
	maxWidth        = 75
	maxHeight       = 15	
	buttonWidth     = 11
	buttonHeight    = 1
	borderChar      = "-"
    headerBar       = "-"
    headerText      = "tinyC2"
    defaultHelpMsg  = "Welcome to tinyC2!"
)

var (
	headerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = headerBar
        b.Left = headerBar
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
    
    helpBoxStyle = lipgloss.NewStyle().
        Width(maxWidth).
        Height(1).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(secondaryColor))

	bigBoxStyle = lipgloss.NewStyle().
        Width(maxWidth).
        Height(maxHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderForeground(lipgloss.Color(secondaryColor))

	buttonStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.HiddenBorder())
	buttonFocusStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(primaryColor))
    
    inputLabelStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(primaryColor)).
        Bold(true)
    inputTextBoxStyle = lipgloss.NewStyle()
    inputBigBoxStyle = lipgloss.NewStyle().
        Width(maxWidth).
        Height(maxHeight).
        Align(lipgloss.Top, lipgloss.Top).
        BorderForeground(lipgloss.Color(secondaryColor))    

    tableStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(secondaryColor))
)

type state uint
const (
    configState         state = iota
    listenerNewState
    listenersState
    rootState
)

type button struct{
	do 	 func() tea.Msg
	text string
}

type inputSaveMsg   string
func saveConfig() tea.Msg {
    return inputSaveMsg("Config")
}
type inputCancelMsg string
func cancelConfig() tea.Msg {
    return inputCancelMsg("Config")
}
type setHelpMsg     string
type setStateMsg    string
func TODOButton() tea.Msg {
    return setStateMsg("TODO")
}
func toRootState() tea.Msg {
    return setStateMsg("Root")
}
func toConfigState() tea.Msg {
    return setStateMsg("Config")
}
func toListenersState() tea.Msg {
    return setStateMsg("Listeners")
}
func toNewListenerState() tea.Msg {
    return setStateMsg("NewListener")
}

func NextFocus(cur, max int) int {
	f := cur + 1
	if f >= max {
		return 0
	}
	return f
}

func GetButtonViewComponent(buttons []button, focus int) string {
    var bview string
    for i, b := range buttons {
        if i == focus {
            bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
            buttonFocusStyle.Render(b.text))
        } else {
            bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
            buttonStyle.Render(b.text))
        }
    }
    return bview
}

func GetInputViewComponent(labels []string, inputs []textinput.Model) string {
    var iview string
    var temp string
    for x, l := range labels {
        temp = lipgloss.JoinVertical(lipgloss.Top,
            inputLabelStyle.Render(l),
            inputTextBoxStyle.Render(inputs[x].View()))
            iview = lipgloss.JoinVertical(lipgloss.Top, iview, temp)
    }
    return inputBigBoxStyle.Render(iview)
}

func GetHeaderViewComponent() string {
    lline := strings.Repeat(borderChar, 3)
	text := headerStyle.Render(headerText)
	rline := strings.Repeat(borderChar, maxWidth - (len(lline) + len(headerText) + 4))
	return lipgloss.JoinHorizontal(lipgloss.Center, lline, text, rline)
}

func GetFooterViewComponent() string {
	line := strings.Repeat(borderChar, maxWidth)
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}

func GetDemoTable() table.Model {
    numCol := 5
    tWidth := (maxWidth / numCol) - ((2 * numCol) / numCol)
    tHeight := maxHeight - 4

    headers := []string{"Id", "Name", "Ip", "Port", "Status"}
    var cols []table.Column
    for _, h := range headers {
        col := table.Column{Title: h, Width: tWidth}
        cols = append(cols, col)
    }
    
    rows := []table.Row {
        {"1", "Tokyo",      "127.0.0.1", "81", "1"},
        {"2", "Colombo",    "127.0.0.2", "82", "1"},
        {"3", "Toronto",    "127.0.0.3", "83", "1"},
        {"4", "New York",   "127.0.0.4", "84", "1"},
    }

    t := table.New(
        table.WithColumns(cols),
        table.WithRows(rows),
        table.WithFocused(false),
        table.WithHeight(tHeight),
    )
    
    s := table.DefaultStyles()
    s.Header = s.Header.
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(secondaryColor)).
        BorderBottom(true).
        Bold(false)
    s.Selected = s.Selected.
        Foreground(lipgloss.Color(secondaryColor)).
        Background(lipgloss.Color(primaryColor)).
        Bold(false).
        BorderStyle(lipgloss.NormalBorder())
    t.SetStyles(s)

    return t
}
