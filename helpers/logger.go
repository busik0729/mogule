package helpers

import (
	"log"
	"os"
)

func LogToFile(message string) {
	logFilePath := "/home/busik/web/crm/log.txt"
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		file, _ := os.Create(logFilePath)
		defer file.Close()
		f, _ = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "prefix", log.LstdFlags)

	logger.Println(message)
}

func DaemonLogToFile(message string) {
	logFilePath := "/home/busik/web/crm/daemon_log.txt"
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		file, _ := os.Create(logFilePath)
		defer file.Close()
		f, _ = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "prefix", log.LstdFlags)

	logger.Println(message)
}
