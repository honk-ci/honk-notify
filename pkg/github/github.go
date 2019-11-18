package github

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/github"
)

// FetchComments looks at a specific org and returns a list of a bunch of comments
func FetchComments(client *http.Client, org string) []*github.Event {
	ghClient := github.NewClient(client)
	listOpts := github.ListOptions{PerPage: 50}
	events, _, _ := ghClient.Activity.ListEventsForOrganization(context.Background(), org, &listOpts)
	return events
}

// WatchGithub looks at a list of orgs and watches for a honking comment and triggers an event
func WatchGithub(client *http.Client, orgs []string, c chan string) {
	latestComment := make(map[string]*github.Event)
	go func() {
		log.Println("Now watching GitHub...")
		for {
			for _, org := range orgs {
				events := FetchComments(client, org)
				for _, event := range events {
					if event.GetType() == "IssueCommentEvent" && event.GetID() > latestComment[org].GetID() {
						latestComment[org] = event
						comment, _ := event.ParsePayload()
						body := comment.(*github.IssueCommentEvent).GetComment().GetBody()
						if strings.Contains(body, "/honk") {
							c <- "GitHub: " + body
						}

					}
				}
			}
			time.Sleep(20 * time.Second)
		}
	}()
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Stopping GitHub Watch...")
	close(c)
}
