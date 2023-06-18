package main

import (
	"github.com/TorkhanNetwork/Networking/golang/response_system"
	"github.com/TorkhanNetwork/Networking/golang/server"
	"github.com/google/uuid"
	"github.com/kataras/golog"
)

type SimpleListener struct {
}

func (l SimpleListener) OnEncData(e server.EncryptedDataReceivedEvent) {
	if e.DecryptedData == "?platform" {
		e.ServerWorker.SendData("Golang bebew", e.MsgUUID, true)
		r := e.ServerWorker.SendData("?platform", uuid.Nil, true)
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
