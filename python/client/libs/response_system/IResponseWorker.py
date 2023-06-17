import threading
import uuid
from response_system import ResponseManager, Response


class IResponseWorker(threading.Thread):
    def __init__(self):
        threading.Thread.__init__(self)
        self.responseManager = ResponseManager.ResponseManager()

    def sendData(self, data: str, msgUUID: uuid.UUID = None, encrypt: bool = True) -> Response.Response:
        pass

    def sendCommand(self, command: str, args: tuple = None) -> Response.Response:
        pass
