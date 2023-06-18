package client

import (
	"strings"

	"github.com/google/uuid"
)

type CommandReceivedEvent struct {
	SocketWorker SocketWorker
	MsgUUID      uuid.UUID
	Command      string
	Args         []string
}

func NewCommandReceivedEvent(socketWorker *SocketWorker, data string, prefix string, splitter string, msgUUID uuid.UUID) CommandReceivedEvent {
	splitted := strings.Split(data[len(prefix):], splitter)
	args := make([]string, 0)
	if len(splitted) >= 2 {
		args = strings.Split(data[len(prefix)+len(splitter)+len(splitted[0]):], splitter)
	}
	return CommandReceivedEvent{
		SocketWorker: *socketWorker,
		MsgUUID:      msgUUID,
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
	MsgUUID             uuid.UUID
	Data, DecryptedData string
}

func NewEncryptedDataReceivedEvent(socketWorker *SocketWorker, data, decryptedData string, msgUUID uuid.UUID) EncryptedDataReceivedEvent {
	return EncryptedDataReceivedEvent{SocketWorker: *socketWorker, MsgUUID: msgUUID, Data: data, DecryptedData: decryptedData}
}

type RawDataReceivedEvent struct {
	SocketWorker SocketWorker
	MsgUUID      uuid.UUID
	Data         string
}

func NewRawDataReceivedEvent(socketWorker *SocketWorker, data string, msgUUID uuid.UUID) RawDataReceivedEvent {
	return RawDataReceivedEvent{SocketWorker: *socketWorker, MsgUUID: msgUUID, Data: data}
}

type ServerDisconnectEvent struct {
	SocketWorker SocketWorker
}

func NewServerDisconnectEvent(socketWorker *SocketWorker) ServerDisconnectEvent {
	return ServerDisconnectEvent{SocketWorker: *socketWorker}
}
