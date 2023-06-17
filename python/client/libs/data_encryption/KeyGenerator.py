from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.serialization import PublicFormat, Encoding
import base64


class KeyGenerator:
    def __init__(self, asyncKeySize: int = 2048, syncKeySize: int = 16, iv_size: int = 16, asyncAlgorithm=hashes.SHA256(), syncAlgorithm: str = "aes_256_cbc"):
        self.publicKey = None
        self.privateKey = None
        self.asyncKeySize = asyncKeySize
        self.syncKeySize = syncKeySize
        self.iv_size = iv_size
        self.asyncAlgorithm = asyncAlgorithm
        self.syncAlgorithm = syncAlgorithm
        self.secretKey = None
        self.secretKeyIv = None

    def generateKeys(self, secretKey: bool, asymetricalKeys: bool):
        if secretKey:
            self.secretKey = get_random_bytes(self.syncKeySize)
            self.secretKeyIv = get_random_bytes(self.iv_size)
        if asymetricalKeys:
            self.privateKey = rsa.generate_private_key(
                public_exponent=65537, key_size=2048, backend=default_backend())
            self.publicKey = self.privateKey.public_key()

    def getSecretKey(self):
        return AES.new(self.secretKey, AES.MODE_CBC, self.secretKeyIv)

    def getStringPublicKey(self):
        return base64.b64encode(self.publicKey.public_bytes(encoding=Encoding.DER, format=PublicFormat.SubjectPublicKeyInfo)).decode(
            'utf-8')
