package net.kio.its.logger;

import java.io.IOException;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.text.SimpleDateFormat;
import java.util.Date;

public class Logger {

    private final ILogger iLogger;
    private final FileLogger fileLogger;

    public Logger(ILogger iLogger) {
        this.iLogger = iLogger;
        this.fileLogger = new FileLogger(this);
        fileLogger.initFileLogger();
        System.setErr(new PrintStream(new OutputStream() {

            private boolean errorInProgress = false;

            @Override
            public void write(int b) throws IOException {
                write(new byte[]{(byte) b}, 0, 1);
            }

            @Override
            public void write(byte[] b, int off, int len) throws IOException {
                String string = new String(b, StandardCharsets.UTF_8);
                if (string.startsWith("Exception") || (string.startsWith("java") && !errorInProgress)) {
                    errorInProgress = string.startsWith("Exception");
                    log(LogType.INTERNAL_ERROR, "Exception caught in system error output : ");
                }
                System.out.write(ConsoleColors.RED.getBytes(StandardCharsets.UTF_8));
                System.out.write(b, off, len);
                System.out.write(ConsoleColors.RESET.getBytes(StandardCharsets.UTF_8));
            }
        }));
    }

    public void log(String message) {
        log(LogType.INFO, message);
    }

    public void log(LogType logType, String message) {
        if (logType == LogType.DEBUG && !iLogger.isDebug()) return;
        System.out.println(logType.getPrefix() + getCurrentStringTime() + "  |  " + ConsoleColors.YELLOW_BOLD + iLogger.getName() + ConsoleColors.RESET + logType.getColor() + "  |\t" + message + ConsoleColors.RESET);
        fileLogger.writeLog(logType.getPrefix().substring(logType.getPrefix().indexOf('[')) + getCurrentStringTime() + "  |  " + iLogger.getName() + "  |    " + message);
    }

    private String getCurrentStringTime() {
        return getCurrentStringTime(false);
    }

    String getCurrentStringTime(boolean zipFormat) {
        return new SimpleDateFormat(zipFormat ? "yyyy-MM-dd_HH-mm-ss" : "yyyy-MM-dd HH:mm:ss").format(new Date(System.currentTimeMillis()));
    }

}
