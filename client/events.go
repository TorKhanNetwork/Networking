package client

import (
	"regexp"
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
	splitted := strings.SplitN(data[len(prefix):], " ", 2)
	args := make([]string, 0)
	r, _ := regexp.Compile(`"((\\")|[^"])+"`)
	bArgs := r.FindAll([]byte(data), -1)
	for _, bArg := range bArgs {
		args = append(args, string(bArg)[1:len(bArg)-1])
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
