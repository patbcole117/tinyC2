package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	focus	int
	buttons	[]button
	bigBox	string
}

func NewRootModel() RootModel {
	butt := []button {
        {text: "Agents", do: TODOButton},
        {text: "Nodes",	do: toNodesState},
		{text: "cli", do: TODOButton},
		{text: "Config", do: toConfigState},
    }
	return RootModel {
		focus:		0,
		buttons: 	butt,
		bigBox: 	GetRandomBanner(),
	}
}

func (m RootModel) Update(msg tea.Msg) (RootModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
        case "ctrl+c", "esc":
            return m, tea.Quit
		case "enter":
			cmds = append(cmds, m.buttons[m.focus].do)
		case "tab":
			m.focus = NextFocus(m.focus, len(m.buttons))
		}
	}
	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	b := GetButtonViewComponent(m.buttons, m.focus)
	return lipgloss.JoinVertical(lipgloss.Top,
		GetHeaderViewComponent(),
		bigBoxStyle.Render(m.bigBox),
		GetFooterViewComponent(), b)
}
