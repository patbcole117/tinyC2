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
	maxWidth        = 75
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

	bigBoxStyle = lipgloss.NewStyle().
        Width(maxWidth).
        Height(maxHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderForeground(lipgloss.Color(primaryColor))

	focusButtonStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(focusColor))

    baseTableStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(primaryColor))

	buttonStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.HiddenBorder())
)

type sessionState uint
const (
	agentsState         sessionState = iota
	agentInfoState
	agentEditState
    cliState
	listenersState
	listenerEditState
    listenerInfoState
    listenerNewState
	mainState
)

type button struct {
    text    string
    state   sessionState
}

var (
    mainButtons = []button {
        {text: "Agents",     state: agentsState},
        {text: "Listeners",  state: listenersState},
        {text: "cli",        state: cliState},
    }

    listenersButtons    = []button {
        {text: "New",       state: listenerNewState},
        {text: "Edit",      state: listenerEditState},
        {text: "Info",      state: listenerInfoState},
        {text: "Delete",    state: listenersState},
    }

    listenerEditButtons = []button {
        {text: "Save",      state: listenersState},
        {text: "Cancel",    state: listenersState},
    }
)

type MainModel struct {
	state           sessionState
	focus           int 
    buttons         []button     
    bigBox          string
    listenersTable  table.Model
}



func NewModel() MainModel {
	m := MainModel{
		state:      mainState,
		focus:      0, 
        buttons:    mainButtons,
        bigBox:     "TODO: bigBox string.",
	}
        m.listenersTable = m.getDemoTable()
	return m
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		cmds []tea.Cmd
	)
    prevState := m.state
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
            switch m.state {
            case mainState:
                return m, tea.Quit
            default:
                m.state     = mainState
                m.buttons   = mainButtons
            }
			
		case "tab":
			m.focus = m.NextFocus()
        case "enter", "n":
            m.state, m.buttons, m.focus = m.NextState()
		}
	}
    if m.state != prevState {
        m.focus = 0
        switch m.state {
        case listenersState:
            m.listenersTable.Focus()
            m.listenersTable.SetCursor(0)
        case listenerEditState:
            m.listenersTable.Blur()
        default:
            m.listenersTable.Blur()
            m.listenersTable.SetCursor(0)
        }
    }
    m.listenersTable, cmd = m.listenersTable.Update(msg)
    cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) NextFocus() int {
    f := m.focus + 1
    if f >= len(m.buttons) {
        return 0
    }
	return f
}

func (m MainModel) NextState() (sessionState, []button, int) {
    f := 0
    s := m.buttons[m.focus].state
    var b []button
    switch s {
    case listenersState:
        b = listenersButtons
    case listenerEditState:
        b = listenerEditButtons
    default:
        b = mainButtons
    }
    if s == m.state {
        f = m.focus
    }
    return s, b, f
} 

func (m MainModel) View() string {
    b := m.getButtonRow()
	switch m.state {
	case mainState:
        return lipgloss.JoinVertical(lipgloss.Top, m.getHeader(headerText),
            bigBoxStyle.Render(m.bigBox), m.getFooter(), b)
    case listenersState:
        return lipgloss.JoinVertical(lipgloss.Top, m.getHeader(headerText),
            baseTableStyle.Render(m.listenersTable.View()), m.getFooter(), b)
    case listenerNewState:
        return fmt.Sprintf("TODO: New listener.")
    case listenerEditState:
        return lipgloss.JoinVertical(lipgloss.Top, m.getHeader(headerText),
            bigBoxStyle.Render(m.listenersTable.SelectedRow()[1]), m.getFooter(), b)
    case listenerInfoState:
        return fmt.Sprintf("TODO: View info of listener %s.",
            m.listenersTable.SelectedRow()[1])
	}
	return ""
}

func (m MainModel) getButtonRow() string {
    var bview string
    for i, b := range m.buttons {
        if i == m.focus {
            bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
            focusButtonStyle.Render(b.text))
        } else {
            bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
            buttonStyle.Render(b.text))
        }
    }
    return bview
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
    numCol := 5
    tWidth := (maxWidth / numCol) - ((2 * numCol) / numCol)
    tHeight := maxHeight - 4

	columns := []table.Column{
		{Title: "Id",       Width: tWidth},
		{Title: "Name",     Width: tWidth},
		{Title: "Ip",       Width: tWidth},
		{Title: "Port",     Width: tWidth},
		{Title: "Status",   Width: tWidth},
	}

	rows := []table.Row{
		{"1", "Tokyo",      "127.0.0.1",    "80", "1"},
		{"2", "Delhi",      "127.0.0.2",    "81", "1"},
		{"3", "Shanghai",   "127.0.0.3",    "82", "1"},
		{"4", "Dhaka",      "127.0.0.4",    "84", "1"},
		{"5", "Sao",        "127.0.0.5",    "85", "1"},
		{"6", "Colombo",    "127.0.0.6",    "86", "1"},
		{"7", "Cairo",      "127.0.0.7",    "87", "1"},
		{"8", "Beijing",    "127.0.0.8",    "88", "1"},
		{"9", "Mumbai",     "127.0.0.9",    "89", "1"},
		{"10", "Osaka",     "127.0.0.10",   "90", "1"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
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
		Bold(false).
        BorderStyle(lipgloss.NormalBorder())
	t.SetStyles(s)

	return t
}
