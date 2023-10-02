package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/TorKhanNetwork/Networking/libs/data_encryption"
	"github.com/TorKhanNetwork/Networking/libs/events_system"
	"github.com/TorKhanNetwork/Networking/libs/response_system"
	"github.com/google/uuid"
	"github.com/kataras/golog"
)

const INPUT_BUFFER_SIZE = 65536

type ServerWorker struct {
	Id                                           int
	server                                       *Server
	commandPrefix, MacAddress                    string
	connection                                   *net.TCPConn
	responseManager                              response_system.ResponseManager
	connected, authenticated, connectionProtocol bool
	quitChan                                     chan bool
}

func NewServerWorker(id int, server *Server, connection *net.TCPConn) ServerWorker {
	return ServerWorker{
		Id:              id,
		server:          server,
		commandPrefix:   COMMAND_PREFIX,
		connection:      connection,
		responseManager: response_system.NewResponseManager(),
	}
}

func (serverWorker ServerWorker) GetName() string {
	return fmt.Sprintf("  |  %s  |  ServerWorker #%d ", serverWorker.server.Name, serverWorker.Id)
}

func (serverWorker *ServerWorker) GetResponseManager() response_system.ResponseManager {
	return serverWorker.responseManager
}

func (serverWorker *ServerWorker) StartWorker() {
	serverWorker.quitChan = make(chan bool)
	serverWorker.connected = true
	go serverWorker.handleClientSocket()
}

func (serverWorker *ServerWorker) StopWorker() {
	if serverWorker.connected {
		serverWorker.server.RemoveServerWorker(*serverWorker)
		serverWorker.quitChan <- true
	}
}

func (serverWorker *ServerWorker) DisconnectSocket() error {
	if !serverWorker.connected {
		return fmt.Errorf("%s - Socket isn't connected", serverWorker.GetName())
	}
	serverWorker.SendData("disconnect", uuid.Nil, serverWorker.authenticated)
	err := serverWorker.connection.Close()
	if err != nil {
		return err
	}
	serverWorker.connected = false
	serverWorker.StopWorker()
	return nil
}

func (serverWorker *ServerWorker) handleClientSocket() {
	buf := make([]byte, INPUT_BUFFER_SIZE)
	serverWorker.startConnectionProtocol()
	for {
		mLen, err := serverWorker.connection.Read(buf)
		if err != nil {
			break
		}
		select {
		case <-serverWorker.quitChan:
			serverWorker.DisconnectSocket()
		default:
			for _, line := range strings.Split(string(buf[:mLen]), "\n") {
				serverWorker.onLineRead(line)
			}
		}
	}
	var e interface{} = NewClientSocketDisconnectEvent(serverWorker)
	serverWorker.server.EventsManager.CallEvent((*events_system.Event)(&e))
}

func (serverWorker *ServerWorker) onLineRead(line string) {
	line = strings.TrimSpace(line)
	if line != "" {
		if decryptedData, err := data_encryption.Decrypt(line, serverWorker.server.keysGenerator); err == nil {
			serverWorker.onDataReceived([]string{decryptedData, line}, true)
		} else {
			serverWorker.onDataReceived([]string{line}, false)
		}
	}

}

func (serverWorker *ServerWorker) onDataReceived(data []string, encrypted bool) {
	var msgUUID uuid.UUID
	var err error
	if strings.HasPrefix(data[0], "response:") {
		msgUUID, err = uuid.Parse(data[0][9:45])
		if err != nil {
			golog.Errorf("%s - Unable to parse response UUID : %s", serverWorker.GetName(), err)
			return
		}
		data[0] = data[0][45:]
		golog.Debugf("%s - Response received for UUID %s", serverWorker.GetName(), msgUUID.String())
		go serverWorker.responseManager.OnResponseReceived(msgUUID, data[0])
	} else {
		var err error
		msgUUID, err = uuid.Parse(data[0][:36])
		if err != nil {
			golog.Errorf("%s - Unable to parse response UUID : %s", serverWorker.GetName(), err)
			return
		}
		data[0] = data[0][36:]
	}

	if serverWorker.connectionProtocol {
		golog.Debugf("%s - Data received (encrypted=%t, uuid=%s) : %s", serverWorker.GetName(), encrypted, msgUUID, data[0])
		if !serverWorker.onConnectionProtocolDataReceived(data[0]) {
			golog.Errorf("%s - Unable to authenticate the client", serverWorker.GetName())
			err := serverWorker.DisconnectSocket()
			if err != nil {
				golog.Errorf("%s - Unable to disconnect socket : %s", serverWorker.GetName(), err)
			}
		}
	} else {
		if serverWorker.commandPrefix != "" && strings.HasPrefix(data[0], serverWorker.commandPrefix) {
			serverWorker.onCommandReceived(data[0], msgUUID, encrypted)
		} else {
			golog.Debugf("%s - Data received (encrypted=%t, uuid=%s) : %s", serverWorker.GetName(), encrypted, msgUUID, data[0])
			if data[0] == "disconnect" {
				serverWorker.StopWorker()
				var e events_system.Event = NewClientDisconnectEvent(serverWorker)
				serverWorker.server.EventsManager.CallEvent(&e)
			}
		}
	}

}

