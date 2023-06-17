import base64
import M2Crypto


def decryptPublicKey(publicKeyStr):
    bio = M2Crypto.BIO.MemoryBuffer(
        ("-----BEGIN PUBLIC KEY-----\n" + publicKeyStr + "\n-----END PUBLIC KEY-----").encode())
    return M2Crypto.RSA.load_pub_key_bio(bio)


def encryptSecretKey(keyGenerator) -> str:
    return base64.b64encode(keyGenerator.publicKey.public_encrypt(keyGenerator.secretKeyIv + keyGenerator.secretKey, M2Crypto.RSA.pkcs1_padding)).decode('utf-8')


# def decryptSecretKey(secretKey: list, privateKey):
#     return [privateKey.decrypt(base64.b64decode(byteElement.encode('utf-8')), padding=padding.OAEP(
#         mgf=padding.MGF1(algorithm=hashes.SHA256()),
#         algorithm=hashes.SHA256(),
#         label=None
#     )) for byteElement in secretKey]


def pad(byte_array: bytearray, keyGenerator):
    pad_len = keyGenerator.syncKeySize - \
        len(byte_array) % keyGenerator.syncKeySize
    return byte_array + (bytes([pad_len]) * pad_len)


def unpad(byte_array: bytearray):
    return byte_array[:-ord(byte_array[-1:])]


def encrypt(data: str, keyGenerator):
    try:
        return base64.b64encode(keyGenerator.getSecretKey().encrypt(
            pad(data.encode(), keyGenerator))).decode()
    except Exception:
        return None


def decrypt(data: str, keyGenerator):
    try:
        return unpad(
            keyGenerator.getSecretKey().decrypt(
                base64.b64decode(
                    data.encode()
                )
            ).decode()
        )
    except Exception:
        return None
