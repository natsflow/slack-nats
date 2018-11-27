package event

import (
	"fmt"
	"github.com/natsflow/slack-nats/pkg/nats"
	"github.com/nlopes/slack"
)

// TODO publish other data types as well?
func Handler(n nats.Publisher, sc *slack.RTM) {
	for ev := range sc.IncomingEvents {
		switch ev := ev.Data.(type) {
		case *slack.MessageEvent:
			// e.g. slack.event.message
			nats.Publish(n, fmt.Sprintf("slack.event.%s", ev.Type), ev)
		}
	}
}
