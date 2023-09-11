package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type state uint
const (
	listenersState	state = iota
	rootState
)

type apiConfig struct {
	apiIp	string
	apiPort	string
	apiVer  string
}

type MainModel struct {
	config		apiConfig
	state		state
	rootModel	RootModel
}

func NewMainModel() MainModel {
	return MainModel {
		state:	rootState,
		config: apiConfig{apiIp: "127.0.0.1", apiPort: "8000", apiVer: "v1"},
		rootModel: NewRootModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case stateChange:
		switch msg {
		case "Listeners":
			m.state = listenersState
		default:
			m.state = rootState
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
            return m, tea.Quit
		}
	}

	switch m.state {
	case rootState:
		m.rootModel, cmd = m.rootModel.Update(msg)
	}

	
    cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var s string
	switch m.state {
	case listenersState:
		m.rootModel.bigBox = "ListenersState"
		s = m.rootModel.View()
	default:
		s = m.rootModel.View()
	}
	return s
}

func KickOff() {
	p := tea.NewProgram(NewMainModel())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
