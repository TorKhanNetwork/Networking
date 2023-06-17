package client

import (
	"strings"
)

type CommandReceivedEvent struct {
	SocketWorker SocketWorker
	Command      string
	Args         []string
}

func NewCommandReceivedEvent(socketWorker *SocketWorker, data string, prefix string, splitter string) CommandReceivedEvent {
	splitted := strings.Split(data[len(prefix):], splitter)
	args := make([]string, 0)
	if len(splitted) >= 2 {
		args = strings.Split(data[len(prefix)+len(splitter)+len(splitted[0]):], splitter)
	}
	return CommandReceivedEvent{
		SocketWorker: *socketWorker,
		Command:      splitted[0],
		Args:         args,
	}
}

type ConnectionProtocolSuccessEvent struct {
	SocketWorker SocketWorker
}

func NewConnectionProtocolSuccessEvent(socketWorker *SocketWorker) ConnectionProtocolSuccessEvent {
	return ConnectionProtocolSuccessEvent{SocketWorker: *socketWorker}
}

type EncryptedDataReceivedEvent struct {
	SocketWorker        SocketWorker
	Data, DecryptedData string
}

func NewEncryptedDataReceivedEvent(socketWorker *SocketWorker, data, decryptedData string) EncryptedDataReceivedEvent {
	return EncryptedDataReceivedEvent{SocketWorker: *socketWorker, Data: data, DecryptedData: decryptedData}
}

type RawDataReceivedEvent struct {
	SocketWorker SocketWorker
	Data         string
}

func NewRawDataReceivedEvent(socketWorker *SocketWorker, data string) RawDataReceivedEvent {
	return RawDataReceivedEvent{SocketWorker: *socketWorker, Data: data}
}

type ServerDisconnectEvent struct {
	SocketWorker SocketWorker
}

func NewServerDisconnectEvent(socketWorker *SocketWorker) ServerDisconnectEvent {
	return ServerDisconnectEvent{SocketWorker: *socketWorker}
}
