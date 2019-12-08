package twitter

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/coreos/pkg/flagutil"
)

// WatchTwitter takes a channel, and honks
func WatchTwitter(topics []string, c chan interface{}) {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	api := anaconda.NewTwitterApiWithCredentials(*accessToken, *accessSecret, *consumerKey, *consumerSecret)
	streamValues := url.Values{}
	streamValues.Set("track", strings.Join(topics, ","))
	streamValues.Set("stall_warnings", "true")
	log.Println("Starting Stream...")
	s := api.PublicStreamFilter(streamValues)

	go func() {
		for t := range s.C {
			switch v := t.(type) {
			case anaconda.Tweet:
				c <- v
			}
			time.Sleep(1)
		}
	}()

	log.Println("Now watching Twitter...")

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Stopping Twitter Watch...")
}
