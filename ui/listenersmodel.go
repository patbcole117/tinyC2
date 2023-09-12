package ui

import (
	tea "github.com/charmbracelet/bubbletea"
 	"github.com/charmbracelet/lipgloss"
 	"github.com/charmbracelet/bubbles/table"
)

type TableModel struct {
	focus int
	buttons []button
	table	table.Model
}

func NewTableModel(butt []button) TableModel {
	return TableModel {
		focus: 0,
		buttons: butt,
		table: GetDemoTable(),
	}
}

func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "esc":
            return m, toRootState
        case "enter":
            return m, m.buttons[m.focus].do
        case "tab":
            m.focus = NextFocus(m.focus, len(m.buttons))   
        }
    }
    m.table, cmd = m.table.Update(msg)
    cmds = append(cmds, cmd)
    return m, tea.Batch(cmds...)
}

func (m TableModel) View() string {
    b := GetButtonViewComponent(m.buttons, m.focus)
    return lipgloss.JoinVertical(lipgloss.Top,
        GetHeaderViewComponent(),
        tableStyle.Render(m.table.View()),
        GetFooterViewComponent(), b)
}
