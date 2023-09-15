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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/patbcole117/tinyC2/node"
)


type button struct {
	do   func() tea.Msg
	text string
}


type trigNewListenerMsg string
func trigNewListener() tea.Msg {return trigNewListenerMsg("NewListener")}
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
			return newInfoMsg(msg)
		}

		body, err := json.Marshal(n)
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "json.Marshal", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "http.NewRequest", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Close = true
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "client.Do", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = fmt.Sprintf(`{"ERROR": "io.ReadAll", "Msg": "%s"}`, err)
				return newInfoMsg(msg)
			}
	
			jsonStr := string(body)
			msg = fmt.Sprintf(`{"SUCCESS": "%s"}`, jsonStr)
	
		} else {
			//The status is not Created. print the error.
			msg = fmt.Sprintf(`{"ERROR": "resp.StatusCode", "Msg": "%s"}`, resp.Status)
			return newInfoMsg(msg)
		}
		return newInfoMsg(msg)
	}
}

type trigDeleteListenerMsg string
func trigDeleteListener() tea.Msg {return trigDeleteListenerMsg("DeleteListener")}
func DeleteListener(id string, c apiConfig) tea.Cmd {
	return func() tea.Msg {
		var msg string
		url := "http://" + c.apiIp + ":" + c.apiPort + "/" + c.apiVer + "/l/delete"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(id)))
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "http.NewRequest", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}
		req.Close = true
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "client.Do", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = fmt.Sprintf(`{"ERROR": "io.ReadAll", "Msg": "%s"}`, err)
				return newInfoMsg(msg)
			}
	
			jsonStr := string(body)
			msg = fmt.Sprintf(`{"SUCCESS": "%s"}`, jsonStr)
	
		} else {
			//The status is not Created. print the error.
			msg = fmt.Sprintf(`{"ERROR": "resp.StatusCode", "Msg": "%s"}`, resp.Status)
			return newInfoMsg(msg)
		}
		return newInfoMsg(msg)
	}
}

type trigUpdateNodeMsg string
func trigUpdateNode() tea.Msg {return trigUpdateNodeMsg("UpdateNode")}
func UpdateNode(id, name, ip, port string, c apiConfig) tea.Cmd {
	return func() tea.Msg {
		var msg string
		url := "http://" + c.apiIp + ":" + c.apiPort + "/" + c.apiVer + "/l/update"
		n := node.NewNode()
		n.Id = id
		n.Name = name
		n.Ip = ip
		p, err := strconv.Atoi(port)
		n.Port = p
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "strconv.Atoi", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}

		body, err := json.Marshal(n)
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "json.Marshal", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "http.NewRequest", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}
		req.Close = true
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			msg = fmt.Sprintf(`{"ERROR": "client.Do", "Msg": "%s"}`, err)
			return newInfoMsg(msg)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = fmt.Sprintf(`{"ERROR": "io.ReadAll", "Msg": "%s"}`, err)
				return newInfoMsg(msg)
			}
	
			jsonStr := string(body)
			msg = fmt.Sprintf(`{"SUCCESS": "%s"}`, jsonStr)
	
		} else {
			b, _ := io.ReadAll(resp.Body)
			//The status is not Created. print the error.
			msg = fmt.Sprintf(`{"ERROR": "resp.StatusCode", "Msg": "%s"}`, string(b))
			return newInfoMsg(msg)
		}
		return newInfoMsg(msg)
	}
}

type syncNodesMsg []node.Node
func SyncNodes(c apiConfig)  tea.Cmd {
	return func() tea.Msg {
		var msg string
		var nodes []node.Node
		url := "http://" + c.apiIp + ":" + c.apiPort + "/" + c.apiVer + "/l"
	
		resp, err := http.Get(url)
			if err != nil {
				msg = fmt.Sprintf(`{"ERROR": " http.Get", "Msg": "%s"}`, resp.Status)
			return newInfoMsg(msg)
			}
			defer resp.Body.Close()
	
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = fmt.Sprintf(`{"ERROR": "io.ReadAll", "Msg": "%s"}`, resp.Status)
			return newInfoMsg(msg)
			}
	
			err = json.Unmarshal(body, &nodes)
			if err != nil {
				msg = fmt.Sprintf(`{"ERROR": "json.Unmarshal", "Msg": "%s"}`, resp.Status)
			return newInfoMsg(msg)
			}

			

			return syncNodesMsg(nodes)
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


type newInfoMsg string
func setInfoMsg(msg string) tea.Cmd {
	return func() tea.Msg {
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
func toListenersState() tea.Msg {
	return setStateMsg("Listeners")
}
func toListenersEditState() tea.Msg {
	return setStateMsg("ListenersEdit")
}
func toListenersNewState() tea.Msg {
	return setStateMsg("ListenersNew")
}
func toListenersInfoState() tea.Msg {
	return setStateMsg("ListenersInfo")
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
