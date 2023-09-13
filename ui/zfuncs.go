package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/patbcole117/tinyC2/node"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type button struct {
	do   func() tea.Msg
	text string
}

type input struct {
	label   string
	textBox textinput.Model
}

type dbMsg string
func trigNewListener() tea.Msg {
	return dbMsg("NewListener")
}
func NewListener(name, ip, port string, c apiConfig) tea.Cmd {
	return func() tea.Msg {
		var msg string
		url := "http://" + c.apiIp + ":" + c.apiPort + "/" + c.apiVer + "/l/new"
		n := node.NewNode()
		n.Name = name
		n.Ip = ip
		p, err := strconv.Atoi(port)
		n.Port = p
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "strconv.Atoi", "Msg": "%s"}`, err)
			return dbMsg(msg)
		}

		body, err := json.Marshal(n)
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "json.Marshal", "Msg": "%s"}`, err)
			return dbMsg(msg)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "http.NewRequest", "Msg": "%s"}`, err)
			return dbMsg(msg)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Close = true

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "client.Do", "Msg": "%s"}`, err)
			return dbMsg(msg)
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusCreated {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				msg = fmt.Sprintf(`{"ERROR": "io.ReadAll", "Msg": "%s"}`, err)
				return dbMsg(msg)
			}
	
			jsonStr := string(body)
			msg = fmt.Sprintf(`{"SUCCESS": "%s"}`, jsonStr)
	
		} else {
			//The status is not Created. print the error.
			msg = fmt.Sprintf(`{"ERROR": "resp.StatusCode", "Msg": "%s"}`, res.Status)
			return dbMsg(msg)
		}
		return setInfoMsg(msg)
	}
}

type inputCancelMsg string
func cancelConfig() tea.Msg {
	return inputCancelMsg("Config")
}

type inputSaveMsg string
func saveConfig() tea.Msg {
	return inputSaveMsg("Config")
}

type setInfoMsg string
func changeInfoMsg(msg string) tea.Cmd {
	return func() tea.Msg {
		return setInfoMsg(msg)
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
func toListenersState() tea.Msg {
	return setStateMsg("Listeners")
}
func toNewListenerState() tea.Msg {
	return setStateMsg("NewListener")
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

func GetInputViewComponent(inputs []input) string {
	var iview string
	var temp string
	for x, i := range inputs {
		temp = lipgloss.JoinVertical(lipgloss.Top,
			inputLabelStyle.Render(i.label),
			inputTextBoxStyle.Render(inputs[x].textBox.View()))
		iview = lipgloss.JoinVertical(lipgloss.Top, iview, temp)
	}
	return inputBigBoxStyle.Render(iview)
}

func GetHeaderViewComponent() string {
	lline := strings.Repeat(borderChar, 3)
	text := headerStyle.Render(headerText)
	rline := strings.Repeat(borderChar, maxWidth-(len(lline)+len(headerText)+4))
	return lipgloss.JoinHorizontal(lipgloss.Center, lline, text, rline)
}

func GetFooterViewComponent() string {
	line := strings.Repeat(borderChar, maxWidth)
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}

func GetDemoTable() table.Model {
	numCol := 5
	tWidth := (maxWidth / numCol) - ((2 * numCol) / numCol)
	tHeight := maxHeight - 4

	headers := []string{"Id", "Name", "Ip", "Port", "Status"}
	var cols []table.Column
	for _, h := range headers {
		col := table.Column{Title: h, Width: tWidth}
		cols = append(cols, col)
	}

	rows := []table.Row{
		{"1", "Tokyo", "127.0.0.1", "81", "1"},
		{"2", "Colombo", "127.0.0.2", "82", "1"},
		{"3", "Toronto", "127.0.0.3", "83", "1"},
		{"4", "New York", "127.0.0.4", "84", "1"},
	}

	t := table.New(
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
func GetRandomBanner() string {
	return getBanner(rand.Intn(len(banners)))
}

func getBanner(i int) string {
	return banners[i]
}

var banners = [...]string{
	`   **   **                     ******   **** 
	/**  //            **   **  **////** */// *
   ****** ** *******  //** **  **    // /    /*
  ///**/ /**//**///**  //***  /**          *** 
	/**  /** /**  /**   /**   /**         *//  
	/**  /** /**  /**   **    //**    ** *     
	//** /** ***  /**  **      //****** /******
	 //  // ///   //  //        //////  ////// `,
	`::::::::::: ::::::::::: ::::    ::: :::   :::  ::::::::   ::::::::  
	 :+:         :+:     :+:+:   :+: :+:   :+: :+:    :+: :+:    :+: 
	 +:+         +:+     :+:+:+  +:+  +:+ +:+  +:+              +:+  
	 +#+         +#+     +#+ +:+ +#+   +#++:   +#+            +#+    
	 +#+         +#+     +#+  +#+#+#    +#+    +#+          +#+      
	 #+#         #+#     #+#   #+#+#    #+#    #+#    #+#  #+#       
	 ###     ########### ###    ####    ###     ########  ########## `,
	`######    ####    ##  ##   ##  ##    ####     ####   
	 ##       ##     ### ##   ##  ##   ##  ##   ##  ##  
	 ##       ##     ######   ##  ##   ##           ##  
	 ##       ##     ######    ####    ##          ##   
	 ##       ##     ## ###     ##     ##        ##     
	 ##       ##     ##  ##     ##     ##  ##   ##      
	 ##      ####    ##  ##     ##      ####    ######  `,
	`     >=>                               >=>             
	 >=>    >>                      >=>   >=>  >=>>=>  
   >=>>==>     >==>>==>  >=>   >=> >=>        >>   >=> 
	 >=>   >=>  >=>  >=>  >=> >=>  >=>            >=>  
	 >=>   >=>  >=>  >=>    >==>   >=>           >=>   
	 >=>   >=>  >=>  >=>     >=>    >=>   >=>  >=>     
	  >=>  >=> >==>  >=>    >=>       >===>   >======> `,
	`      mm      db                           .g8"""bgd          
	  MM                                 .dP'      M          
	mmMMmm   7MM   7MMpMMMb.   7M'    MF'dM'          pd*"*b. 
	  MM      MM    MM    MM    VA   ,V  MM          (O)   j8 
	  MM      MM    MM    MM     VA ,V   MM.             ,;j9 
	  MM      MM    MM    MM      VVV     Mb.     ,'  ,-='    
	   Mbmo .JMML..JMML  JMML.    ,V        "bmmmd'  Ammmmmmm 
								 ,V                           
	OOb"                            `,
	`             __  __  
	|_ .  _     /     _) 
	|_ | | ) \/ \__  /__ 
			 /           `,
	`                88                              ,ad8888ba,    ad888888b,  
	     ,d     ""                             d8"'     "8b  d8"     "88  
	     88                                   d8'                    a8P  
       MM88MMM  88  8b,dPPYba,   8b       d8  88                  ,d8P"   
	     88     88  88P'    "8a   8b     d8'  88                a8P"      
	     88     88  88       88    8b   d8'   Y8,             a8P'        
	     88,    88  88       88     8b,d8'     Y8a.    .a8P  d8"          
 	    "Y888  88  88       88      Y88'        "Y8888Y"'   88888888888  
 							   d8'                              
d8'`,
}
