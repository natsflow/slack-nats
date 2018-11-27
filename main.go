package main

import (
	nt "github.com/nats-io/go-nats"
	"github.com/natsflow/slack-nats/pkg/channel"
	"github.com/natsflow/slack-nats/pkg/chat"
	"github.com/natsflow/slack-nats/pkg/event"
	"github.com/natsflow/slack-nats/pkg/nats"
	"github.com/nlopes/slack"
	"os"
)

func main() {
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nt.DefaultURL
	}
	n := nats.NewConnection(natsURL)
	defer n.Close()

	s := slack.New(os.Getenv("SLACK_TOKEN")).NewRTM()
	go s.ManageConnection()

	go chat.PostMessageHandler(n, s)
	go event.Handler(n, s)
	go channel.JoinHandler(n, s)
	go channel.LeaveHandler(n, s)

	select {}
}
