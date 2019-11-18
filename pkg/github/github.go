package github

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/github"
)

// FetchComments looks at a specific repo and returns a list of a bunch of comments
func FetchComments(repo string) []*github.RepositoryComment {
	s := strings.Split(repo, "/")
	client := github.NewClient(nil)
	listOpts := github.ListOptions{}
	comments, _, _ := client.Repositories.ListComments(context.Background(), s[0], s[1], &listOpts)
	return comments
}

// WatchGithub looks at a list of repos and watches for a honking comment and triggers an event
func WatchGithub(repos []string, c chan string) {
	latestComment := make(map[string]*github.RepositoryComment)
	go func() {
		for {
			for _, repo := range repos {
				comments := FetchComments(repo)
				for _, comment := range comments {
					if comment.GetID() > latestComment[repo].GetID() {
						latestComment[repo] = comment
						if strings.Contains(comment.GetBody(), "/honk") {
							c <- "GitHub: " + comment.GetBody()
						}
					}
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Stopping GitHub Watch...")
	close(c)
}
