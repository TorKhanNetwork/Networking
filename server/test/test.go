package main

import (
	"strings"

	"github.com/TorKhanNetwork/Networking/response_system"
	"github.com/TorKhanNetwork/Networking/server"
	"github.com/kataras/golog"
)

type SimpleListener struct {
}

func (l SimpleListener) OnCommand(e server.CommandReceivedEvent) {
	golog.Debugf("args: %s", strings.Join(e.Args, " | "))
	if e.Command == "platform" {
		e.ServerWorker.SendData("Golang bebew", e.MsgUUID, true)
		r := e.ServerWorker.SendCommand("platform")
		r.WaitReply(false, func(worker response_system.IResponseWorker, s string) {
			golog.Debugf("askip sa plateforme c'est %s", s)
		})
	}
}

func main() {
	golog.SetLevel("debug")
	server := server.NewServer("Petit Serveur de test", 40000)
	server.EventsManager.RegisterListener(SimpleListener{})
	server.StartSocketListener()
	server.Join()
}
