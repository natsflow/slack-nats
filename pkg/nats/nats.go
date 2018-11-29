package nats

import (
	"github.com/nats-io/go-nats"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	l, _ := zap.NewProduction()
	logger = l.Sugar()
}

type Publisher interface {
	Publish(subject string, v interface{}) error
}

type Subscriber interface {
	Subscribe(subject string, cb nats.Handler) (*nats.Subscription, error)
}

type PubSub interface {
	Publisher
	Subscriber
}

func NewConnection(url string) *nats.EncodedConn {
	nc, err := nats.Connect(url)
	if err != nil {
		logger.Fatalf("Failed to connect to nats on %q: %v", url, err)
	}
	logger.Infof("Connected to nats %s", url)

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		logger.Fatalf("Failed to get json encoder: %v", err)
	}
	return ec
}

// common pub/sub & log patterns
func Subscribe(n Subscriber, subject string, handler nats.Handler) {
	if _, err := n.Subscribe(subject, handler); err != nil {
		logger.Fatalf("failed to subscribe to subject=%s: %v", subject, err)
	}
	logger.Infof("Subscribed to subject=%s", subject)
}

func PublishReply(n Publisher, subject, reply string, resp interface{}) {
	if err := n.Publish(reply, resp); err != nil {
		logger.Errorf("could not publish to nats subject=%s reply=%s: %v", subject, reply, err)
	}
}

func Publish(n Publisher, subject string, event interface{}) {
	if err := n.Publish(subject, event); err != nil {
		logger.Errorf("could not publish to nats subject=%s: %v", subject, err)
	}
}
