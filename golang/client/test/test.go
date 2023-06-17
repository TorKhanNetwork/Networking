package main

import (
	"strings"

	"github.com/TorkhanNetwork/Networking/golang/client"
	"github.com/kataras/golog"
)

type SimpleListener struct {
	client client.Client
}

func (l SimpleListener) OnCommand(e client.CommandReceivedEvent) {
	golog.Debug(l.client.Name + " command : " + e.Command + " " + strings.Join(e.Args, " "))
}

func (l SimpleListener) OnEncData(e client.EncryptedDataReceivedEvent) {
	golog.Debug(l.client.Name + " enc data : " + e.DecryptedData)
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
