package net.kio.its.server;

import net.kio.its.event.EventsManager;
import net.kio.its.logger.ILogger;
import net.kio.its.logger.Logger;
import net.kio.its.server.events.ClientSocketOpenedEvent;
import net.kio.its.server.events.ServerWorkerBoundEvent;
import net.kio.security.dataencryption.KeysGenerator;
import org.bouncycastle.jce.provider.BouncyCastleProvider;

import java.io.IOException;
import java.net.Socket;
import java.security.Security;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Random;

public class Server implements ILogger {

    private final String name;

    private static final String version = "1.0";
    private final SocketThreadListener socketThreadListener;
    private final int port;

    private final Logger logger;
    private final List<ServerWorker> clients;
    private final EventsManager eventsManager;

    private final KeysGenerator keysGenerator;

    private final boolean debug;

    public Server(String name, int port, boolean debug) throws IOException {
        this.name = name;
        this.debug = debug;
        this.port = port;
        this.logger = new Logger(this);
        this.socketThreadListener = new SocketThreadListener(this);
        this.clients = new ArrayList<>();
        this.eventsManager = new EventsManager();
        this.keysGenerator = new KeysGenerator();

        // CRYPTOGRAPHIC PROVIDER
        Security.addProvider(new BouncyCastleProvider());
        keysGenerator.generateKeys(false, true);
    }

    public Server(String name, int port) throws IOException {
        this(name, port, false);
    }

    public static void main(String[] args) throws IOException {

        Server server = new Server("Main Test Server", 40000, true);
        server.startSocketListener();

    }

    public void startSocketListener() {
        if (!socketThreadListener.isAlive()) socketThreadListener.start();
    }

    public void handleSocketConnection(Socket socket) throws IOException {
        ClientSocketOpenedEvent event = new ClientSocketOpenedEvent(socket);
        eventsManager.callEvent(event);
        if (!event.isCancelled()) {
            ServerWorker serverWorker = new ServerWorker(this, socket, keysGenerator);
            ServerWorkerBoundEvent serverWorkerBoundEvent = new ServerWorkerBoundEvent(serverWorker);
            eventsManager.callEvent(serverWorkerBoundEvent);
            if (!serverWorkerBoundEvent.isCancelled()) {
                clients.add(serverWorker);
                serverWorker.start();
            }
        } else {
            socket.close();
        }
    }

    public static String generateRequestSeparator() {
        int leftLimit = 48; // numeral '0'
        int rightLimit = 122; // letter 'z'
        int targetStringLength = 20;
        Random random = new Random();

        return random.ints(leftLimit, rightLimit + 1)
                .filter(i -> (i <= 57 || i >= 65) && (i <= 90 || i >= 97))
                .limit(targetStringLength)
                .collect(StringBuilder::new, StringBuilder::appendCodePoint, StringBuilder::append)
                .toString();
    }

    public static String getVersion() {
        return version;
    }

    public int getPort() {
        return port;
    }

    public Logger getLogger() {
        return logger;
    }

    public SocketThreadListener getSocketThreadListener() {
        return socketThreadListener;
    }

    public List<ServerWorker> getClients() {
        return clients;
    }

    public EventsManager getEventsManager() {
        return eventsManager;
    }

    public static String getCurrentStringTime() {
        return new SimpleDateFormat("yyyy-MM-dd HH:mm:ss").format(new Date(System.currentTimeMillis()));
    }

    @Override
    public boolean isDebug() {
        return debug;
    }

    @Override
    public String getName() {
        return name;
    }
}
