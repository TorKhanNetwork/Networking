package server

import (
	"net"
	"strings"

	"github.com/google/uuid"
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
	MsgUUID      uuid.UUID
	Command      string
	Args         []string
}

func NewCommandReceivedEvent(serverWorker *ServerWorker, data string, prefix string, splitter string, msgUUID uuid.UUID) CommandReceivedEvent {
	splitted := strings.Split(data[len(prefix):], splitter)
	args := make([]string, 0)
	if len(splitted) >= 2 {
		args = strings.Split(data[len(prefix)+len(splitter)+len(splitted[0]):], splitter)
	}
	return CommandReceivedEvent{
		ServerWorker: *serverWorker,
		MsgUUID:      msgUUID,
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
	MsgUUID             uuid.UUID
	Data, DecryptedData string
}

func NewEncryptedDataReceivedEvent(serverWorker *ServerWorker, data, decryptedData string, msgUUID uuid.UUID) EncryptedDataReceivedEvent {
	return EncryptedDataReceivedEvent{ServerWorker: *serverWorker, MsgUUID: msgUUID, Data: data, DecryptedData: decryptedData}
}

type RawDataReceivedEvent struct {
	ServerWorker ServerWorker
	MsgUUID      uuid.UUID
	Data         string
}

func NewRawDataReceivedEvent(serverWorker *ServerWorker, data string, msgUUID uuid.UUID) RawDataReceivedEvent {
	return RawDataReceivedEvent{ServerWorker: *serverWorker, MsgUUID: msgUUID, Data: data}
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
