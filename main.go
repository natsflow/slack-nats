package main

import (
	"github.com/nats-io/go-nats"
	"github.com/natsflow/slack-nats/pkg/slack"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")

	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}
	n := newNatsConn(natsURL)
	defer n.Close()

	go slack.EventStream(n, slackToken)
	go slack.ReqHandler(n.Conn, slackToken)

	select {}
}

func newNatsConn(url string) *nats.EncodedConn {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("url", url).
			Msg("Failed to connect to NATS")
	}
	log.Info().Str("url", url).Msg("Connected to NATS")

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create NATS json connection")
	}
	return ec
}
