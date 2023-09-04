package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/patbcole117/tinyC2/listener"
)

type sessionState uint

const (
	agentsView		sessionState = iota
	listenersView
	mainView
)

var (
	textBoxModelStyle = lipgloss.NewStyle().
		Width(15).
		Height(5).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.HiddenBorder())
	focusedtextBoxModelStyle = lipgloss.NewStyle().
		Width(15).
		Height(5).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))
	rightDisplayTextBoxModelStyle = lipgloss.NewStyle().
		Width(100).
		Height(19).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type textBoxModel struct {
	text	string	
}

type mainModel struct {
	state   			sessionState
	focus				string
	rightDisplay 		textBoxModel
	agentsBox		textBoxModel
	listenersBox		textBoxModel
	exitBox 			textBoxModel
	listenersModel		listener.ListenersModel
}

func newModel() mainModel {
	m := mainModel{state: mainView, focus: "agents"}
	m.rightDisplay		= textBoxModel{text: "TODO: ASCII art."}
	m.agentsBox 		= textBoxModel{text: "Agents"}
	m.listenersBox	= textBoxModel{text: "Listeners"}
	m.exitBox		= textBoxModel{text: "Exit"}
	m.listenersModel	= listener.NewModel()
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case listenersView:
			if !m.listenersModel.Back {
				newListenersModel, newCmd := m.listenersModel.Update(msg)
				m.listenersModel = newListenersModel
				cmd = newCmd
			} else {
				m.state = mainView
				m.listenersModel.Back = false
				m.listenersModel.Focus = "table"
				m.listenersModel.HelpBox.Text = "Select a listener to configure."
				m.focus = "agents"
			}
		default:
			m.state = mainView
			switch msg.String() {
			case "tab":
				switch m.focus {
				case "agents":
					m.focus = "listeners"
					m.rightDisplay.text = "Display and configure listeners."
				case "listeners":
					m.focus = "exit"
					m.rightDisplay.text = "Exit."
				default:
					m.focus = "agents"
					m.rightDisplay.text = "Display and configure agents."
	
				}
	
			case "enter", "n":
				switch m.focus {
				case "agents":
					m.state = agentsView
				case "listeners":
					m.state = listenersView
				case "exit":
					m.rightDisplay.text = "Goodbye!"
					return m, tea.Quit
	
				}
			}
		}

		switch msg.String() {
		case "ctrl+c", "esc", "q":
			if m.state == mainView {
				m.rightDisplay.text = "Goodbye!"
				return m, tea.Quit
			}
				m.state = mainView
		}
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	switch m.state {
	case agentsView:
		//TODO
		s = m.MainView()
	case listenersView:
		s = m.listenersModel.View()
	default:
		s = m.MainView()
	}
	return s

}

func (m mainModel) MainView() string {
	var s string
	var lbox string
	if m.focus == "agents" {
		lbox = lipgloss.JoinVertical(lipgloss.Top, focusedtextBoxModelStyle.Render(m.agentsBox.text),
		textBoxModelStyle.Render(m.listenersBox.text), textBoxModelStyle.Render(m.exitBox.text))
	} else if m.focus == "listeners" {
		lbox = lipgloss.JoinVertical(lipgloss.Top, textBoxModelStyle.Render(m.agentsBox.text),
		focusedtextBoxModelStyle.Render(m.listenersBox.text),	textBoxModelStyle.Render(m.exitBox.text))
	} else {
		lbox = lipgloss.JoinVertical(lipgloss.Top, textBoxModelStyle.Render(m.agentsBox.text),
		textBoxModelStyle.Render(m.listenersBox.text), focusedtextBoxModelStyle.Render(m.exitBox.text))
	}
	s += lipgloss.JoinHorizontal(lipgloss.Top, lbox, rightDisplayTextBoxModelStyle.Render(m.rightDisplay.text))
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: select %s • q: exit\n", m.focus))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == agentsView {
		return "agents"
	} else if m.state == listenersView {
		return "listeners"
	} else {
		return "exit"
	}
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}