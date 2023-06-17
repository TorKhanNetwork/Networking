package net.kio.its.server.events;

import net.kio.its.server.ServerWorker;

import java.util.UUID;

public class CommandReceivedEvent extends DataEvent {

    private final String command;
    private final String[] args;

    public CommandReceivedEvent(ServerWorker serverWorker, String data, String prefix, String splitter, UUID messageUUID) {
        super(serverWorker, data, messageUUID);
        String[] splitted = data.substring(prefix.length()).split(splitter);
        this.command = splitted[0];
        this.args = (splitted.length >= 2) ? data.substring(prefix.length() + splitter.length() + command.length()).split(splitter) : new String[]{};
    }

    public String getCommand() {
        return command;
    }

    public String[] getArgs() {
        return args;
    }

    @Override
    public String getData() {
        return super.getData();
    }
}
