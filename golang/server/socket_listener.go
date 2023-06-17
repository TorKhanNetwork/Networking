package server

import (
	"net"

	"github.com/kataras/golog"
)

type SocketListener struct {
	Alive      bool
	connection *net.TCPListener
	exit       chan int
}

func NewSocketListener() SocketListener {
	return SocketListener{Alive: false, exit: make(chan int)}
}

func (socketListener *SocketListener) Start(server *Server) {
	go func() {
		c, err := net.ListenTCP("tcp", &net.TCPAddr{Port: server.port})
		if err != nil {
			golog.Error(server.Name+" - Unable to listen sockets\n", err)
			return
		}
		socketListener.connection = c
		socketListener.Alive = true
		golog.Info(server.Name + " - Listening for clients connections")
		for {
			select {
			case <-socketListener.exit:
				c.Close()
				return
			default:
				socket, err := c.AcceptTCP()
				if err != nil {
					golog.Error("Unable to accept incoming connection")
				} else {
					golog.Debug(server.Name + " - New socket accepted from " + socket.RemoteAddr().String())
					server.HandleSocketConnection(socket)
				}
			}
		}
	}()
}
