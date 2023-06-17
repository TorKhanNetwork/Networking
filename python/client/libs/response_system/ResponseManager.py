class ResponseManager:
    def __init__(self):
        self.responsesWaitingList = dict()

    def waitResponse(self, response):
        self.responsesWaitingList.__setitem__(response.msgUUID, response)

    def onResponseReceived(self, msgUUID, message: str):
        if msgUUID in self.responsesWaitingList:
            self.responsesWaitingList.get(msgUUID).acceptReply(message)
            self.responsesWaitingList.pop(msgUUID)
