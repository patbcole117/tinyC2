package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	primaryColor    = "240"
	focusColor      = "69"
	maxWidth        = 100
	maxHeight       = 15
	buttonWidth     = 15
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

	infoBoxStyle = lipgloss.NewStyle().
        Width(maxWidth).
        Height(maxHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderForeground(lipgloss.Color(primaryColor))

	focusTextBoxStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(focusColor))

    baseTableStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(primaryColor))

	textBoxStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.HiddenBorder())
)

type sessionState uint
const (
	mainAgentsState sessionState = iota
	agentState
    mainCLIState
	listenersState
    listenerNewState
    listenerInfoState
	listenerEditState
	mainState
)

type viewFocus uint
const (
	mainAgentsFocus viewFocus = iota
	mainListenersFocus
	mainCLIFocus
    listenersNewFocus
    listenersEditFocus
    listenersDeleteFocus
    listenersInfoFocus
    
)

type textBoxModel struct {
	text string
}

type MainModel struct {
	state               sessionState
	focus               viewFocus
	mainAgentsBox      textBoxModel
	mainListenersBox    textBoxModel
	mainCLIBox         textBoxModel
	infoBox             textBoxModel
    listenersNewBox     textBoxModel      
    listenersEditBox    textBoxModel
    listenersDeleteBox  textBoxModel
    listenersInfoBox    textBoxModel

}

func NewModel() MainModel {
	m := MainModel{
		state:                  mainState,
		focus:                  mainAgentsFocus,
		mainAgentsBox:          textBoxModel{text: "Agents"},
		mainListenersBox:       textBoxModel{text: "Listeners"},
		mainCLIBox:             textBoxModel{text: "cli"},
		infoBox:                textBoxModel{text: "TODO infoBox.text"},
        listenersNewBox:        textBoxModel{text: "New"},   
        listenersEditBox:       textBoxModel{text: "Edit"},
        listenersDeleteBox:     textBoxModel{text: "Delete"},
        listenersInfoBox:       textBoxModel{text: "View"},
	}
	return m
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		//cmd tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
            switch m.state {
            case mainState:
                return m, tea.Quit
            default:
                m.state = mainState
                m.focus = mainAgentsFocus
            }
			
		case "tab":
			m.focus = m.NextFocus()
        case "enter", "n":
            m.state, m.focus = m.NextState()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m MainModel) NextFocus() viewFocus {
    s := m.state
    f := m.focus
	switch s {
	case mainState:
		switch m.focus {
		case mainAgentsFocus:
			f = mainListenersFocus
		case mainListenersFocus:
			f = mainCLIFocus
		case mainCLIFocus:
			f = mainAgentsFocus
		}
    case listenersState:
        switch m.focus {
        case listenersNewFocus:
            f = listenersEditFocus
        case listenersEditFocus: 
            f = listenersInfoFocus
        case listenersInfoFocus:
            f = listenersDeleteFocus
        case listenersDeleteFocus:
            f = listenersNewFocus
        }
	}
	return f
}

func (m MainModel) NextState() (sessionState, viewFocus) {
    s := m.state
    f := m.focus
    switch f {
    case mainAgentsFocus:
        s = mainAgentsState
    case mainListenersFocus:
        s = listenersState
        f = listenersNewFocus
    case mainCLIFocus:
        s = mainCLIState
    case listenersNewFocus:
        s = listenerNewState
    case listenersEditFocus:
        s = listenerEditState
    case listenersDeleteFocus:
        fmt.Println("[!] TODO: Delete listener.")
    case listenersInfoFocus:
        s = listenerInfoState
    }
    return s, f
} 

func (m MainModel) View() string {
	switch m.state {
	case mainState:
		return m.mainView()
    case listenersState:
        return m.listenersView()
	}
	return ""
}

func (m MainModel) mainView() string {
	var buttons string
	switch m.focus {
	case mainAgentsFocus:
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, focusTextBoxStyle.Render(m.mainAgentsBox.text),
			textBoxStyle.Render(m.mainListenersBox.text), textBoxStyle.Render(m.mainCLIBox.text))
	case mainListenersFocus:
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, textBoxStyle.Render(m.mainAgentsBox.text),
			focusTextBoxStyle.Render(m.mainListenersBox.text), textBoxStyle.Render(m.mainCLIBox.text))
	case mainCLIFocus:
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, textBoxStyle.Render(m.mainAgentsBox.text),
			textBoxStyle.Render(m.mainListenersBox.text), focusTextBoxStyle.Render(m.mainCLIBox.text))
	}
	s := lipgloss.JoinVertical(lipgloss.Top, m.getHeader(headerText), infoBoxStyle.Render(m.infoBox.text), m.getFooter(), buttons)
	return s
}

func (m MainModel) listenersView() string {
	var buttons string
	switch m.focus {
	case listenersNewFocus:
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, focusTextBoxStyle.Render(m.listenersNewBox.text),
			textBoxStyle.Render(m.listenersEditBox.text), textBoxStyle.Render(m.listenersInfoBox.text), textBoxStyle.Render(m.listenersDeleteBox.text))
	case listenersEditFocus:
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, textBoxStyle.Render(m.listenersNewBox.text),
            focusTextBoxStyle.Render(m.listenersEditBox.text), textBoxStyle.Render(m.listenersInfoBox.text), textBoxStyle.Render(m.listenersDeleteBox.text))
	case listenersInfoFocus:
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, textBoxStyle.Render(m.listenersNewBox.text),
			textBoxStyle.Render(m.listenersEditBox.text), focusTextBoxStyle.Render(m.listenersInfoBox.text), textBoxStyle.Render(m.listenersDeleteBox.text))
    case listenersDeleteFocus:
        buttons = lipgloss.JoinHorizontal(lipgloss.Top, textBoxStyle.Render(m.listenersNewBox.text),
            textBoxStyle.Render(m.listenersEditBox.text), textBoxStyle.Render(m.listenersInfoBox.text), focusTextBoxStyle.Render(m.listenersDeleteBox.text))
	}

	s := lipgloss.JoinVertical(lipgloss.Top, m.getHeader(headerText), baseTableStyle.Render(m.getDemoTable().View()), m.getFooter(), buttons)
	return s
}

func (m MainModel) getHeader(s string) string {
    lline := strings.Repeat(borderChar, 3)
	text := headerStyle.Render(s)
	rline := strings.Repeat(borderChar, maxWidth - (len(lline) + len(headerText) + 4))
	return lipgloss.JoinHorizontal(lipgloss.Center, lline, text, rline)
}

func (m MainModel) getFooter() string {
	line := strings.Repeat(borderChar, maxWidth)
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}

func (m MainModel) getDemoTable() table.Model {
    numCol := 4
    tWidth := (maxWidth / numCol) - ((2 * numCol) / numCol)
    tHeight := maxHeight - 4

	columns := []table.Column{
		{Title: "Rank",         Width: tWidth},
		{Title: "City",         Width: tWidth},
		{Title: "Country",      Width: tWidth},
		{Title: "Population",   Width: tWidth},
	}

	rows := []table.Row{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
		{"7", "Cairo", "Egypt", "21,750,020"},
		{"8", "Beijing", "China", "21,333,332"},
		{"9", "Mumbai", "India", "20,961,472"},
		{"10", "Osaka", "Japan", "19,059,856"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tHeight),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(primaryColor)).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(primaryColor)).
		Background(lipgloss.Color(focusColor)).
		Bold(false)
	t.SetStyles(s)

	return t
}
