package net.kio.its.logger;

public enum LogType {

    INFO(ConsoleColors.GREEN, ConsoleColors.GREEN_BOLD + "[INFO] "),
    DEBUG(ConsoleColors.YELLOW, ConsoleColors.YELLOW_BOLD + "[DEBUG] "),
    CRITICAL(ConsoleColors.PURPLE, ConsoleColors.PURPLE_BOLD + "[CRITICAL] "),
    INTERNAL_ERROR(ConsoleColors.RED, ConsoleColors.RED_BOLD + "[INTERNAL ERROR] ");

    private final String prefix;
    private final String color;

    LogType(String consoleColor, String prefix) {
        this.color = consoleColor;
        this.prefix = prefix;
    }

    public String getPrefix() {
        return prefix + color;
    }

    public String getColor() {
        return color;
    }
}
