package listener

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	listenersView	sessionState = iota
	mainView
)

var (
	baseTableStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	FocusedBaseTableStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))
	textBoxModelStyle = lipgloss.NewStyle().
		Width(15).
		Height(5).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.HiddenBorder())
	FocusedTextBoxModelStyle = lipgloss.NewStyle().
		Width(15).
		Height(5).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))
	helpTextBoxModelStyle = lipgloss.NewStyle().
		Width(59).
		Height(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type textBoxModel struct {
	Text	string	
}

type ListenersModel struct {
	state   		sessionState
	Focus			string
	Back			bool
	table			table.Model
	backBox   		textBoxModel
	newBox			textBoxModel
	HelpBox			textBoxModel
}

func NewModel() ListenersModel {
	m := ListenersModel{state: listenersView, Focus: "table", Back: false}
	m.table		= getDemoTable()
	m.backBox		= textBoxModel{Text: "Back"}
	m.newBox 		= textBoxModel{Text: "New"}
	m.HelpBox	= textBoxModel{Text: "Display and configure listener information."}
	return m
}

func (m ListenersModel) Init() tea.Cmd {
	return nil
}

func (m ListenersModel) Update(msg tea.Msg) (ListenersModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.Focus {
			case "table":
				m.Focus = "back"
				m.HelpBox.Text = "Go back?"
			case "new":
				m.Focus = "table"
				m.table.Focus()
				m.HelpBox.Text = "Select a listener to configure."
			default:
				m.Focus = "new"
				m.HelpBox.Text = "Create a new listener."
			}

		case "enter", "n":
			switch m.Focus {
			case "table":
				m.HelpBox.Text = fmt.Sprintf("TODO: Configure %s", m.table.SelectedRow())
			case "new":
				m.HelpBox.Text = "TODO: Create new listener."
			default:
				m.Back = true
				m.HelpBox.Text = "Are you sure?"
			}
		}

		if m.Focus == "table" {
			newModel, newCmd := m.table.Update(msg)
			m.table = newModel
			cmd = newCmd
		}
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ListenersModel) View() string {
	var s string
	var butBox string
	var ctrlBox string
	switch m.Focus{
	case "table":
		butBox = lipgloss.JoinVertical(lipgloss.Top, textBoxModelStyle.Render(m.newBox.Text), textBoxModelStyle.Render(m.backBox.Text) )
		ctrlBox = lipgloss.JoinHorizontal(lipgloss.Top, butBox, FocusedBaseTableStyle.Render(m.table.View()) + "\n")
	case "new":
		butBox = lipgloss.JoinVertical(lipgloss.Top, FocusedTextBoxModelStyle.Render(m.newBox.Text), textBoxModelStyle.Render(m.backBox.Text) )
		ctrlBox = lipgloss.JoinHorizontal(lipgloss.Top, butBox, baseTableStyle.Render(m.table.View()) + "\n")
	default:
		butBox = lipgloss.JoinVertical(lipgloss.Top, textBoxModelStyle.Render(m.newBox.Text), FocusedTextBoxModelStyle.Render(m.backBox.Text) )
		ctrlBox = lipgloss.JoinHorizontal(lipgloss.Top, butBox, baseTableStyle.Render(m.table.View()) + "\n")
	}
	s += lipgloss.JoinVertical(lipgloss.Top, helpTextBoxModelStyle.Render(m.HelpBox.Text), ctrlBox)
	s += helpStyle.Render(fmt.Sprintf("\ntab: Focus next • n: select %s • q: exit\n", m.Focus))
	return s
}

func getDemoTable() table.Model {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
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
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("69")).
		Bold(false)
	t.SetStyles(s)

	return t
}