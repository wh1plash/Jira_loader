package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

func seveOrderFile(b []byte, fileName string, taskID string) (string, error) {
	currentDate := time.Now().Format("2006-01-02")
	dir := filepath.Join("./downloads", currentDate, taskID)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal("error to create folder to download order", err)
	}

	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("error to create file:", err)
	}
	defer file.Close()

	if _, err := file.Write(b); err != nil {
		log.Fatal("can`t write file:", err)
	}

	return filePath, nil
}

func task(q Queue) []string {
	var tasksUrl []string
	for _, queue := range q.Values {
		if queue.Status.Name == "Open" {
			tasksUrl = append(tasksUrl, queue.Self)
		}
	}
	return tasksUrl
}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
