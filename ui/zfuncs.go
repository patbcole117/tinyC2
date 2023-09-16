package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/patbcole117/tinyC2/node"
)

type button struct {
	do   func() tea.Msg
	text string
}

type inputCancelMsg string
func cancelConfig() tea.Msg {
	return inputCancelMsg("Config")
}

type inputSaveMsg string
func saveConfig() tea.Msg {
	return inputSaveMsg("Config")
}

type newInfoMsg string
func setInfoMsg(msg string) tea.Cmd {
	return func() tea.Msg {
		return newInfoMsg(msg)
	}
}

type syncNodesMsg []node.Node
func SyncNodes(c apiConfig) tea.Cmd {
	return func() tea.Msg {
		var msg string
		var nodes []node.Node
		url := "http://" + c.apiIp + ":" + c.apiPort + "/" + c.apiVer + "/l"

		resp, err := http.Get(url)
		if err != nil {
			msg = errMsg("SyncNodes:http.Get", resp.Status)
			return newInfoMsg(msg)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			msg = errMsg("SyncNodes:io.ReadAll", resp.Status)
			return newInfoMsg(msg)
		}

		err = json.Unmarshal(body, &nodes)
		if err != nil {
			msg = errMsg("SyncNodes:json.Unmarshal", resp.Status)
			return newInfoMsg(msg)
		}
		return syncNodesMsg(nodes)
	}
}

type trigDeleteNodeMsg string
func trigDeleteNode() tea.Msg { return trigDeleteNodeMsg("DeleteNode") }
func DeleteNode(id string, c apiConfig) tea.Cmd {
	return func() tea.Msg {
		var msg string
		url := "http://" + c.apiIp + ":" + c.apiPort + "/" + c.apiVer + "/l/delete"

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(id)))
		if err != nil {
			msg = errMsg("DeleteNode:http.NewRequest", err.Error())
			return newInfoMsg(msg)
		}
		req.Close = true
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			msg = errMsg("DeleteNode:client.Do", err.Error())
			return newInfoMsg(msg)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = errMsg("DeleteNode:io.ReadAll", err.Error())
				return newInfoMsg(msg)
			}
			msg = sucMsg("DELETE", string(body))
		} else {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = errMsg("DeleteNode:io.ReadAll", string(b)+err.Error())
				return newInfoMsg(msg)
			}
		}
		return newInfoMsg(msg)
	}
}

type trigNewNodeMsg string
func trigNewNode() tea.Msg { return trigNewNodeMsg("NewNode") }
func NewNode(name, ip, port string, c apiConfig) tea.Cmd {
	return func() tea.Msg {
		var msg string
		if name == "" || ip == "" || port == "" {
			msg = errMsg("NewNode:InputValidation","EMPTY FIELD")
			return newInfoMsg(msg)
		}
		
		url := "http://" + c.apiIp + ":" + c.apiPort + "/" + c.apiVer + "/l/new"
		n := node.NewNode()
		n.Name = name
		n.Ip = ip
		p, err := strconv.Atoi(port)
		n.Port = p
		if err != nil {
			msg = errMsg("NewNode:strconv.Atoi", err.Error())
			return newInfoMsg(msg)
		}

		body, err := json.Marshal(n)
		if err != nil {
			msg = errMsg("NewNode:json.Marshal", err.Error())
			return newInfoMsg(msg)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			msg = errMsg("NewNode:http.NewRequest", err.Error())
			return newInfoMsg(msg)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Close = true
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			msg = errMsg("NewNode:client.Do", err.Error())
			return newInfoMsg(msg)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = errMsg("NewNode:io.ReadAll", err.Error())
				return newInfoMsg(msg)
			}
			msg = sucMsg("NEW NODE", string(body))
		} else {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = errMsg("NewNode:io.ReadAll", string(b)+err.Error())
				return newInfoMsg(msg)
			}
		}
		return newInfoMsg(msg)
	}
}

type trigToggleNodeMsg string
func trigStartNode() tea.Msg { return trigToggleNodeMsg("START")}
func trigStopNode() tea.Msg { return trigToggleNodeMsg("STOP")}
type trigUpdateNodeMsg string
func trigUpdateNode() tea.Msg { return trigUpdateNodeMsg("UpdateNode") }
func UpdateNode(n node.Node, c apiConfig) tea.Cmd {
	return func() tea.Msg {
		var msg string
		url := "http://" + c.apiIp + ":" + c.apiPort + "/" + c.apiVer + "/l/update"

		body, err := json.Marshal(n)
		if err != nil {
			msg = errMsg("UpdateNode:json.Marshal", err.Error())
			return newInfoMsg(msg)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			msg = errMsg("UpdateNode:http.NewRequest", err.Error())
			return newInfoMsg(msg)
		}
		req.Close = true
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			msg = errMsg("UpdateNode:client.Do", err.Error())
			return newInfoMsg(msg)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = errMsg("UpdateNode:io.ReadAll", err.Error())
				return newInfoMsg(msg)
			}
			msg = sucMsg("UPDATE", string(body))
		} else {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = errMsg("UpdateNode:io.ReadAll", string(b)+err.Error())
				return newInfoMsg(msg)
			}
		}
		return newInfoMsg(msg)
	}
}

type setStateMsg string
func TODOButton() tea.Msg {
	return setStateMsg("TODO")
}
func toRootState() tea.Msg {
	return setStateMsg("Root")
}
func toConfigState() tea.Msg {
	return setStateMsg("Config")
}
func toNodesState() tea.Msg {
	return setStateMsg("Nodes")
}
func toNodesEditState() tea.Msg {
	return setStateMsg("NodesEdit")
}
func toNodesNewState() tea.Msg {
	return setStateMsg("NodesNew")
}
func toNodesInfoState() tea.Msg {
	return setStateMsg("NodesInfo")
}

func NextFocus(cur, max int) int {
	f := cur + 1
	if f >= max {
		return 0
	}
	return f
}

func GetButtonViewComponent(buttons []button, focus int) string {
	var bview string
	for i, b := range buttons {
		if i == focus {
			bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
				buttonFocusStyle.Render(b.text))
		} else {
			bview = lipgloss.JoinHorizontal(lipgloss.Top, bview,
				buttonStyle.Render(b.text))
		}
	}
	return bview
}
func GetFooterViewComponent() string {
	line := strings.Repeat(borderChar, maxWidth)
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}
func GetHeaderViewComponent() string {
	lline := strings.Repeat(borderChar, 3)
	text := headerStyle.Render(headerText)
	rline := strings.Repeat(borderChar, maxWidth-(len(lline)+len(headerText)+4))
	return lipgloss.JoinHorizontal(lipgloss.Center, lline, text, rline)
}

func errMsg(strErr, msg string) string {
	return fmt.Sprintf(`{"ERROR": "%s", "Msg": "%s"}`, strErr, msg)
}

func sucMsg(strSuc, msg string) string {
	return fmt.Sprintf(`{"SUCCESS": "%s", "Msg": "%s"}`, strSuc, msg)
}

func infMsg(strInf, msg string) string {
	return fmt.Sprintf(`{"INFO": "%s", "Msg": "%s"}`, strInf, msg)
}