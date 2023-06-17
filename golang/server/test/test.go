package main

import (
	"github.com/TorkhanNetwork/Networking/golang/server"
	"github.com/kataras/golog"
)

func main() {
	golog.SetLevel("debug")
	server := server.NewServer("Petit Serveur de test", 40000)
	server.StartSocketListener()
	server.Join()
}
