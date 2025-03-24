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
	transitionWaitID := os.Getenv("transitionWaitID")
	queues := client.GetQueues(os.Getenv("queueUrl"))
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-sigch:
			break loop
		case <-ticker.C:
			tasks := task(queues)
			for _, i := range tasks {
				t := client.GetTask(i)

				//orderFile := client.GetTaskDescription(t.Self)
				for _, a := range t.Fields.Attachment {
					orderFile := client.GetAttachment(a.FileName, a.Content, t.Key)
					resp := client.loadOrder(orderFile)
					fmt.Println(resp)
				}
				client.setStatus(t.Self, transitionWaitID)
				client.addComment(t.Self)

				//token := client.GetLoaderToken()
				//fmt.Println("Token:", token)

			}
		}
	}
	//https://jira.symboltransport.com/rest/api/2/issue/PER-3
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
