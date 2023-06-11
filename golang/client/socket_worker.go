package client

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/TorkhanNetwork/Networking/golang/data_encryption"
	"github.com/google/uuid"
	"github.com/kataras/golog"
)

const INPUT_BUFFER_SIZE = 65536

type SocketWorker struct {
	Id                                           int
	client                                       Client
	targetIp, commandPrefix, requestSeparator    string
	targetPort                                   int
	connection                                   net.TCPConn
	connected, authenticated, connectionProtocol bool
	keyGenerator                                 data_encryption.KeyGenerator
	quitChan                                     chan bool
}

func NewSocketWorker(id int, client Client, targetIp string, targetPort int) SocketWorker {
	return SocketWorker{
		Id:           id,
		client:       client,
		targetIp:     targetIp,
		targetPort:   targetPort,
		keyGenerator: data_encryption.NewGenerator(),
	}
}

func (socketWorker SocketWorker) GetName() string {
	return "SocketWorker #" + strconv.Itoa(socketWorker.Id)
}

func (socketWorker *SocketWorker) StartWorker() {
	if err := socketWorker.ConnectSocket(); err != nil {
		golog.Error(socketWorker.GetName()+" - Unbale to start worker\n", err)
		return
	}
	socketWorker.quitChan = make(chan bool)
	go socketWorker.handleServerSocket()
}

func (socketWorker *SocketWorker) StopWorker() {
	if socketWorker.connected {
		socketWorker.SendCommand("disconnect")
		socketWorker.quitChan <- true
	}
}

func (socketWorker *SocketWorker) ConnectSocket() error {
	if socketWorker.connected {
		return errors.New(socketWorker.GetName() + " - Socket is already connected")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", socketWorker.targetIp+":"+strconv.Itoa(socketWorker.targetPort))
	if err != nil {
		return err
	}
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	if connection == nil {
		return errors.New(socketWorker.GetName() + " - Unable to connect the socket")
	}
	connection.SetDeadline(time.Time{})
	socketWorker.connection = *connection
	socketWorker.connected = true
	return nil
}

func (socketWorker *SocketWorker) DisconnectSocket() error {
	if !socketWorker.connected {
		return errors.New(socketWorker.GetName() + " - Socket isn't connected")
	}
	err := socketWorker.connection.Close()
	if err != nil {
		return err
	}
	socketWorker.StopWorker()
	socketWorker.connected = false
	socketWorker.client.RemoveSocketWorker(*socketWorker)
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
			socketWorker.DisconnectSocket()
		default:
			for _, line := range strings.Split(string(buf[:mLen]), "\n") {
				socketWorker.onLineRead(line)
			}
		}
	}

}

func (socketWorker *SocketWorker) onLineRead(line string) {
	line = strings.TrimSpace(line)
	if line != "" {
		if decryptedData, err := data_encryption.Decrypt(line, socketWorker.keyGenerator); err == nil {
			socketWorker.onDataReceived([]string{decryptedData, line}, true)
		} else {
			socketWorker.onDataReceived([]string{line}, false)
		}
	}

}

func (socketWorker *SocketWorker) onDataReceived(data []string, encrypted bool) {
	var msgUUID uuid.UUID
	if strings.HasPrefix(data[0], "response:") {
		msgUUID, err := uuid.Parse(data[0][9:45])
		if err != nil {
			golog.Error(socketWorker.GetName()+" - Unable to parse response UUID\n", err)
			return
		}
		data[0] = data[0][45:]
		golog.Debug(socketWorker.GetName() + " - Response received for UUID " + msgUUID.String())
	} else {
		var err error
		msgUUID, err = uuid.Parse(data[0][:36])
		if err != nil {
			golog.Error(socketWorker.GetName()+" - Unable to parse response UUID\n", err)
			return
		}
		data[0] = data[0][36:]
	}

	if socketWorker.commandPrefix != "" && strings.HasPrefix(data[0], socketWorker.commandPrefix) {
		socketWorker.onCommandReceived(data[0], msgUUID, encrypted)
	} else {
		golog.Debug(socketWorker.GetName() + " - Data received (encrypted=" + strconv.FormatBool(encrypted) + ", uuid=" + msgUUID.String() + ") : " + data[0])
	}

	// TODO event

	if socketWorker.connectionProtocol {
		if !socketWorker.onConnectionProtocolDataReceived(data[0]) {
			golog.Error(socketWorker.GetName() + " - Unable to connect to the server")
			socketWorker.StopWorker()
		}
	}

}

