package ui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/patbcole117/tinyC2/node"
)

type state uint
const (
    configState         state = iota
    nodesState
    nodeEditState
    nodeInfoState
    nodeNewState
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
    nodes           []node.Node
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
	return SyncNodes(m.config)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
    curState := m.state
	switch msg := msg.(type) {
    case inputSaveMsg:
        switch msg {
        case "Config":
            cmds = append(cmds, toRootState)
            m.config.apiIp      = m.inputModel.inputs[0].textBox.Value()
            m.config.apiPort    = m.inputModel.inputs[1].textBox.Value()
            m.config.apiVer     = m.inputModel.inputs[2].textBox.Value()
            infoMsg := infMsg("CONFIG", fmt.Sprintf(`"Ip": "%s", "Port":"%s", "Ver":"%s"`,
                        m.config.apiIp, m.config.apiPort, m.config.apiVer))
            cmds = append(cmds, setInfoMsg(infoMsg))
        }
    case inputCancelMsg:
        switch msg {
        case "Config":
            cmds = append(cmds, toRootState)
        }
    case newInfoMsg:
        time.Sleep(500 * time.Millisecond)
        m.infoMsg = string(msg)
    case trigNewNodeMsg:
        cmds = append(cmds, toNodesState)
        cmd = NewNode(m.inputModel.inputs[0].textBox.Value(), m.inputModel.inputs[1].textBox.Value(),
                    m.inputModel.inputs[2].textBox.Value(), m.config)
        cmds = append(cmds, cmd)
    case trigDeleteNodeMsg:
        var infoMsg string
        if len(m.nodes) == 0 {
            infoMsg := "No nodes."
            cmds = append(cmds, setInfoMsg(infoMsg))
        } else {
            id := m.tableModel.table.SelectedRow()[0]
            for i := range m.nodes {
                if m.nodes[i].Id == id {
                    err := m.nodes[i].SrvStop(); if err != nil {
                        infoMsg = errMsg("trigToggleNode:SrvStop", m.nodes[i].Id)
                    }
                    cmds = append(cmds, DeleteNode(m.tableModel.table.SelectedRow()[0], m.config))
                }
            }
        }
        cmds = append(cmds, setInfoMsg(infoMsg))
    case trigToggleNodeMsg:
        var infoMsg string
        if len(m.nodes) == 0 {
            infoMsg = "No nodes."
        } else {
            id := m.tableModel.table.SelectedRow()[0]
            for i := range m.nodes {
                if m.nodes[i].Id == id {
                    if string(msg) == "START"{
                        err := m.nodes[i].SrvStart(); if err != nil {
                            infoMsg = errMsg("trigToggleNode:SrvStart",err.Error())
                        }
                    } else {
                        err := m.nodes[i].SrvStop(); if err != nil {
                            infoMsg = errMsg("trigToggleNode:SrvStop", err.Error())
                        }
                    }
                    cmds = append(cmds, UpdateNode(m.nodes[i], m.config))
                }
            }  
        }
        cmds = append(cmds, setInfoMsg(infoMsg))
    case trigUpdateNodeMsg:
        cmds = append(cmds, toNodesState)
        var infoMsg string
        if len(m.nodes) == 0 {
            infoMsg := "No nodes."
            cmds = append(cmds, setInfoMsg(infoMsg))
        } else {
            id := m.tableModel.table.SelectedRow()[0]
            for i := range m.nodes {
                if m.nodes[i].Id == id {
                    if m.inputModel.inputs[0].textBox.Value() != "" {
                        m.nodes[i].Name = m.inputModel.inputs[0].textBox.Value()
                    }
                    if m.inputModel.inputs[1].textBox.Value() != "" {
                        m.nodes[i].Ip = m.inputModel.inputs[1].textBox.Value()
                    }
                    if m.inputModel.inputs[2].textBox.Value() != "" {
                        m.nodes[i].Port, _ = strconv.Atoi(m.inputModel.inputs[2].textBox.Value())
                    }
                    err := m.nodes[i].SrvStop(); if err != nil {
                        infoMsg = errMsg("trigToggleNode:SrvStop", m.nodes[i].Id)
                    }
                    cmds = append(cmds, UpdateNode(m.nodes[i], m.config))
                }
            }
        }
        cmds = append(cmds, setInfoMsg(infoMsg))
    case syncNodesMsg:
        for j := range msg {
            for i := range m.nodes      {
                if msg[j].Id == m.nodes[i].Id {
                    msg[j].Server = m.nodes[i].Server
                    if msg[j].Ip != m.nodes[i].Ip || msg[j].Port != m.nodes[i].Port {
                        infoMsg := infMsg("REBOOT REQUIRED", m.nodes[i].Id)
                        cmds = append(cmds, setInfoMsg(infoMsg))
                    }
                }
            }
        }
        if len(msg) != len(m.nodes) {
            infoMsg := infMsg("NODE SYNC", fmt.Sprintf("%d -> %d", len(m.nodes), len(msg)))
            cmds = append(cmds, setInfoMsg(infoMsg))
        }
        m.nodes = msg
    case setStateMsg:
		switch msg {
        case "Config":
            m.state = configState
            ins, butts := m.MakeInputModelComponents()
            m.inputModel = NewInputModel(ins, butts)
		case "Nodes":
            m.state = nodesState
            cmds = append(cmds, SyncNodes(m.config))
            t, butts := m.MakeTableModelComponents()
            m.tableModel = NewTableModel(butts, t)
        case "NodesNew":
            m.state = nodeNewState
            m.infoMsg = "Create a new node."
            ins, butts := m.MakeInputModelComponents()
            m.inputModel = NewInputModel(ins, butts)
        case "NodesInfo":
            if len(m.nodes) != 0 {
                m.state = nodeInfoState
                m.tableModel.table.SelectedRow()
                ins, butts := m.MakeInputModelComponents()
                m.inputModel = NewInputModel(ins, butts)
            } else {
                cmds = append(cmds, setInfoMsg("There are no nodes."))
            }
        case "NodesEdit":
            m.state = nodeEditState
            ins, butts := m.MakeInputModelComponents()
            m.inputModel = NewInputModel(ins, butts)
		default:
            m.state = rootState
            cmds = append(cmds, setInfoMsg(defaultinfoMsg))
		}
	}

    if m.state != curState {
        m.tableModel.table.Blur()
        m.tableModel.focus = 0
        m.rootModel.bigBox = GetRandomBanner()
        m.rootModel.focus = 0
    }

	switch m.state {
    case configState, nodeEditState, nodeInfoState, nodeNewState:
        m.inputModel, cmd = m.inputModel.Update(msg)
        cmds = append(cmds, cmd)
    case nodesState:
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
    case nodeEditState, nodeInfoState, nodeNewState, configState:
        s = m.inputModel.View()
	case nodesState:
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
    case nodeEditState:
        butts = []button {
            {text: "Save", do: trigUpdateNode},
            {text: "Cancel", do: toNodesState},
        }
        inputs = []input {
            {label: "Name", textBox: textinput.New()},
            {label: "Ip", textBox: textinput.New()},
            {label: "Port", textBox: textinput.New()},
        }
        var p []string
        if len(m.nodes) == 0 {
            p = []string {"Name", "Ip", "Port"}
        } else {
            p = []string{
                m.tableModel.table.SelectedRow()[1],
                m.tableModel.table.SelectedRow()[2],
                m.tableModel.table.SelectedRow()[3],
            }
        }
        for i := range inputs {
            inputs[i].textBox.Placeholder = p[i]
            inputs[i].textBox.CharLimit   = 25
            inputs[i].textBox.Width       = 25
            inputs[i].textBox.Prompt      = "# "
        }
    case nodeInfoState:
        var n node.Node
        for _, v := range m.nodes {
            if v.Id == m.tableModel.table.SelectedRow()[0]{
                n = v
                break
            }
        }
        butts = []button {
            {text: "Return", do: toNodesState},
        }
        inputs = []input {
            {label: "_id", textBox: textinput.New()},
            {label: "Name", textBox: textinput.New()},
            {label: "Ip", textBox: textinput.New()},
            {label: "Port", textBox: textinput.New()},
            {label: "Status", textBox: textinput.New()},
            {label: "Dob", textBox: textinput.New()},
            {label: "Hello", textBox: textinput.New()},
        }
        p := []string{n.Id, n.Name, n.Ip, strconv.Itoa(n.Port), strconv.Itoa(n.Status),
            n.Dob.Format(time.RFC3339), n.Hello.Format(time.RFC3339),}
        for i := range inputs {
            inputs[i].textBox.Placeholder = p[i]
            inputs[i].textBox.CharLimit   = 50
            inputs[i].textBox.Width       = 50
            inputs[i].textBox.Prompt      = "# "
        }
    case nodeNewState:
        butts = []button {
            {text: "Save", do: trigNewNode},
            {text: "Cancel", do: toNodesState},
        }
        inputs = []input {
            {label: "Name", textBox: textinput.New()},
            {label: "Ip", textBox: textinput.New()},
            {label: "Port", textBox: textinput.New()},
        }
        var p []string
        if len(m.nodes) == 0 {
            p = []string {"Name", "Ip", "Port"}
        } else {
            p = []string{
                m.tableModel.table.SelectedRow()[1],
                m.tableModel.table.SelectedRow()[2],
                m.tableModel.table.SelectedRow()[3],
            }
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

func (m MainModel) MakeTableModelComponents() (table.Model, []button) {
    butts := []button {
        {text: "New", do: toNodesNewState},
        {text: "Edit", do: toNodesEditState},
        {text: "Info", do: toNodesInfoState},
        {text: "Start", do: trigStartNode},
        {text: "Stop", do: trigStopNode},
        {text: "Delete", do: trigDeleteNode},
    }
    
    t := m.GetNodesTable()
    return t, butts
}

func (m MainModel) GetNodesTable() table.Model {
    var t table.Model
	numCol := 5
	tWidth := (maxWidth / numCol) - ((2 * numCol) / numCol)
	tHeight := maxHeight - 4

    cols := []table.Column {
        {Title: "Id", Width: tWidth},
        {Title: "Name", Width: tWidth},
        {Title: "Ip", Width: tWidth},
        {Title: "Port", Width: tWidth},
        {Title: "Status", Width: tWidth},
    }

	var rows []table.Row
	for _, n := range m.nodes {
		r := table.Row {n.Id, n.Name, n.Ip, strconv.Itoa(n.Port), strconv.Itoa(n.Status)}
		rows = append(rows, r)
	}
	t = table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(tHeight),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(secondaryColor)).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(secondaryColor)).
		Background(lipgloss.Color(primaryColor)).
		Bold(false).
		BorderStyle(lipgloss.NormalBorder())
	t.SetStyles(s)

	return t
}

func KickOff() {
	p := tea.NewProgram(NewMainModel())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
