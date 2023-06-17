import uuid


class Response:
    def __init__(self, responseWorker, msgUUID: uuid.UUID):
        self.responseWorker = responseWorker
        if msgUUID is not None:
            self.msgUUID = msgUUID
        else:
            self.msgUUID = uuid.uuid4()

    def acceptReply(self, message: str):
        pass
