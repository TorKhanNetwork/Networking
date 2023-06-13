package main

import (
	"github.com/TorkhanNetwork/Networking/golang/client"
	"github.com/kataras/golog"
)

func main() {
	golog.SetLevel("debug")
	client := client.NewClient("Petit Client de test")
	worker := client.AddSocketWorker("localhost", 40000)
	worker.StartWorker()
	client.JoinWorkers()
}
