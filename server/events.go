package server

import (
	"net"
	"regexp"
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

func NewCommandReceivedEvent(serverWorker *ServerWorker, data string, prefix string, msgUUID uuid.UUID) CommandReceivedEvent {
	splitted := strings.SplitN(data[len(prefix):], " ", 2)
	args := make([]string, 0)
	r, _ := regexp.Compile(`"((\\")|[^"])+"`)
	bArgs := r.FindAll([]byte(data), -1)
	for _, bArg := range bArgs {
		args = append(args, string(bArg)[1:len(bArg)-1])
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
