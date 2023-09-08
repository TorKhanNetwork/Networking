package client

import (
	"github.com/TorkhanNetwork/Networking/events_system"
	"github.com/kataras/golog"
)

const CLIENT_VERSION = "1.0"

type Client struct {
	Name             string
	MacAddress       string
	SocketWorkerList []SocketWorker
	EventsManager    events_system.EventsManager
	count            int
	exit             chan int
}

func NewClient(name string) Client {
	macAddress, err := GetMacAddress()
	if err != nil {
		golog.Fatal("Unable to find MAC Address")
	}
	return Client{
		Name:             name,
		MacAddress:       macAddress,
		SocketWorkerList: make([]SocketWorker, 0),
		EventsManager:    events_system.NewEventsManager(),
		exit:             make(chan int),
		count:            0,
	}
}

func (client *Client) AddSocketWorker(ip string, port int) SocketWorker {
	client.count++
	socketWorker := NewSocketWorker(client.count, client, ip, port)
	client.SocketWorkerList = append(client.SocketWorkerList, socketWorker)
	return socketWorker
}

func (client *Client) RemoveSocketWorker(socketWorker SocketWorker) {
	for i, s := range client.SocketWorkerList {
		if s.Id == socketWorker.Id {
			client.SocketWorkerList[i] = client.SocketWorkerList[len(client.SocketWorkerList)-1]
			client.SocketWorkerList = client.SocketWorkerList[:len(client.SocketWorkerList)-1]
		}
	}
	if len(client.SocketWorkerList) == 0 {
		client.exit <- 1
	}
}

func (client *Client) JoinWorkers() {
	if len(client.SocketWorkerList) == 0 {
		return
	}
	<-client.exit
}
