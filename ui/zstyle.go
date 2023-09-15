package ui
import (
	"math/rand"
	"strconv"

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
    defaultinfoMsg  = "Welcome to tinyC2!"
)

var (
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
        
    inputBigBoxStyle = lipgloss.NewStyle().
        Width(maxWidth).
        Height(maxHeight).
        Align(lipgloss.Top, lipgloss.Top).
        BorderForeground(lipgloss.Color(secondaryColor)) 
    inputLabelStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(primaryColor)).
        Bold(true)
    inputTextBoxStyle = lipgloss.NewStyle()
       
    tableStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(secondaryColor))
)