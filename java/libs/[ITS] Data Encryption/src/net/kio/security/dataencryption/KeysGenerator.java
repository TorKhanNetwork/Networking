package net.kio.security.dataencryption;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.security.*;
import java.util.Base64;
import java.util.Random;

public class KeysGenerator {
    private final int asyncKeySize;
    private final int syncKeySize;
    private final String asyncAlgorithm;
    private final String syncAlgorithm;
    private KeyPairGenerator keyPairGenerator;
    private SecretKey secretKey;
    private byte[] secretKeyIv;
    private PublicKey publicKey;
    private PrivateKey privateKey;

    public KeysGenerator(int asyncKeySize, int syncKeySize, String asyncAlgorithm, String syncAlgorithm) {
        this.secretKey = null;
        this.publicKey = null;
        this.privateKey = null;
        this.asyncKeySize = asyncKeySize;
        this.syncKeySize = syncKeySize;
        this.asyncAlgorithm = asyncAlgorithm;
        this.syncAlgorithm = syncAlgorithm;
        this.initializeGenerators();
    }

    public KeysGenerator() {
        this(2048, 16);
    }

    public KeysGenerator(int asyncKeySize, int syncKeySize) {
        this(asyncKeySize, syncKeySize, "RSA/None/PKCS1PADDING", "AES/CBC/PKCS5Padding");
    }

    private void initializeGenerators() {
        try {
            this.keyPairGenerator = KeyPairGenerator.getInstance(this.asyncAlgorithm.split("/")[0]);
            this.keyPairGenerator.initialize(this.asyncKeySize);
        } catch (NoSuchAlgorithmException var2) {
            var2.printStackTrace();
        }

    }

    public void generateKeys(boolean secretKey, boolean asynchronousKeys) {
        if (secretKey) {
            this.secretKeyIv = new byte[16];
            new Random().nextBytes(this.secretKeyIv);
            byte[] secretKeyBytes = new byte[this.syncKeySize];
            new Random().nextBytes(secretKeyBytes);
            this.secretKey = new SecretKeySpec(secretKeyBytes, 0, this.syncKeySize, "AES");
        }

        if (asynchronousKeys) {
            KeyPair keyPair = this.keyPairGenerator.generateKeyPair();
            this.privateKey = keyPair.getPrivate();
            this.publicKey = keyPair.getPublic();
        }

    }

    public PrivateKey getPrivateKey() {
        return this.privateKey;
    }

    public PublicKey getPublicKey() {
        return this.publicKey;
    }

    public void setPublicKey(PublicKey publicKey) {
        this.publicKey = publicKey;
    }

    public String getStringPublicKey() {
        return Base64.getEncoder().encodeToString(publicKey.getEncoded());

    }

    public SecretKey getSecretKey() {
        return this.secretKey;
    }

    public void setSecretKey(SecretKey secretKey) {
        this.secretKey = secretKey;
    }

    public byte[] getSecretKeyIv() {
        return secretKeyIv;
    }

    public void setSecretKeyIv(byte[] secretKeyIv) {
        this.secretKeyIv = secretKeyIv;
    }

    public String getAsyncAlgorithm() {
        return asyncAlgorithm;
    }

    public String getSyncAlgorithm() {
        return syncAlgorithm;
    }

    public KeysGenerator clone() throws CloneNotSupportedException {
        return (KeysGenerator) super.clone();
    }
}
