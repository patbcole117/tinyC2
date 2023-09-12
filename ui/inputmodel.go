package ui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textinput"
)

type InputModel struct {
    focus   int
    buttons []button
    labels  []string
    inputs  []textinput.Model
}

func NewInputModel(labs, placeholders []string, f [](func() tea.Msg)) InputModel {
    var inModels []textinput.Model = make([]textinput.Model, len(labs))
    butt := []button {
        {text: "Save", do: f[0]},
        {text: "Cancel", do: f[1]},
    }
    
    for i := range labs {
        inModels[i] = textinput.New()
        inModels[i].Placeholder = placeholders[i]
        inModels[i].CharLimit   = 25
        inModels[i].Width       = 25
        inModels[i].Prompt      = "# "
    }
    return InputModel {
        focus:      0,
        buttons:    butt,
        labels:     labs,
        inputs:     inModels,
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
            m.inputs[x].Prompt = "> "
            m.inputs[x].Focus() 
        } else {
            m.inputs[x].Prompt = "# "
            m.inputs[x].Blur()
        }
    }
    
    for i := range m.inputs {
        m.inputs[i], cmd = m.inputs[i].Update(msg)
        cmds = append(cmds, cmd)
    }
    return m, tea.Batch(cmds...)
}

func (m InputModel) View() string {
    return lipgloss.JoinVertical(lipgloss.Top,
        GetHeaderViewComponent(),
        GetInputViewComponent(m.labels, m.inputs),
        GetFooterViewComponent(),
        GetButtonViewComponent(m.buttons, m.focus))
}
