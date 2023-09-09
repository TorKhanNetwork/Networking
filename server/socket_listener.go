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
			golog.Errorf("%s - Unable to listen sockets : %s", server.Name, err)
			return
		}
		socketListener.connection = c
		socketListener.Alive = true
		golog.Infof("%s - Listening for clients connections", server.Name)
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
					server.HandleSocketConnection(socket)
				}
			}
		}
	}()
}
