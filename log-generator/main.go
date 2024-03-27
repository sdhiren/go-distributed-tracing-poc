package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	logFilePath = "../app/app4.log"
	maxFileSize = 1 << 30 // 1 GB
	logEntry    = `{"time":"%s","level":"INFO","msg":"testing splunk retention policy ","trace_id":"ad29519f158159a4ebd15ca6813fskjd74","span_id":"e03c5defe36d9b24","method_name":"CallApi2","class_name":"ApiController1"}`
)

func main() {
	file, err := os.Create(logFilePath)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer file.Close()

	start := time.Now()

	for {

		loc, _ := time.LoadLocation("Asia/Kolkata")
		entry := fmt.Sprintf(logEntry, time.Now().In(loc).Format(time.RFC3339Nano))
	
		_, err = file.WriteString(entry + "\n")
		if err != nil {
			log.Fatalf("Failed to write log entry to file: %v", err)
		}

		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatalf("Failed to get file info: %v", err)
		}
		if fileInfo.Size() >= maxFileSize {
			break
		}
	}

	duration := time.Since(start)
	fmt.Printf("Time taken: %v\n", duration)
}
