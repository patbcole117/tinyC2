package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type apiConfig struct {
	apiIp	string
	apiPort	string
	apiVer  string
}

type MainModel struct {
	config		    apiConfig
	state		    state
    helpMsg         string
	rootModel   	RootModel
    listenersModel  ListenersModel
    inputModel      InputModel
}

func NewMainModel() MainModel {
	return MainModel {
		state:	        rootState,
		config:         apiConfig{apiIp: "127.0.0.1", apiPort: "8000", apiVer: "v1"},
        helpMsg:        defaultHelpMsg,
		rootModel:      NewRootModel(),
        listenersModel: NewListenersModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
    curState := m.state
	switch msg := msg.(type) {
    case inputSaveMsg:
        switch msg {
        case "Config":
            m.config.apiIp      = m.inputModel.inputs[0].Value()
            m.config.apiPort    = m.inputModel.inputs[1].Value()
            m.config.apiVer     = m.inputModel.inputs[2].Value()
            m.state = rootState
        }
    case inputCancelMsg:
        switch msg {
        case "Config":
            m.state = rootState 
        }
	case setStateMsg:
		switch msg {
        case "Config":
            l := []string{"Ip", "Port", "Version"}
            p := []string{m.config.apiIp, m.config.apiPort, m.config.apiVer}
            f := [](func() tea.Msg){saveConfig, cancelConfig}
            m.inputModel = NewInputModel(l, p, f)
            m.state = configState
        case "NewListener":
            l := []string{"Name", "Ip", "Port"}
            p := []string{m.listenersModel.table.SelectedRow()[1],
                m.listenersModel.table.SelectedRow()[2],
                m.listenersModel.table.SelectedRow()[3],
            }
            f := [](func() tea.Msg){TODOButton, TODOButton}
            m.inputModel = NewInputModel(l, p, f)
            m.helpMsg = "Create a new listener."
            m.state = listenerNewState
		case "Listeners":
            m.listenersModel.table.SetCursor(0)
            m.helpMsg = "View and configure listeners."
			m.state = listenersState
		default:
            m.helpMsg = defaultHelpMsg
			m.state = rootState
		}
    case setHelpMsg:
        m.helpMsg = string(msg)
	}

    if m.state != curState {
        m.listenersModel.table.Blur()
        m.listenersModel.focus = 0
        m.rootModel.bigBox = GetRandomBanner()
        m.rootModel.focus = 0
    }

	switch m.state {
    case configState:
        m.inputModel, cmd = m.inputModel.Update(msg)
    case listenerNewState:
        m.inputModel, cmd = m.inputModel.Update(msg)
    case listenersState:
        m.listenersModel.table.Focus()
        m.listenersModel, cmd = m.listenersModel.Update(msg)
	default:
		m.rootModel, cmd = m.rootModel.Update(msg)
	}

	
    cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var s string
	switch m.state {
    case listenerNewState, configState:
        s = m.inputModel.View()
	case listenersState:
		s = m.listenersModel.View()
	default:
		s = m.rootModel.View()
	}
	return lipgloss.JoinVertical(lipgloss.Top, s, helpBoxStyle.Render(m.helpMsg))
}

func KickOff() {
	p := tea.NewProgram(NewMainModel())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
