from zope.interface import Interface


class Event(Interface):
    def __init__(self):
        self.name = None

    def getName(self):
        if self.name is None:
            return Interface.__name__
        else:
            return self.name
