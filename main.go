package main

import (
	"github.com/nats-io/go-nats"
	"github.com/natsflow/slack-nats/pkg/slack"
	"go.uber.org/zap"
	"os"
)

var logger *zap.SugaredLogger

func init() {
	l, _ := zap.NewProduction()
	logger = l.Sugar()
}

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
		logger.Fatalf("Failed to connect to nats on %q: %v", url, err)
	}
	logger.Infof("Connected to nats %s", url)

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		logger.Fatalf("Failed to create nats json connection: %v", err)
	}
	return ec
}
