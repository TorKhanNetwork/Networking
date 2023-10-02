package main

import (
	"strings"

	"github.com/TorKhanNetwork/Networking/client"
	"github.com/TorKhanNetwork/Networking/response_system"
	"github.com/kataras/golog"
)

type SimpleListener struct {
	client client.Client
}

func (l SimpleListener) OnAuth(e client.ConnectionProtocolSuccessEvent) {
	r := e.SocketWorker.SendCommand("platform", "test", "\"blou\"", `"blou2"`)
	err := r.WaitReply(true, func(worker response_system.IResponseWorker, s string) {
		golog.Debugf("askip sa plateforme c'est %s", s)
	})
	if err != nil {
		golog.Warnf("%s", err)
	}
}

func (l SimpleListener) OnCommand(e client.CommandReceivedEvent) {
	golog.Debugf("args: %s", strings.Join(e.Args, " | "))
	if e.Command == "platform" {
		e.SocketWorker.SendData("Golang bebew", e.MsgUUID, true)
	}
}

func (l SimpleListener) OnClose(e client.ServerSocketDisconnectEvent) {
	golog.Debug(l.client.Name + " illegal disconnection from " + e.SocketWorker.GetName())
}

func main() {
	golog.SetLevel("debug")
	client := client.NewClient("Petit Client de test")
	client.EventsManager.RegisterListener(SimpleListener{client: client})
	worker := client.AddSocketWorker("127.0.0.1", 40000)
	worker.StartWorker()
	client.JoinWorkers()
}
