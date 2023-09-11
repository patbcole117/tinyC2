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
        {text: "Agents", do: rootListenersButton},
        {text: "Listeners",	do: rootListenersButton},
		{text: "cli", do: rootListenersButton},
		{text: "Config", do: rootListenersButton},
    }
	return RootModel {
		focus:		0,
		buttons: 	butt,
		bigBox: 	"TODO",
	}
}

type stateChange string
func rootListenersButton() tea.Msg {
	return stateChange("Listeners")
}

func (m RootModel) Update(msg tea.Msg) (RootModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.buttons[m.focus].do
		case "tab":
			m.focus = NextFocus(m.focus, len(m.buttons))
		}
	}

    cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	b := GetButtonViewComponent(m.buttons, m.focus)
	return lipgloss.JoinVertical(lipgloss.Top,
		GetHeaderViewComponent(headerText),
		bigBoxStyle.Render(m.bigBox),
		GetFooterViewComponent(), b)
}