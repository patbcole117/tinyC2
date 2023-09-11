package ui

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	secondaryColor    = "240"
	primaryColor      = strconv.Itoa(rand.Intn(230))
	maxWidth        = 75
	maxHeight       = 15	
	buttonWidth     = 11
	buttonHeight    = 1
	borderChar      = "-"
    headerBar       = "-"
    headerText      = "tinyC2"
)

var (
	headerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = headerBar
        b.Left = headerBar
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	bigBoxStyle = lipgloss.NewStyle().
        Width(maxWidth).
        Height(maxHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderForeground(lipgloss.Color(secondaryColor))

	buttonStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.HiddenBorder())
	focusButtonStyle = lipgloss.NewStyle().
        Width(buttonWidth).
        Height(buttonHeight).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color(primaryColor))

    baseTableStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(secondaryColor))

	inputLabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(primaryColor)).
		Bold(true)
	inputBoxStyle = lipgloss.NewStyle()
    inputBigBoxStyle = lipgloss.NewStyle().
        Width(maxWidth).
        Height(maxHeight).
        Align(lipgloss.Top, lipgloss.Top).
        BorderForeground(lipgloss.Color(secondaryColor))
)

type sessionState uint
const (
	agentsState         sessionState = iota
	agentInfoState
	agentEditState
    cliState
	configEditState
	listenersState
	listenerEditState
    listenerInfoState
    listenerNewState
	mainState
)

var (
    mainButtons = []button {
        {text: "Agents",     	state: agentsState},
        {text: "Listeners",  	state: listenersState},
		{text: "cli",        	state: cliState},
		{text: "Config", 		state: configEditState},
    }

    listenersButtons    = []button {
        {text: "New",       state: listenerNewState},
        {text: "Edit",      state: listenerEditState},
        {text: "Info",      state: listenerInfoState},
        {text: "Start",     state: listenersState},
        {text: "Stop",      state: listenersState},
        {text: "Delete",    state: listenersState},
    }

    listenerEditButtons = []button {
        {text: "Save",      state: listenersState},
        {text: "Cancel",    state: listenersState},
    }

	configEditButtons = []button {
        {text: "Save",      state: mainState},
        {text: "Cancel",    state: mainState},
    }
)

type Config struct {
	apiIp	string
	apiPort	string
	apiVer  string
}

type MainModel struct {
	state           sessionState
	focus           int 
	config			Config
    buttons         []button     
    bigBox          string
    listenersTable  table.Model
	InputBoxes	InputModel
}

type InputModel struct {
	Labels	[]string
	Inputs 	[]textinput.Model
	Focus 	int
	Err 	error
}

type button struct {
    text    string
    state   sessionState
}

func NewModel() MainModel {
	m := MainModel {
		state:      mainState,
		focus:      0, 
        buttons:    mainButtons,
        bigBox:     GetRandomBanner(),
		config: Config{apiIp: "127.0.0.1", apiPort: "8000", apiVer: "v1"},
	}
    m.listenersTable = m.getDemoTableViewComponent()
	return m
}

