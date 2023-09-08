package main

import (
	"github.com/TorkhanNetwork/Networking/client"
	"github.com/TorkhanNetwork/Networking/response_system"
	"github.com/kataras/golog"
)

type SimpleListener struct {
	client client.Client
}

func (l SimpleListener) OnAuth(e client.ConnectionProtocolSuccessEvent) {
	r := e.SocketWorker.SendCommand("platform")
	err := r.WaitReply(true, func(worker response_system.IResponseWorker, s string) {
		golog.Debugf("askip sa plateforme c'est %s", s)
	})
	if err != nil {
		golog.Warnf("%s", err)
	}
	golog.Debug("blou")
}

func (l SimpleListener) OnCommand(e client.CommandReceivedEvent) {
	if e.Command == "platform" {
		e.SocketWorker.SendData("Golang bebew", e.MsgUUID, true, true)
	}
}

func (l SimpleListener) OnClose(e client.ServerDisconnectEvent) {
	golog.Debug(l.client.Name + " disconnection from " + e.SocketWorker.GetName())
}

func main() {
	golog.SetLevel("debug")
	client := client.NewClient("Petit Client de test")
	client.EventsManager.RegisterListener(SimpleListener{client: client})
	worker := client.AddSocketWorker("192.168.1.22", 40000)
	worker.StartWorker()
	client.JoinWorkers()
}
