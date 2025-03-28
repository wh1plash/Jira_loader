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

	var (
		transitionWaitID  = os.Getenv("TRANSITION_WAIT_ID")
		transitionDoneID  = os.Getenv("TRANSITION_DONE_ID")
		ticker            = time.NewTicker(getTickerInterval())
		waitStatusComment = os.Getenv("WAITSTATUS_COMMENT")
	)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-sigch:
			break loop
		case <-ticker.C:
			queues := client.GetQueues(os.Getenv("QUEUE_URL"))
			tasks := task(queues)
			for _, i := range tasks {
				t := client.GetTask(i)
				client.setStatus(t.Self, transitionWaitID)
				client.addComment(t.Self, waitStatusComment)
				//orderFile := client.GetTaskDescription(t.Self)
				for _, a := range t.Fields.Attachment {
					orderFile := client.GetAttachment(a.FileName, a.Content, t.Key)
					result := client.loadOrder(orderFile)
					if err := client.addCommentWithResult(t.Self, result); err != nil {
						fmt.Println("error to add comment of result of load order to Jira issue:", err)
					}
				}
				client.setStatus(t.Self, transitionDoneID)

			}
		}
	}
}

func main() {
	baseurl := os.Getenv("BASE_URL")
	client := NewHTTPClient(baseurl)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)

	go runFetcher(client, sigch)

	<-sigch
	fmt.Println("Exit programm")
}
