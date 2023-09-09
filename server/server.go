package server

import (
	"net"

	"github.com/TorkhanNetwork/Networking/data_encryption"
	"github.com/TorkhanNetwork/Networking/events_system"
	"github.com/kataras/golog"
)

const (
	SERVER_VERSION = "1.0.0"
	COMMAND_PREFIX = "!"
)

type Server struct {
	Name           string
	socketListener SocketListener
	port           int
	Clients        []ServerWorker
	EventsManager  events_system.EventsManager
	keysGenerator  data_encryption.KeysGenerator
	count          int
	idle           chan int
	terminate      chan int
}

func NewServer(name string, port int) Server {
	kg := data_encryption.NewGenerator()
	ok := ReadAsyncKeys(&kg)
	if !ok {
		kg.GenerateKeys(false, true)
		WriteAsyncKeys(kg)
	}
	return Server{
		Name:           name,
		socketListener: SocketListener{},
		port:           port,
		Clients:        make([]ServerWorker, 0),
		EventsManager:  events_system.NewEventsManager(),
		keysGenerator:  kg,
	}
}

func (server *Server) StartSocketListener() {
	if !server.socketListener.Alive {
		server.socketListener.Start(server)
	}
}

func (server *Server) HandleSocketConnection(connection *net.TCPConn) {
	event := NewClientSocketConnectEvent(connection)
	var e interface{} = event
	server.EventsManager.CallEvent((*events_system.Event)(&e))
	if !event.Cancel {
		server.count++
		serverWorker := NewServerWorker(server.count, server, connection)
		golog.Infof("%s - New socket accepted from %s -> #%s", server.Name, connection.RemoteAddr().String(), serverWorker.Id)
		event := NewServerWorkerBoundEvent(&serverWorker)
		e = event
		if !event.Cancel {
			server.Clients = append(server.Clients, serverWorker)
			serverWorker.StartWorker()
		}
	}
}

func (server *Server) RemoveServerWorker(serverWorker ServerWorker) {
	for i, s := range server.Clients {
		if s.Id == serverWorker.Id {
			server.Clients[i] = server.Clients[len(server.Clients)-1]
			server.Clients = server.Clients[:len(server.Clients)-1]
		}
	}
	if len(server.Clients) == 0 {
		server.idle <- 1
	}
}

func (server *Server) JoinWorkers() {
	if len(server.Clients) == 0 {
		return
	}
	<-server.idle
}

func (server *Server) Join() {
	<-server.terminate
}
