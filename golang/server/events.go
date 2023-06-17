package server

import (
	"net"
	"strings"
)

type ClientDisconnectEvent struct {
	ServerWorker ServerWorker
}

func NewClientDisconnectEvent(serverWorker *ServerWorker) ClientDisconnectEvent {
	return ClientDisconnectEvent{ServerWorker: *serverWorker}
}

type ClientSocketConnectEvent struct {
	Connection *net.TCPConn
	Cancel     bool
}

func NewClientSocketConnectEvent(connection *net.TCPConn) ClientSocketConnectEvent {
	return ClientSocketConnectEvent{Connection: connection, Cancel: false}
}

type CommandReceivedEvent struct {
	ServerWorker ServerWorker
	Command      string
	Args         []string
}

func NewCommandReceivedEvent(serverWorker *ServerWorker, data string, prefix string, splitter string) CommandReceivedEvent {
	splitted := strings.Split(data[len(prefix):], splitter)
	args := make([]string, 0)
	if len(splitted) >= 2 {
		args = strings.Split(data[len(prefix)+len(splitter)+len(splitted[0]):], splitter)
	}
	return CommandReceivedEvent{
		ServerWorker: *serverWorker,
		Command:      splitted[0],
		Args:         args,
	}
}

type ConnectionProtocolSuccessEvent struct {
	ServerWorker ServerWorker
}

func NewConnectionProtocolSuccessEvent(serverWorker *ServerWorker) ConnectionProtocolSuccessEvent {
	return ConnectionProtocolSuccessEvent{ServerWorker: *serverWorker}
}

type EncryptedDataReceivedEvent struct {
	ServerWorker        ServerWorker
	Data, DecryptedData string
}

func NewEncryptedDataReceivedEvent(serverWorker *ServerWorker, data, decryptedData string) EncryptedDataReceivedEvent {
	return EncryptedDataReceivedEvent{ServerWorker: *serverWorker, Data: data, DecryptedData: decryptedData}
}

type RawDataReceivedEvent struct {
	ServerWorker ServerWorker
	Data         string
}

func NewRawDataReceivedEvent(serverWorker *ServerWorker, data string) RawDataReceivedEvent {
	return RawDataReceivedEvent{ServerWorker: *serverWorker, Data: data}
}

type ServerWorkerBoundEvent struct {
	ServerWorker ServerWorker
	Cancel       bool
}

func NewServerWorkerBoundEvent(serverWorker *ServerWorker) ServerWorkerBoundEvent {
	return ServerWorkerBoundEvent{
		ServerWorker: *serverWorker,
		Cancel:       false,
	}
}
