package net.kio.security.dataencryption;

import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.InvalidAlgorithmParameterException;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;

public class EncryptedRequestManager {

    public static String encryptSecretKey(KeysGenerator keysGenerator) {
        try {
            return Base64.getEncoder().encodeToString(encryptSecretKey0(keysGenerator));
        } catch (NoSuchAlgorithmException | InvalidKeyException | IllegalBlockSizeException | BadPaddingException | NoSuchPaddingException var3) {
            var3.printStackTrace();
            return null;
        }
    }

    private static byte[] encryptSecretKey0(KeysGenerator keysGenerator) throws NoSuchPaddingException, NoSuchAlgorithmException, InvalidKeyException, BadPaddingException, IllegalBlockSizeException {
        Cipher cipher = Cipher.getInstance(keysGenerator.getAsyncAlgorithm());
        cipher.init(1, keysGenerator.getPublicKey());
        byte[] secretKeyInfo = new byte[keysGenerator.getSecretKeyIv().length + keysGenerator.getSecretKey().getEncoded().length];
        System.arraycopy(keysGenerator.getSecretKeyIv(), 0, secretKeyInfo, 0, keysGenerator.getSecretKeyIv().length);
        System.arraycopy(keysGenerator.getSecretKey().getEncoded(), 0, secretKeyInfo, keysGenerator.getSecretKeyIv().length, keysGenerator.getSecretKey().getEncoded().length);
        return cipher.doFinal(secretKeyInfo);
    }

    public static void decryptSecretKey(String data, KeysGenerator keysGenerator) {
        try {
            decryptSecretKey0(Base64.getDecoder().decode(data), keysGenerator);
        } catch (NoSuchAlgorithmException | InvalidKeyException | IllegalBlockSizeException | BadPaddingException | NoSuchPaddingException | InvalidAlgorithmParameterException var3) {
            var3.printStackTrace();
        }
    }

    private static void decryptSecretKey0(byte[] data, KeysGenerator keysGenerator) throws NoSuchPaddingException, NoSuchAlgorithmException, InvalidKeyException, BadPaddingException, IllegalBlockSizeException, InvalidAlgorithmParameterException {
        Cipher cipher = Cipher.getInstance(keysGenerator.getAsyncAlgorithm());
        cipher.init(Cipher.DECRYPT_MODE, keysGenerator.getPrivateKey());
        byte[] decryptedKey = cipher.doFinal(data);
        byte[] iv = new byte[16];
        System.arraycopy(decryptedKey, 0, iv, 0, 16);
        keysGenerator.setSecretKeyIv(iv);
        keysGenerator.setSecretKey(new SecretKeySpec(decryptedKey, 16, 16, "AES"));
    }


    public static String encrypt(String data, KeysGenerator keysGenerator) {
        try {
            return encrypt0(data, keysGenerator);
        } catch (InvalidKeyException | NoSuchPaddingException | BadPaddingException | IllegalBlockSizeException | NoSuchAlgorithmException | InvalidAlgorithmParameterException var3) {
            return null;
        }
    }

    private static String encrypt0(String data, KeysGenerator keysGenerator) throws NoSuchAlgorithmException, InvalidKeyException, NoSuchPaddingException, BadPaddingException, IllegalBlockSizeException, InvalidAlgorithmParameterException {
        Cipher aesCipher = Cipher.getInstance(keysGenerator.getSyncAlgorithm());
        aesCipher.init(Cipher.ENCRYPT_MODE, keysGenerator.getSecretKey(), new IvParameterSpec(keysGenerator.getSecretKeyIv()));
        byte[] byteCipherText = aesCipher.doFinal(data.getBytes(StandardCharsets.UTF_8));
        return Base64.getEncoder().encodeToString(byteCipherText);
    }

    public static String decrypt(String data, KeysGenerator keysGenerator) {
        try {
            return decrypt0(Base64.getDecoder().decode(data), keysGenerator);
        } catch (NoSuchPaddingException | NullPointerException | InvalidAlgorithmParameterException | InvalidKeyException | IllegalBlockSizeException | BadPaddingException | NoSuchAlgorithmException | IllegalArgumentException var3) {
            return null;
        }
    }

    private static String decrypt0(byte[] data, KeysGenerator keysGenerator) throws InvalidKeyException, NoSuchPaddingException, NoSuchAlgorithmException, BadPaddingException, IllegalBlockSizeException, NullPointerException, IllegalArgumentException, InvalidAlgorithmParameterException {
        Cipher aesCipher = Cipher.getInstance(keysGenerator.getSyncAlgorithm());
        aesCipher.init(Cipher.DECRYPT_MODE, keysGenerator.getSecretKey(), new IvParameterSpec(keysGenerator.getSecretKeyIv()));
        byte[] bytePlainText = aesCipher.doFinal(data);
        return new String(bytePlainText, StandardCharsets.UTF_8);
    }
}

