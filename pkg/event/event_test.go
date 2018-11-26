package event

import (
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHandlerShouldPublishIncomingMessages(t *testing.T) {
	slackRTM := &slack.RTM{}
	slackRTM.IncomingEvents = make(chan slack.RTMEvent)
	msgEvent := &slack.MessageEvent{
		Msg: slack.Msg{
			Channel: "CDNPXK2KT",
			Text:    "blah blah blah",
		},
	}
	ev := slack.RTMEvent{
		Type: "chat",
		Data: msgEvent,
	}
	n := natsPubStub{}
	actualSubject := ""
	actualResp := &slack.MessageEvent{}
	// EventHandler() (and therefore Publish()) will be run in a goroutine, so we need to use a mutex to signal when
	// Publish() has actually been called
	mutex := &sync.Mutex{}
	mutex.Lock()
	n.publish = func(subject string, v interface{}) error {
		actualSubject = subject
		actualResp = v.(*slack.MessageEvent)
		mutex.Unlock()
		return nil
	}

	go Handler(n, slackRTM)
	slackRTM.IncomingEvents <- ev

	// wait until Publish() has been called, otherwise actualSubject and actualResp will not have been set
	mutex.Lock()
	assert.Equal(t, "slack.event.message", actualSubject)
	assert.Equal(t, msgEvent, actualResp)
}

type natsPubStub struct {
	publish func(subject string, v interface{}) error
}

func (n natsPubStub) Publish(subject string, v interface{}) error {
	return n.publish(subject, v)
}