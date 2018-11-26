package event

import (
	"github.com/natsflow/slack-nats/pkg/nats"
	"github.com/nlopes/slack"
)

// TODO publish other data types as well?
func Handler(n nats.Publisher, sc *slack.RTM) {
	for ev := range sc.IncomingEvents {
		switch ev := ev.Data.(type) {
		case *slack.MessageEvent:
			nats.Publish(n, "slack.event.message", ev)
		}
	}
}
