package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func configure() {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig("consumer-key", "consumer-secret")
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken("access-token", "access-token-secret")
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
	}

	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"cat"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()

}

func main() {
	fmt.Println("Go-Twitter Bot v0.01")
	configure()
}
