package events

import (
	"net"
	"strings"
)

type CommandReceivedEvent struct {
	SocketWorker interface{}
	Command      string
	Args         []string
}

func NewCommandReceivedEvent(socketWorker interface{}, data string, prefix string, splitter string) CommandReceivedEvent {
	splitted := strings.Split(data[len(prefix):], splitter)
	args := make([]string, 0)
	if len(splitted) >= 2 {
		args = strings.Split(data[len(prefix)+len(splitter)+len(splitted[0]):], splitter)
	}
	return CommandReceivedEvent{
		SocketWorker: socketWorker,
		Command:      splitted[0],
		Args:         args,
	}
}

type ConnectionProtocolSuccessEvent struct {
	SocketWorker interface{}
}

func NewConnectionProtocolSuccessEvent(socketWorker interface{}) ConnectionProtocolSuccessEvent {
	return ConnectionProtocolSuccessEvent{SocketWorker: socketWorker}
}

type EncryptedDataReceivedEvent struct {
	SocketWorker        interface{}
	Data, DecryptedData string
}

func NewEncryptedDataReceivedEvent(socketWorker interface{}, data, decryptedData string) EncryptedDataReceivedEvent {
	return EncryptedDataReceivedEvent{SocketWorker: socketWorker, Data: data, DecryptedData: decryptedData}
}

type RawDataReceivedEvent struct {
	SocketWorker interface{}
	Data         string
}

func NewRawDataReceivedEvent(socketWorker interface{}, data string) RawDataReceivedEvent {
	return RawDataReceivedEvent{SocketWorker: socketWorker, Data: data}
}

type ServerSocketClosedEvent struct {
	SocketWorker interface{}
	Connection   net.TCPConn
}

func NewServerSocketClosedEvent(socketWorker interface{}, connection net.TCPConn) ServerSocketClosedEvent {
	return ServerSocketClosedEvent{SocketWorker: socketWorker, Connection: connection}
}
