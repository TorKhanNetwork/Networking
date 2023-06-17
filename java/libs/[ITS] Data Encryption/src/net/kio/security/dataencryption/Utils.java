package net.kio.security.dataencryption;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

public class Utils {
    public Utils() {
    }

    public static String hashPassword(String password) {
        StringBuilder hashText = null;

        try {
            MessageDigest md = MessageDigest.getInstance("MD5");
            byte[] messageDigest = md.digest(password.getBytes(StandardCharsets.UTF_8));
            BigInteger no = new BigInteger(1, messageDigest);
            hashText = new StringBuilder(no.toString(16));

            while (hashText.length() < 32) {
                hashText.insert(0, "0");
            }
        } catch (NoSuchAlgorithmException var5) {
            var5.printStackTrace();
        }

        return hashText != null ? hashText.toString() : null;
    }
}
