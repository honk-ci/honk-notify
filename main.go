package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/oauth2"

	"github.com/honk-ci/honk-notify/pkg/github"
	"github.com/honk-ci/honk-notify/pkg/honk"
	"github.com/honk-ci/honk-notify/pkg/twitter"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	c := make(chan string)

	go twitter.WatchTwitter([]string{
		"kubernetes honk",
		"kubecon honk",
		"kcsna2019",
		"kubekhan",
	}, c)
	go github.WatchGithub(tc, []string{
		"kubernetes",
		"kubernetes-sigs",
		"honk-ci",
	}, c)

	go func() {
		for {
			val := <-c
			log.Println(val)
			honk.Honk()
		}
	}()
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Stopping Honk-Alert...")
	close(c)
}
