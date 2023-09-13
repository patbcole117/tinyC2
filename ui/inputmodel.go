package ui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type InputModel struct {
    focus   int
    buttons []button
    inputs  []input
}

func NewInputModel(ins []input, butts []button) InputModel {
    return InputModel {
        focus:      0,
        buttons:    butts,
        inputs:     ins,
    }
}

func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "esc":
            return m, m.buttons[1].do
        case "enter":
            if m.focus >= len(m.buttons) {
                m.focus = NextFocus(m.focus, len(m.buttons) + len(m.inputs))
            } else {
                cmds = append(cmds, m.buttons[m.focus].do)
            }
        case "tab":
                m.focus = NextFocus(m.focus, len(m.buttons) + len(m.inputs))
        }
    }
        
    inFoc :=  m.focus - len(m.buttons)
    for x := range m.inputs {
        if x == inFoc {
            m.inputs[x].textBox.Prompt = "> "
            m.inputs[x].textBox.Focus() 
        } else {
            m.inputs[x].textBox.Prompt = "# "
            m.inputs[x].textBox.Blur()
        }
    }
    
    for i := range m.inputs {
        m.inputs[i].textBox, cmd = m.inputs[i].textBox.Update(msg)
        cmds = append(cmds, cmd)
    }
    return m, tea.Batch(cmds...)
}

func (m InputModel) View() string {
    return lipgloss.JoinVertical(lipgloss.Top,
        GetHeaderViewComponent(),
        GetInputViewComponent(m.inputs),
        GetFooterViewComponent(),
        GetButtonViewComponent(m.buttons, m.focus))
}