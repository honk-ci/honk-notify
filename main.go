package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/honk-ci/honk-notify/pkg/github"
	"github.com/honk-ci/honk-notify/pkg/honk"
	"github.com/honk-ci/honk-notify/pkg/twitter"
)

func main() {
	c := make(chan string)

	go twitter.WatchTwitter([]string{
		"kubernetes honk",
		"kubecon honk",
		"kcsna2019",
		"kubekhan",
	}, c)
	go github.WatchGithub([]string{
		"kubernetes/kubernetes",
		"kubernetes/test-infra",
		"kubernetes/release",
		"kubernetes/sig-release",
		"kubernetes/community",
		"kubernetes-sigs/contributor-playground",
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