func (serverWorker *ServerWorker) onCommandReceived(command string, msgUUID uuid.UUID, encrypted bool) {
	golog.Debugf("%s - Command Received (encrypted=%t, uuid=%s) : %s", serverWorker.GetName(), encrypted, msgUUID, command)
	var event interface{} = NewCommandReceivedEvent(serverWorker, command, serverWorker.commandPrefix, msgUUID)
	serverWorker.server.EventsManager.CallEvent((*events_system.Event)(&event))
}

func (serverWorker *ServerWorker) startConnectionProtocol() {
	serverWorker.connectionProtocol = true
}

func (serverWorker *ServerWorker) onConnectionProtocolDataReceived(data string) bool {
	if strings.HasPrefix(data, "version:") {
		split := strings.SplitN(data[8:], ".", 3)
		correct := len(split) == 3 && strings.Join(split[:2], ".") == strings.Join(strings.SplitN(SERVER_VERSION, ".", 3)[:2], ".")
		serverWorker.SendData("version:"+strconv.FormatBool(correct), uuid.Nil, false)
		if !correct {
			return false
		}
	} else if strings.HasPrefix(data, "macAddress:") {
		serverWorker.MacAddress = strings.ToUpper(data[11:])
		publicKey, err := ExportRsaPublicKeyToString(&serverWorker.server.keysGenerator.PublicKey)
		if err != nil {
			return false
		}
		serverWorker.SendData("publicKey:"+publicKey, uuid.Nil, false)
	} else if strings.HasPrefix(data, "secretKey:") {
		err := data_encryption.DecryptSecretKey(data[10:], &serverWorker.server.keysGenerator)
		connected := err == nil
		if connected {
			serverWorker.SendData("commandPrefix:"+serverWorker.commandPrefix, uuid.Nil, true)
		}
		serverWorker.SendData("connected:"+strconv.FormatBool(connected), uuid.Nil, true)
		serverWorker.connectionProtocol = !connected
		if connected {
			var e interface{} = NewConnectionProtocolSuccessEvent(serverWorker)
			serverWorker.server.EventsManager.CallEvent((*events_system.Event)(&e))
		}
		return connected
	} else {
		return false
	}
	return true
}

func (serverWorker *ServerWorker) SendData(data string, responseUUID uuid.UUID, encrypt bool) *response_system.Response {
	rawData := data
	response := response_system.NewResponse(serverWorker, responseUUID)
	data = response.UUID.String() + data
	if responseUUID != uuid.Nil {
		data = "response:" + data
	}
	if encrypt {
		var err error
		data, err = data_encryption.Encrypt(data, serverWorker.server.keysGenerator)
		if err != nil {
			golog.Errorf("%s - Unable to encrypt data : %s", serverWorker.GetName(), err)
			return nil
		}
	}
	_, err := serverWorker.connection.Write([]byte(data + "\n"))
	if err != nil {
		golog.Errorf("%s - Unable to send data : %s", serverWorker.GetName(), err)
		return nil
	}
	golog.Debugf("%s - Data sent (encrypted=%t, uuid=%s) : %s", serverWorker.GetName(), encrypt, response.UUID, rawData)
	return &response
}

func (serverWorker *ServerWorker) SendCommand(command string, args ...string) *response_system.Response {
	data := serverWorker.commandPrefix + command
	for _, arg := range args {
		arg = strings.ReplaceAll(arg, "\"", "\\\"")
		data += " \"" + arg + "\""
	}
	return serverWorker.SendData(data, uuid.Nil, true)
}
