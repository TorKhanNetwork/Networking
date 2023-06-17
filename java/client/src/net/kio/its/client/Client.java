package net.kio.its.client;

import net.kio.its.event.EventsManager;
import net.kio.its.logger.ILogger;
import net.kio.its.logger.LogType;
import net.kio.its.logger.Logger;
import org.bouncycastle.jce.provider.BouncyCastleProvider;

import java.io.IOException;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.security.Security;
import java.text.SimpleDateFormat;
import java.util.*;

public class Client implements ILogger {

    private static final String version = "1.0";
    private final String name;
    private static String macAddress;
    private final boolean debug;

    private final Logger logger;
    private final EventsManager eventsManager;
    private final List<SocketWorker> socketWorkerList;


    public Client(String name, boolean debug) {
        this.name = name;
        this.debug = debug;
        this.logger = new Logger(this);
        setMacAddress();
        eventsManager = new EventsManager();
        socketWorkerList = new ArrayList<>();

        // CRYPTOGRAPHIC PROVIDER
        Security.addProvider(new BouncyCastleProvider());
    }

    public Client(String name) {
        this(name, false);
    }

    public static void main(String[] args) {
        Client client = new Client("Main Client", true);

        // Socket Workers
        try {
            SocketWorker socketWorker = client.addSocketWorker("localhost", 40000);
            socketWorker.startWorker();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public void setMacAddress() {
        try {

            String firstInterface = null;
            Map<String, String> addressByNetwork = new HashMap<>();
            Enumeration<NetworkInterface> networkInterfaces = NetworkInterface.getNetworkInterfaces();

            while (networkInterfaces.hasMoreElements()) {
                NetworkInterface network = networkInterfaces.nextElement();

                byte[] bmac = network.getHardwareAddress();
                if (bmac != null) {
                    StringBuilder sb = new StringBuilder();
                    for (int i = 0; i < bmac.length; i++) {
                        sb.append(String.format("%02X%s", bmac[i], (i < bmac.length - 1) ? "-" : ""));
                    }

                    if (!sb.toString().isEmpty()) {
                        addressByNetwork.put(network.getName(), sb.toString());
                    }

                    if (!sb.toString().isEmpty() && firstInterface == null) {
                        firstInterface = network.getName();
                    }
                }
            }

            if (firstInterface != null) {
                macAddress = addressByNetwork.get(firstInterface);
            } else {
                macAddress = randomMACAddress();
                logger.log(LogType.INTERNAL_ERROR, "Unable to find any hardware address on this device, using random bytes : " + macAddress);
            }

        } catch (SocketException e) {
            e.printStackTrace();
        }

    }

    private String randomMACAddress() {
        Random rand = new Random();
        byte[] macAddr = new byte[6];
        rand.nextBytes(macAddr);

        macAddr[0] = (byte) (macAddr[0] & (byte) 254);

        StringBuilder sb = new StringBuilder(18);
        for (byte b : macAddr) {

            if (sb.length() > 0)
                sb.append(":");

            sb.append(String.format("%02x", b));
        }


        return sb.toString();
    }

    public static String getMacAddress() {
        return macAddress;
    }

    public static String getVersion() {
        return version;
    }

    public EventsManager getEventsManager() {
        return eventsManager;
    }

    public SocketWorker addSocketWorker(String ip, int port) throws IOException {
        SocketWorker socketWorker = new SocketWorker(this, ip, port);
        socketWorkerList.add(socketWorker);
        return socketWorker;
    }

    public List<SocketWorker> getSocketWorkerList() {
        return socketWorkerList;
    }

    public static String getCurrentStringTime() {
        return new SimpleDateFormat("yyyy-MM-dd HH:mm:ss").format(new Date(System.currentTimeMillis()));
    }

    public Logger getLogger() {
        return logger;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public boolean isDebug() {
        return debug;
    }
}
