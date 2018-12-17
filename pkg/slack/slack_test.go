package slack

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
			Type:    "message",
			Channel: "CDNPXK2KT",
			Text:    "blah blah blah",
		},
	}
	ev := slack.RTMEvent{
		Type: "message",
		Data: msgEvent,
	}
	// eventHandler() (and therefore Publish()) will be run in a goroutine, so we need to use a mutex to signal when
	// Publish() has actually been called
	mutex := &sync.Mutex{}
	mutex.Lock()
	n := &natsPubStub{mutex: mutex}

	go eventHandler(n, slackRTM)
	slackRTM.IncomingEvents <- ev

	// wait until Publish() has been called, otherwise actualSubject and actualResp will not have been set
	mutex.Lock()
	assert.Equal(t, "slack.event.message", n.actualSubject)
	assert.Equal(t, msgEvent, n.actualResp.(*slack.MessageEvent))
}

type natsPubStub struct {
	actualSubject string
	actualResp    interface{}
	mutex         *sync.Mutex
}

func (n *natsPubStub) Publish(subject string, v interface{}) error {
	n.actualSubject = subject
	n.actualResp = v
	n.mutex.Unlock()
	return nil
}
