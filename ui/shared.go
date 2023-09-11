package ui

import (
	"strings"
	"math/rand"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
)

var (
	headerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = headerBar
        b.Left = headerBar
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

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
	focusButtonStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(primaryColor))
)

type button struct{
	do 	 func() tea.Msg
	text string
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
            focusButtonStyle.Render(b.text))
        } else {
            bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
            buttonStyle.Render(b.text))
        }
    }
    return bview
}

func GetHeaderViewComponent(s string) string {
    lline := strings.Repeat(borderChar, 3)
	text := headerStyle.Render(s)
	rline := strings.Repeat(borderChar, maxWidth - (len(lline) + len(headerText) + 4))
	return lipgloss.JoinHorizontal(lipgloss.Center, lline, text, rline)
}

func GetFooterViewComponent() string {
	line := strings.Repeat(borderChar, maxWidth)
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}