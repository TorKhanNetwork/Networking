package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/TorkhanNetwork/Networking/golang/data_encryption"
	"github.com/TorkhanNetwork/Networking/golang/events_system"
	"github.com/TorkhanNetwork/Networking/golang/response_system"
	"github.com/google/uuid"
	"github.com/kataras/golog"
)

const INPUT_BUFFER_SIZE = 65536

type ServerWorker struct {
	Id                                           int
	server                                       *Server
	commandPrefix, requestSeparator, MacAddress  string
	connection                                   *net.TCPConn
	responseManager                              response_system.ResponseManager
	connected, authenticated, connectionProtocol bool
	quitChan                                     chan bool
}

func NewServerWorker(id int, server *Server, connection *net.TCPConn) ServerWorker {
	return ServerWorker{
		Id:               id,
		server:           server,
		commandPrefix:    "!",
		requestSeparator: server.GenerateRequestSeparator(),
		connection:       connection,
		responseManager:  response_system.NewResponseManager(),
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
	go serverWorker.handleClientSocket()
}

func (serverWorker *ServerWorker) StopWorker() {
	if serverWorker.connected {
		serverWorker.SendData("disconnect", uuid.Nil, serverWorker.authenticated)
		serverWorker.quitChan <- true
	}
}

func (serverWorker *ServerWorker) DisconnectSocket() error {
	if !serverWorker.connected {
		return fmt.Errorf("%s - Socket isn't connected", serverWorker.GetName())
	}
	err := serverWorker.connection.Close()
	if err != nil {
		return err
	}
	serverWorker.StopWorker()
	serverWorker.connected = false
	serverWorker.server.RemoveServerWorker(*serverWorker)
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
	var e interface{} = NewClientDisconnectEvent(serverWorker)
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

	if serverWorker.commandPrefix != "" && strings.HasPrefix(data[0], serverWorker.commandPrefix) {
		serverWorker.onCommandReceived(data[0], msgUUID, encrypted)
	} else {
		golog.Debugf("%s - Data received (encrypted=%t, uuid=%s) : %s", serverWorker.GetName(), encrypted, msgUUID, data[0])
	}

	var e interface{}
	if encrypted {
		e = NewEncryptedDataReceivedEvent(serverWorker, data[1], data[0], msgUUID)
	} else {
		e = NewRawDataReceivedEvent(serverWorker, data[0], msgUUID)
	}
	serverWorker.server.EventsManager.CallEvent((*events_system.Event)(&e))

	if serverWorker.connectionProtocol {
		if !serverWorker.onConnectionProtocolDataReceived(data[0]) {
			golog.Errorf("%s - Unable to establish a connection with the client", serverWorker.GetName())
			serverWorker.StopWorker()
		}
	}

}

func (serverWorker *ServerWorker) onCommandReceived(command string, msgUUID uuid.UUID, encrypted bool) {
	golog.Debugf("%s - Command Received (encrypted=%t, uuid=%s) : %s", serverWorker.GetName(), encrypted, msgUUID, command)
	var event interface{} = NewCommandReceivedEvent(serverWorker, command, serverWorker.commandPrefix, serverWorker.requestSeparator, msgUUID)
	serverWorker.server.EventsManager.CallEvent((*events_system.Event)(&event))
}

func (serverWorker *ServerWorker) startConnectionProtocol() {
	serverWorker.connectionProtocol = true
}

func (serverWorker *ServerWorker) onConnectionProtocolDataReceived(data string) bool {
	if strings.HasPrefix(data, "version:") {
		split := strings.SplitN(data[8:], ".", 3)
		correct := len(split) >= 2 && strings.Join(split[:2], ".") == strings.Join(strings.SplitN(SERVER_VERSION, ".", 3)[:2], ".")
		serverWorker.SendData("version:"+strconv.FormatBool(correct), uuid.Nil, false)
		if !correct {
			return false
		}
		serverWorker.SendData("commandPrefix:"+serverWorker.commandPrefix, uuid.Nil, false)
		serverWorker.SendData("separator:"+serverWorker.requestSeparator, uuid.Nil, false)
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
		serverWorker.SendData("connected:"+strconv.FormatBool(connected), uuid.Nil, false)
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
	return serverWorker.SendData(serverWorker.commandPrefix+command+serverWorker.requestSeparator+strings.Join(args, serverWorker.requestSeparator), uuid.Nil, true)
}
