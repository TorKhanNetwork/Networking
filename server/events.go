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

type ClientSocketDisconnectEvent struct {
	ServerWorker ServerWorker
}

func NewClientSocketDisconnectEvent(serverWorker *ServerWorker) ClientSocketDisconnectEvent {
	return ClientSocketDisconnectEvent{ServerWorker: *serverWorker}
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
	bArgs := r.FindAllString(data, -1)
	for _, bArg := range bArgs {
		s := strings.ReplaceAll(string(bArg), `\"`, `"`)
		if len(s) >= 2 && strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
			args = append(args, s[1:len(s)-1])
		}
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
