package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	mustLoadEnvVariables()
}

func runFetcher(client *HTTPClient, sigch chan os.Signal) {
	transitionID := os.Getenv("transitionID")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-sigch:
			fmt.Println("Stopping Jira fetcher")
			break loop
		case <-ticker.C:
			queues := client.GetQueues(os.Getenv("queueUrl"))
			tasks := task(queues)

			for _, i := range tasks {
				t := client.GetTask(i)
				client.setStatus(t.Self, transitionID)
				client.addComment(t.Self)
				for _, a := range t.Fields.Attachment {
					client.GetAttachment(a.FileName, a.Content)
				}
			}
		}
	}
}

func main() {
	baseurl := os.Getenv("baseUrl")
	client := NewHTTPClient(baseurl)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)

	go runFetcher(client, sigch)

	<-sigch
	fmt.Println("Exit programm")
}