func (socketWorker *SocketWorker) onCommandReceived(command string, msgUUID uuid.UUID, encrypted bool) {
	golog.Debug(socketWorker.GetName() + " - Command Received (encrypted=" + strconv.FormatBool(encrypted) + ", uuid=" + msgUUID.String() + ") : " + command)
	// TODO Event
}

func (socketWorker *SocketWorker) startConnectionProtocol() {
	socketWorker.connectionProtocol = true
	socketWorker.SendData("version:"+CLIENT_VERSION, uuid.NullUUID{}, false)
}

func (socketWorker *SocketWorker) onConnectionProtocolDataReceived(data string) bool {
	if strings.HasPrefix(data, "version:") {
		if version := strings.ToLower(data[8:]) == "true"; version {
			socketWorker.SendData("macAddress:"+socketWorker.client.MacAddress, uuid.NullUUID{}, false)
		} else {
			return false
		}
	} else if strings.HasPrefix(data, "publicKey:") {
		pub, err := ParseRsaPublicKeyFromPemStr(data[10:])
		if err != nil {
			golog.Error(socketWorker.GetName()+" - Failed to parse server public key\n", err)
			return false
		}
		socketWorker.keyGenerator.PublicKey = *pub
		socketWorker.keyGenerator.GenerateKeys(true, false)
		secretKey, err := data_encryption.EncryptSecretKey(socketWorker.keyGenerator)
		if err != nil {
			golog.Error("Unable to encrypt secret key\n", err)
			return false
		}
		socketWorker.SendData("secretKey:"+secretKey, uuid.NullUUID{}, false)
	} else if strings.HasPrefix(data, "separator:") {
		socketWorker.requestSeparator = data[9:]
	} else if strings.HasPrefix(data, "commandPrefix:") {
		socketWorker.commandPrefix = data[14:]
	} else if strings.HasPrefix(data, "connected:") {
		connected := strings.ToLower(data[10:]) == "true"
		socketWorker.authenticated = connected
		socketWorker.connectionProtocol = !connected
		// TODO event
		if connected {
			socketWorker.SendCommand("blou")
		}
		return connected
	} else {
		return false
	}
	return true
}

func (socketWorker *SocketWorker) SendData(data string, responseUUID uuid.NullUUID, encrypt bool) {
	rawData := data
	// TODO response
	data = uuid.New().String() + data
	if encrypt {
		var err error
		data, err = data_encryption.Encrypt(data, socketWorker.keyGenerator)
		if err != nil {
			golog.Error(socketWorker.GetName()+" - Unable to encrypt data\n", err)
			return
		}
	}
	_, err := socketWorker.connection.Write([]byte(data + "\n"))
	if err != nil {
		golog.Error(socketWorker.GetName()+" - Unable to send data\n", err)
		return
	}
	golog.Debug(socketWorker.GetName() + " - Data sent (encrypted=" + strconv.FormatBool(encrypt) + ", uuid=not supported yet) : " + rawData)
}

func (socketWorker *SocketWorker) SendCommand(command string, args ...string) {
	socketWorker.SendData(socketWorker.commandPrefix+command+socketWorker.requestSeparator+strings.Join(args, socketWorker.requestSeparator), uuid.NullUUID{}, true)
}
