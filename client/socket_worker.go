package client

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/TorKhanNetwork/Networking/data_encryption"
	"github.com/TorKhanNetwork/Networking/events_system"
	"github.com/TorKhanNetwork/Networking/response_system"
	"github.com/google/uuid"
	"github.com/kataras/golog"
)

const INPUT_BUFFER_SIZE = 65536

type SocketWorker struct {
	Id                                           int
	client                                       *Client
	targetIp, commandPrefix, requestSeparator    string
	targetPort                                   int
	connection                                   *net.TCPConn
	connected, authenticated, connectionProtocol bool
	keysGenerator                                data_encryption.KeysGenerator
	responseManager                              response_system.ResponseManager
	quitChan                                     chan bool
}

func NewSocketWorker(id int, client *Client, targetIp string, targetPort int) SocketWorker {
	return SocketWorker{
		Id:              id,
		client:          client,
		targetIp:        targetIp,
		targetPort:      targetPort,
		keysGenerator:   data_encryption.NewGenerator(),
		responseManager: response_system.NewResponseManager(),
	}
}

func (socketWorker *SocketWorker) GetName() string {
	return fmt.Sprintf("  |  %s  |  SocketWorker #%d ", socketWorker.client.Name, socketWorker.Id)
}

func (socketWorker *SocketWorker) GetResponseManager() response_system.ResponseManager {
	return socketWorker.responseManager
}

func (socketWorker *SocketWorker) StartWorker() (err error) {
	if err = socketWorker.ConnectSocket(); err != nil {
		return err
	}
	socketWorker.quitChan = make(chan bool)
	go socketWorker.handleServerSocket()
	return
}

func (socketWorker *SocketWorker) StopWorker() {
	if socketWorker.connected {
		socketWorker.client.RemoveSocketWorker(*socketWorker)
		socketWorker.quitChan <- true
	}
}

func (socketWorker *SocketWorker) ConnectSocket() error {
	if socketWorker.connected {
		return errors.New("socket is already connected")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", socketWorker.targetIp+":"+strconv.Itoa(socketWorker.targetPort))
	if err != nil {
		return fmt.Errorf("unable to resolve TCP Address : %s", err)
	}
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return fmt.Errorf("unable to connect the socket : %s", err)
	}
	if connection == nil {
		return errors.New("connection to the socket is nil")
	}
	connection.SetDeadline(time.Time{})
	socketWorker.connection = connection
	socketWorker.connected = true
	return nil
}

func (socketWorker *SocketWorker) DisconnectSocket() error {
	if !socketWorker.connected {
		return errors.New("socket isn't connected")
	}
	socketWorker.SendData("disconnect", uuid.Nil, socketWorker.authenticated)
	err := socketWorker.connection.Close()
	if err != nil {
		return fmt.Errorf("unable to close socket : %s", err)
	}
	socketWorker.connected = false
	socketWorker.StopWorker()
	return nil
}

func (socketWorker *SocketWorker) handleServerSocket() {
	buf := make([]byte, INPUT_BUFFER_SIZE)
	socketWorker.startConnectionProtocol()
	for {
		mLen, err := socketWorker.connection.Read(buf)
		if err != nil {
			break
		}
		select {
		case <-socketWorker.quitChan:
			if err := socketWorker.DisconnectSocket(); err != nil {
				golog.Errorf("%s - %s", socketWorker.GetName(), err)
			}
			return
		default:
			for _, line := range strings.Split(string(buf[:mLen]), "\n") {
				socketWorker.onLineRead(line)
			}
		}
	}
	var e events_system.Event = NewServerSocketDisconnectEvent(socketWorker)
	socketWorker.client.EventsManager.CallEvent(&e)
}

func (socketWorker *SocketWorker) onLineRead(line string) {
	line = strings.TrimSpace(line)
	if line != "" {
		if decryptedData, err := data_encryption.Decrypt(line, socketWorker.keysGenerator); err == nil {
			socketWorker.onDataReceived([]string{decryptedData, line}, true)
		} else {
			socketWorker.onDataReceived([]string{line}, false)
		}
	}

}

