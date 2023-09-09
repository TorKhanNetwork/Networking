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
	bArgs := r.FindAllString(data, -1)
	for _, bArg := range bArgs {
		s := strings.ReplaceAll(string(bArg), `\"`, `"`)
		if len(s) >= 2 && strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
			args = append(args, s[1:len(s)-1])
		}
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

type ServerDisconnectEvent struct {
	SocketWorker SocketWorker
}

func NewServerDisconnectEvent(socketWorker *SocketWorker) ServerDisconnectEvent {
	return ServerDisconnectEvent{SocketWorker: *socketWorker}
}

type ServerSocketDisconnectEvent struct {
	SocketWorker SocketWorker
}

func NewServerSocketDisconnectEvent(socketWorker *SocketWorker) ServerSocketDisconnectEvent {
	return ServerSocketDisconnectEvent{SocketWorker: *socketWorker}
}
