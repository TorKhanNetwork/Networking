package net.kio.its.logger;

import java.io.File;
import java.io.FileOutputStream;
import java.io.FileWriter;
import java.io.IOException;
import java.nio.file.Files;
import java.util.zip.ZipEntry;
import java.util.zip.ZipOutputStream;

public class FileLogger {

    private final Logger logger;

    private File logFile;

    public FileLogger(Logger logger) {
        this.logger = logger;
        Runtime.getRuntime().addShutdownHook(new Thread(this::closeLogFile));
    }

    public void initFileLogger() {
        File dir = new File("./logs/");
        if (!dir.exists()) dir.mkdir();
        logFile = new File("./logs/latest.log");
        try {
            if (!logFile.exists()) {
                logFile.createNewFile();
            }
            new FileWriter(logFile).close();

        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public void writeLog(String line) {
        if (logFile != null && logFile.exists()) {
            try {
                FileWriter fileWriter = new FileWriter(logFile, true);
                fileWriter.write(line + "\n");
                fileWriter.close();
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }

    private void closeLogFile() {
        if (logFile != null && logFile.exists()) {
            File zipFile = new File("./logs/" + logger.getCurrentStringTime(true) + ".zip");
            if (!zipFile.exists()) {
                try {
                    zipFile.createNewFile();
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
            try (
                    FileOutputStream outputStream = new FileOutputStream("./logs/" + logger.getCurrentStringTime(true) + ".zip");
                    ZipOutputStream zipOutputStream = new ZipOutputStream(outputStream);
            ) {
                zipOutputStream.putNextEntry(new ZipEntry(zipFile.getName().replace(".zip", ".log")));
                zipOutputStream.write(Files.readAllBytes(logFile.toPath()));
                zipOutputStream.closeEntry();
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }
}
