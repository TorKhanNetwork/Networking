package response_system

import (
	"github.com/google/uuid"
)

type ResponseManager struct {
	responsesWaitingList map[uuid.UUID]*Response
}

func NewResponseManager() ResponseManager {
	return ResponseManager{
		responsesWaitingList: make(map[uuid.UUID]*Response),
	}
}

func (responseManager *ResponseManager) WaitResponse(response *Response) {
	responseManager.responsesWaitingList[response.UUID] = response
}

func (responseManager *ResponseManager) OnResponseReceived(UUID uuid.UUID, message string) {
	if r, exist := responseManager.responsesWaitingList[UUID]; exist {
		r.AcceptReply(message)
		delete(responseManager.responsesWaitingList, UUID)
	}
}
