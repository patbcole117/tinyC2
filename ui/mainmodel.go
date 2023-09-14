package ui

import (
	"log"
    "fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state uint
const (
    configState         state = iota
    listenerNewState
    listenersState
    rootState
)

type apiConfig struct {
	apiIp	string
	apiPort	string
	apiVer  string
}

type MainModel struct {
	config		    apiConfig
	state		    state
    infoMsg         string
	rootModel   	RootModel
    tableModel      TableModel
    inputModel      InputModel
}

func NewMainModel() MainModel {
	return MainModel {
		state:	        rootState,
		config:         apiConfig{apiIp: "127.0.0.1", apiPort: "8000", apiVer: "v1"},
        infoMsg:        defaultinfoMsg,
		rootModel:      NewRootModel(),
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
            m.config.apiIp      = m.inputModel.inputs[0].textBox.Value()
            m.config.apiPort    = m.inputModel.inputs[1].textBox.Value()
            m.config.apiVer     = m.inputModel.inputs[2].textBox.Value()
            m.state = rootState
            msg := fmt.Sprintf(`{"CONFIG":{"Ip": "%s", "Port":"%s", "Ver":"%s"}}`, m.config.apiIp, m.config.apiPort, m.config.apiVer)
            m.infoMsg = msg
        }
    case inputCancelMsg:
        switch msg {
        case "Config":
            m.state = rootState 
        }
	case setStateMsg:
		switch msg {
        case "Config":
            m.state = configState
            ins, butts := m.MakeInputModelComponents()
            m.inputModel = NewInputModel(ins, butts)
            
        case "NewListener":
            m.state = listenerNewState
            m.infoMsg = "Create a new listener."
            ins, butts := m.MakeInputModelComponents()
            m.inputModel = NewInputModel(ins, butts)
            
		case "Listeners":
            butt := []button {
                {text: "New", do: toNewListenerState},
                {text: "Edit", do: TODOButton},
                {text: "Info", do: TODOButton},
                {text: "Start", do: TODOButton},
                {text: "Stop", do: TODOButton},
                {text: "Delete", do: TODOButton},
            }
            cmds = append(cmds, trigGetAllListeners)
            m.tableModel = NewTableModel(butt)
            m.tableModel.table.SetCursor(0)
            m.infoMsg = "View and configure listeners."
			m.state = listenersState
		default:
            m.infoMsg = defaultinfoMsg
			m.state = rootState
		}
    case trigNewListenerMsg:
            cmd = NewListener(m.inputModel.inputs[0].textBox.Value(), m.inputModel.inputs[1].textBox.Value(), m.inputModel.inputs[2].textBox.Value(), m.config)
            cmds = append(cmds, cmd)
            m.state = listenersState

    case newInfoMsg:
        m.infoMsg = string(msg)
	}

    if m.state != curState {
        m.tableModel.table.Blur()
        m.tableModel.focus = 0
        m.rootModel.bigBox = GetRandomBanner()
        m.rootModel.focus = 0
    }

	switch m.state {
    case configState:
        m.inputModel, cmd = m.inputModel.Update(msg)
        cmds = append(cmds, cmd)
    case listenerNewState:
        m.inputModel, cmd = m.inputModel.Update(msg)
        cmds = append(cmds, cmd)
    case listenersState:
        m.tableModel.table.Focus()
        m.tableModel, cmd = m.tableModel.Update(msg)
        cmds = append(cmds, cmd)
	default:
		m.rootModel, cmd = m.rootModel.Update(msg)
        cmds = append(cmds, cmd)

	}
    
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var s string
	switch m.state {
    case listenerNewState, configState:
        s = m.inputModel.View()
	case listenersState:
		s = m.tableModel.View()
	default:
		s = m.rootModel.View()
	}
	return lipgloss.JoinVertical(lipgloss.Top, s, helpBoxStyle.Render(m.infoMsg))
}

func (m MainModel) MakeInputModelComponents() ([]input, []button) {
    var butts []button
    var inputs []input
    switch m.state {
    case configState:
        butts = []button {
            {text: "Save", do: saveConfig},
            {text: "Cancel", do: cancelConfig},
        }
        inputs = []input {
            {label: "Ip", textBox: textinput.New()},
            {label: "Port", textBox: textinput.New()},
            {label: "Version", textBox: textinput.New()},
        }
        p := []string{m.config.apiIp, m.config.apiPort, m.config.apiVer}
        for i := range inputs {
            inputs[i].textBox.Placeholder = p[i]
            inputs[i].textBox.CharLimit   = 25
            inputs[i].textBox.Width       = 25
            inputs[i].textBox.Prompt      = "# "
        }
    case listenerNewState:
        butts = []button {
            {text: "Save", do: trigNewListener},
            {text: "Cancel", do: TODOButton},
        }
        inputs = []input {
            {label: "Name", textBox: textinput.New()},
            {label: "Ip", textBox: textinput.New()},
            {label: "Port", textBox: textinput.New()},
        }
        p := []string{
            m.tableModel.table.SelectedRow()[1],
            m.tableModel.table.SelectedRow()[2],
            m.tableModel.table.SelectedRow()[3],
        }
        for i := range inputs {
            inputs[i].textBox.Placeholder = p[i]
            inputs[i].textBox.CharLimit   = 25
            inputs[i].textBox.Width       = 25
            inputs[i].textBox.Prompt      = "# "
        }
    }
    return inputs, butts
}

func KickOff() {
	p := tea.NewProgram(NewMainModel())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