func (m InputModel) NewInputModel(labels, placeholders []string) InputModel {
	var inModels []textinput.Model = make([]textinput.Model, len(labels))
	for i := range labels {
		inModels[i] = textinput.New()
		inModels[i].Placeholder = placeholders[i]
		inModels[i].CharLimit = 25
		inModels[i].Width = 25
		inModels[i].Prompt = "# "
	}

	r := InputModel {
		Labels: labels,
		Inputs:	inModels,
		Focus: 0,
		Err: nil,
	}
	return r
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
    prevState := m.state

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
            switch m.state {
            case mainState:
                return m, tea.Quit
            default:
                m.state     = mainState
                m.buttons   = mainButtons
            }
		case "tab":
			m.focus = m.NextFocus()
        case "enter":
			m.HandleButton(&m)
			if m.focus >= len(m.buttons){
				m.focus = m.NextFocus()
			} else {
				m.state, m.buttons, m.focus = m.NextState()
			}

            
		}
	}

	if m.state == listenerEditState || m.state == configEditState || m.state == listenerNewState {
		for i := range m.InputBoxes.Inputs {
			m.InputBoxes.Inputs[i], cmd = m.InputBoxes.Inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

    if m.state != prevState {
        m.listenersTable.Blur()
		m.listenersTable.SetCursor(0)
		m.bigBox = GetRandomBanner()

        switch m.state {
		case configEditState:
			placeholders := []string {m.config.apiIp, m.config.apiPort,
				m.config.apiVer}
				m.InputBoxes = m.InputBoxes.NewInputModel([]string{"Ip",
				"Port", "Version"}, placeholders)
        case listenersState:
            m.listenersTable.Focus()
            m.listenersTable.SetCursor(0)
        case listenerEditState:
			placeholders := []string {
            m.listenersTable.SelectedRow()[1],
            m.listenersTable.SelectedRow()[2],
            m.listenersTable.SelectedRow()[3]}
			m.InputBoxes = m.InputBoxes.NewInputModel([]string{"Name",
            "Ip", "Port"}, placeholders)
		case listenerNewState:
			placeholders := []string {"Bob", "127.0.0.1", "80"}
			m.InputBoxes = m.InputBoxes.NewInputModel([]string{"Name",
            "Ip", "Port"}, placeholders)
        }
    }
    m.listenersTable, cmd = m.listenersTable.Update(msg)
    cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) NextFocus() int {
	f := m.focus + 1
	if m.state == listenerEditState || m.state == configEditState || m.state == listenerNewState {
		for i := range m.InputBoxes.Inputs {
			m.InputBoxes.Inputs[i].Blur()
		}
		if f == len(m.buttons) + len(m.InputBoxes.Inputs){
			return 0
		} else if f >= len(m.buttons) {
			m.InputBoxes.Inputs[f - len(m.buttons)].Focus()
		}
	} else {
		if f >= len(m.buttons) {
			return 0
		}
	}
	return f	
}

func (m MainModel) HandleButton(pm *MainModel) {
	url := "http://" + m.config.apiIp + ":" + m.config.apiPort + "/" + m.config.apiVer + "/"
	switch m.state {
	case configEditState:
		if m.focus == 0 {
			pm.config.apiIp = m.InputBoxes.Inputs[0].Value()
			pm.config.apiPort = m.InputBoxes.Inputs[1].Value()
			pm.config.apiVer = m.InputBoxes.Inputs[2].Value()
		}
	case listenerNewState:
		if m.focus == 0 {
			body := []byte(fmt.Sprintf(`{"name": "%s", "ip": "%s", "port": "%s"}`,
			m.InputBoxes.Inputs[0].Value(),
			m.InputBoxes.Inputs[1].Value(),
			m.InputBoxes.Inputs[2].Value(),))
			r, err := http.NewRequest("POST", url + "l/new", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			client := &http.Client{}
			_, err = client.Do(r)
			if err != nil {
				panic(err)
			}
			defer r.Body.Close()

			//fmt.Println(resp)
		}
	default:

	}
}

func (m MainModel) NextState() (sessionState, []button, int) {
    f := 0
    s := m.buttons[m.focus].state
    var b []button
    switch s {
	case configEditState:
		b = configEditButtons
    case listenersState:
        b = listenersButtons
    case listenerEditState:
        b = listenerEditButtons
	case listenerNewState:
        b = listenerEditButtons
    default:
        b = mainButtons
    }
    if s == m.state {
        f = m.focus
    }
    return s, b, f
} 

func (m MainModel) View() string {
    b := m.getButtonViewComponent()
	switch m.state {
	case configEditState:
		return lipgloss.JoinVertical(lipgloss.Top,
            m.getHeaderViewComponent(headerText),
			m.getInputViewComponent(),
            m.getFooterViewComponent(), b)
	case mainState:
        return lipgloss.JoinVertical(lipgloss.Top,
            m.getHeaderViewComponent(headerText),
            bigBoxStyle.Render(m.bigBox),
            m.getFooterViewComponent(), b)
    case listenersState:
        return lipgloss.JoinVertical(lipgloss.Top,
            m.getHeaderViewComponent(headerText),
            baseTableStyle.Render(m.listenersTable.View()),
            m.getFooterViewComponent(), b)
    case listenerNewState:
        return lipgloss.JoinVertical(lipgloss.Top,
            m.getHeaderViewComponent(headerText),
			m.getInputViewComponent(),
            m.getFooterViewComponent(), b)
    case listenerEditState:
        return lipgloss.JoinVertical(lipgloss.Top,
            m.getHeaderViewComponent(headerText),
			m.getInputViewComponent(),
            m.getFooterViewComponent(), b)
    case listenerInfoState:
        return fmt.Sprintf("TODO: View info of listener %s.",
            m.listenersTable.SelectedRow()[1])
	}
	return ""
}

func (m MainModel) getInputViewComponent() string {
    var iview string
    var temp string
    for i, t := range m.InputBoxes.Labels {
        if i+len(m.buttons) == m.focus {
            m.InputBoxes.Inputs[i].Prompt = "> "
        } else {
            m.InputBoxes.Inputs[i].Prompt = "# "
        }
        temp = lipgloss.JoinVertical(lipgloss.Top,
        inputLabelStyle.Render(t),
        inputBoxStyle.Render(m.InputBoxes.Inputs[i].View()))
		iview = lipgloss.JoinVertical(lipgloss.Top, iview, temp)
    }
	return inputBigBoxStyle.Render(iview)
}

func (m MainModel) getButtonViewComponent() string {
    var bview string
    for i, b := range m.buttons {
        if i == m.focus {
            bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
            focusButtonStyle.Render(b.text))
        } else {
            bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
            buttonStyle.Render(b.text))
        }
    }
    return bview
}

