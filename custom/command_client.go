package main

import (
	"encoding/json"
	"fmt"

	"github.com/hiddify/libcore/bridge"
	"github.com/sagernet/sing-box/experimental/libbox"
)

type CommandClientHandler struct {
	port int64
}

func (cch *CommandClientHandler) Connected() {
	fmt.Println("connected")
}

func (cch *CommandClientHandler) Disconnected(message string) {
	fmt.Printf("disconnected: %s\n", message)
}

func (cch *CommandClientHandler) ClearLog() {
	fmt.Println("clear log")
}

func (cch *CommandClientHandler) WriteLog(message string) {
	fmt.Printf("new log: %s\n", message)
}

func (cch *CommandClientHandler) WriteStatus(message *libbox.StatusMessage) {
	msg, err := json.Marshal(
		map[string]int64{
			"connections-in":  int64(message.ConnectionsIn),
			"connections-out": int64(message.ConnectionsOut),
			"uplink":          message.Uplink,
			"downlink":        message.Downlink,
			"uplink-total":    message.UplinkTotal,
			"downlink-total":  message.DownlinkTotal,
		},
	)
	if err != nil {
		bridge.SendStringToPort(cch.port, fmt.Sprintf("error: %e", err))
	} else {
		bridge.SendStringToPort(cch.port, string(msg))
	}
}

func (cch *CommandClientHandler) WriteGroups(message libbox.OutboundGroupIterator) {
	if message == nil {
		return
	}
	groups := []*OutboundGroup{}
	for message.HasNext() {
		group := message.Next()
		items := group.GetItems()
		groupItems := []*OutboundGroupItem{}
		for items.HasNext() {
			item := items.Next()
			groupItems = append(groupItems,
				&OutboundGroupItem{
					Tag:          item.Tag,
					Type:         item.Type,
					URLTestTime:  item.URLTestTime,
					URLTestDelay: item.URLTestDelay,
				},
			)
		}
		groups = append(groups, &OutboundGroup{Tag: group.Tag, Type: group.Type, Selected: group.Selected, Items: groupItems})
	}
	response, err := json.Marshal(groups)
	if err != nil {
		bridge.SendStringToPort(cch.port, fmt.Sprintf("error: %e", err))
	} else {
		bridge.SendStringToPort(cch.port, string(response))
	}
}

func (cch *CommandClientHandler) InitializeClashMode(modeList libbox.StringIterator, currentMode string) {
}

func (cch *CommandClientHandler) UpdateClashMode(newMode string) {}

type OutboundGroup struct {
	Tag      string               `json:"tag"`
	Type     string               `json:"type"`
	Selected string               `json:"selected"`
	Items    []*OutboundGroupItem `json:"items"`
}

type OutboundGroupItem struct {
	Tag          string `json:"tag"`
	Type         string `json:"type"`
	URLTestTime  int64  `json:"url-test-time"`
	URLTestDelay int32  `json:"url-test-delay"`
}
