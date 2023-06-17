from datetime import datetime


class Logger:

    def __init__(self, ILogger):
        self.ILogger = ILogger

    def log(self, message: str):
        print(datetime.now().strftime("%y-%m-%d %H:%M:%S") +
              "\t|\t" + self.ILogger.name + "   |\t" + message)
