package main

import (
	"github.com/TorkhanNetwork/Networking/golang/client"
	"github.com/TorkhanNetwork/Networking/golang/response_system"
	"github.com/google/uuid"
	"github.com/kataras/golog"
)

type SimpleListener struct {
	client client.Client
}

func (l SimpleListener) OnAuth(e client.ConnectionProtocolSuccessEvent) {
	r := e.SocketWorker.SendData("?platform", uuid.Nil, true)
	err := r.WaitReply(true, func(worker response_system.IResponseWorker, s string) {
		golog.Debugf("askip sa plateforme c'est %s", s)
	})
	if err != nil {
		golog.Errorf("%s", err)
	}
	golog.Debug("blou")
}

func (l SimpleListener) OnEncData(e client.EncryptedDataReceivedEvent) {
	golog.Debug(l.client.Name + " enc data : " + e.DecryptedData)
	if e.DecryptedData == "?platform" {
		e.SocketWorker.SendData("Golang bebew", e.MsgUUID, true)
	}
}

func (l SimpleListener) OnRaw(e client.RawDataReceivedEvent) {
	golog.Debug(l.client.Name + " raw data : " + e.Data)
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
