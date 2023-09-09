package response_system

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type IResponseWorker interface {
	SendData(string, uuid.UUID, bool) *Response
	SendCommand(command string, args ...string) *Response
	GetResponseManager() ResponseManager
}

type Response struct {
	responseWorker IResponseWorker
	UUID           uuid.UUID
	OnReply        func(IResponseWorker, string)
	replied        chan bool
	Timeout        int64
}

func NewResponse(worker IResponseWorker, UUID uuid.UUID) Response {
	if UUID == uuid.Nil {
		UUID = uuid.New()
	}
	return Response{
		responseWorker: worker,
		UUID:           UUID,
		replied:        make(chan bool),
		Timeout:        3000,
	}
}

func (response *Response) WaitReply(blockThread bool, onReply func(IResponseWorker, string)) (err error) {
	response.OnReply = onReply
	r := response.responseWorker.GetResponseManager()
	r.WaitResponse(response)
	if blockThread {
		now := time.Now().UnixMilli()
		for {
			select {
			case <-response.replied:
				return
			default:
				if time.Now().UnixMilli() > now+response.Timeout {
					return fmt.Errorf("no response from server in %d ms", response.Timeout)
				}
			}
		}
	}
	return
}

func (response *Response) AcceptReply(message string) {
	if response.OnReply != nil {
		response.OnReply(response.responseWorker, message)
	}
	response.replied <- true
}