func (socketWorker *SocketWorker) onDataReceived(data []string, encrypted bool) {
	var msgUUID uuid.UUID
	var err error
	if strings.HasPrefix(data[0], "response:") {
		msgUUID, err = uuid.Parse(data[0][9:45])
		if err != nil {
			golog.Errorf("%s - Unable to parse response UUID : %s", socketWorker.GetName(), err)
			return
		}
		data[0] = data[0][45:]
		golog.Debugf("%s - Response received for UUID %s", socketWorker.GetName(), msgUUID)
		go socketWorker.responseManager.OnResponseReceived(msgUUID, data[0])
	} else {
		var err error
		msgUUID, err = uuid.Parse(data[0][:36])
		if err != nil {
			golog.Errorf("%s - Unable to parse response UUID : %s", socketWorker.GetName(), err)
			return
		}
		data[0] = data[0][36:]
	}

	if socketWorker.connectionProtocol {
		golog.Debugf("%s - Data received (encrypted=%t, uuid=%s) : %s", socketWorker.GetName(), encrypted, msgUUID, data[0])
		if !socketWorker.onConnectionProtocolDataReceived(data[0]) {
			golog.Default.Errorf("%s - Unable to authenticate to the server", socketWorker.GetName())
		}
	} else {
		if socketWorker.commandPrefix != "" && strings.HasPrefix(data[0], socketWorker.commandPrefix) {
			socketWorker.onCommandReceived(data[0], msgUUID, encrypted)
		} else {
			golog.Debugf("%s - Data received (encrypted=%t, uuid=%s) : %s", socketWorker.GetName(), encrypted, msgUUID, data[0])
			if data[0] == "disconnect" {
				socketWorker.StopWorker()
				var e events_system.Event = NewServerDisconnectEvent(socketWorker)
				socketWorker.client.EventsManager.CallEvent(&e)
			}
		}
	}

}

func (socketWorker *SocketWorker) onCommandReceived(command string, msgUUID uuid.UUID, encrypted bool) {
	golog.Debugf("%s - Command Received (encrypted=%t, uuid=%s) : %s", socketWorker.GetName(), encrypted, msgUUID, command)
	var event events_system.Event = NewCommandReceivedEvent(socketWorker, command, socketWorker.commandPrefix, socketWorker.requestSeparator, msgUUID)
	socketWorker.client.EventsManager.CallEvent(&event)
}

func (socketWorker *SocketWorker) startConnectionProtocol() {
	socketWorker.connectionProtocol = true
	socketWorker.SendData("version:"+CLIENT_VERSION, uuid.Nil, false)
}

func (socketWorker *SocketWorker) onConnectionProtocolDataReceived(data string) bool {
	if strings.HasPrefix(data, "version:") {
		if version := strings.ToLower(data[8:]) == "true"; version {
			socketWorker.SendData("macAddress:"+socketWorker.client.MacAddress, uuid.Nil, false)
		} else {
			return false
		}
	} else if strings.HasPrefix(data, "publicKey:") {
		pub, err := ParseRsaPublicKeyFromPemStr(data[10:])
		if err != nil {
			golog.Errorf("%s - Failed to parse server public key : %s", socketWorker.GetName(), err)
			return false
		}
		socketWorker.keysGenerator.PublicKey = *pub
		socketWorker.keysGenerator.GenerateKeys(true, false)
		secretKey, err := data_encryption.EncryptSecretKey(socketWorker.keysGenerator)
		if err != nil {
			golog.Errorf("%s - Unable to encrypt secret key : %s", socketWorker.GetName(), err)
			return false
		}
		socketWorker.SendData("secretKey:"+secretKey, uuid.Nil, false)
	} else if strings.HasPrefix(data, "commandPrefix:") {
		socketWorker.commandPrefix = data[14:]
	} else if strings.HasPrefix(data, "connected:") {
		connected := strings.ToLower(data[10:]) == "true"
		socketWorker.authenticated = connected
		socketWorker.connectionProtocol = !connected
		if connected {
			var e events_system.Event = NewConnectionProtocolSuccessEvent(socketWorker)
			socketWorker.client.EventsManager.CallEvent(&e)
		}
		return connected
	} else {
		return false
	}
	return true
}

func (socketWorker *SocketWorker) SendData(data string, responseUUID uuid.UUID, encrypt bool) *response_system.Response {
	rawData := data
	response := response_system.NewResponse(socketWorker, responseUUID)
	data = response.UUID.String() + data
	if responseUUID != uuid.Nil {
		data = "response:" + data
	}
	if encrypt {
		var err error
		data, err = data_encryption.Encrypt(data, socketWorker.keysGenerator)
		if err != nil {
			golog.Errorf("%s - Unable to encrypt data : %s", socketWorker.GetName(), err)
			return nil
		}
	}
	_, err := socketWorker.connection.Write([]byte(data + "\n"))
	if err != nil {
		golog.Errorf("%s - Unable to send data : %s", socketWorker.GetName(), err)
		return nil
	}
	golog.Debugf("%s - Data sent (encrypted=%t, uuid=%s) : %s", socketWorker.GetName(), encrypt, response.UUID, rawData)
	return &response
}

func (socketWorker *SocketWorker) SendCommand(command string, args ...string) *response_system.Response {
	data := socketWorker.commandPrefix + command
	for _, arg := range args {
		arg = strings.ReplaceAll(arg, "\"", "\\\"")
		data += " \"" + arg + "\""
	}
	return socketWorker.SendData(data, uuid.Nil, true)
}