func (m MainModel) getHeaderViewComponent(s string) string {
    lline := strings.Repeat(borderChar, 3)
	text := headerStyle.Render(s)
	rline := strings.Repeat(borderChar, maxWidth - (len(lline) + len(headerText) + 4))
	return lipgloss.JoinHorizontal(lipgloss.Center, lline, text, rline)
}

func (m MainModel) getFooterViewComponent() string {
	line := strings.Repeat(borderChar, maxWidth)
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}

func (m MainModel) getDemoTableViewComponent() table.Model {
    numCol := 5
    tWidth := (maxWidth / numCol) - ((2 * numCol) / numCol)
    tHeight := maxHeight - 4

	columns := []table.Column{
		{Title: "Id",       Width: tWidth},
		{Title: "Name",     Width: tWidth},
		{Title: "Ip",       Width: tWidth},
		{Title: "Port",     Width: tWidth},
		{Title: "Status",   Width: tWidth},
	}

	rows := []table.Row{
		{"1", "Tokyo",      "127.0.0.1",    "80", "1"},
		{"2", "Delhi",      "127.0.0.2",    "81", "1"},
		{"3", "Shanghai",   "127.0.0.3",    "82", "1"},
		{"4", "Dhaka",      "127.0.0.4",    "84", "1"},
		{"5", "Sao",        "127.0.0.5",    "85", "1"},
		{"6", "Colombo",    "127.0.0.6",    "86", "1"},
		{"7", "Cairo",      "127.0.0.7",    "87", "1"},
		{"8", "Beijing",    "127.0.0.8",    "88", "1"},
		{"9", "Mumbai",     "127.0.0.9",    "89", "1"},
		{"10", "Osaka",     "127.0.0.10",   "90", "1"},
	}

	t := table.New(
		table.WithColumns(columns),
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
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
