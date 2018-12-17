package main

import (
	nt "github.com/nats-io/go-nats"
	"github.com/natsflow/slack-nats/pkg/nats"
	"github.com/natsflow/slack-nats/pkg/slack"
	"os"
)

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")

	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nt.DefaultURL
	}
	n := nats.NewConnection(natsURL)
	defer n.Close()

	go slack.ReqHandler(n.Conn, slackToken)
	go slack.EventStream(n, slackToken)

	select {}
}