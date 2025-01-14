package main

import "github.com/sagernet/sing-box/experimental/libbox"

var commandServer *libbox.CommandServer

type CommandServerHandler struct{}

func (csh *CommandServerHandler) ServiceReload() error {
	propagateStatus(Starting)
	if commandServer != nil {
		commandServer.SetService(nil)
		commandServer = nil
	}
	if box != nil {
		box.Close()
		box = nil
	}
	return startService()
}

func (csh *CommandServerHandler) GetSystemProxyStatus() *libbox.SystemProxyStatus {
	return &libbox.SystemProxyStatus{Available: true, Enabled: false}
}

func (csh *CommandServerHandler) SetSystemProxyEnabled(isEnabled bool) error {
	return nil
}

func startCommandServer() error {
	commandServer = libbox.NewCommandServer(&CommandServerHandler{}, 300)
	return commandServer.Start()
}
