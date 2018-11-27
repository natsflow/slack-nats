package chat

import (
	"errors"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPostHandler(t *testing.T) {
	req := PostMessageReq{
		Text: "Hello slack",
		PostMessageParameters: slack.PostMessageParameters{
			Channel: "CDNPXK2KT",
			AsUser:  true,
		},
	}
	s := new(slackMessageMock)
	s.On("PostMessage", req.Channel, req.Text, req.PostMessageParameters).
		Return("CDNPXK2KT", "1541493485.000300", nil)
	expectedResp := PostMessageResp{
		Channel:   "CDNPXK2KT",
		Timestamp: "1541493485.000300",
	}
	n := new(natsPubMock)
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	postMessageHandler(n, s)("slack.chat.postMessage", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

func TestPostHandlerShouldPublishErrorsIfPostFails(t *testing.T) {
	req := PostMessageReq{
		Text: "Hello slack",
		PostMessageParameters: slack.PostMessageParameters{
			Channel: "CDNPXK2KT",
			AsUser:  true,
		},
	}
	s := new(slackMessageMock)
	err := errors.New("failed to post chat")
	s.On("PostMessage", req.Channel, req.Text, req.PostMessageParameters).
		Return("", "", err)
	expectedResp := PostMessageResp{
		Err: "failed to post chat",
	}
	n := new(natsPubMock)
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	postMessageHandler(n, s)("slack.chat.postMessage", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

func TestMessagesShouldAlwaysBePostedAsUser(t *testing.T) {
	req := PostMessageReq{
		Text: "Hello slack",
		PostMessageParameters: slack.PostMessageParameters{
			Channel: "CDNPXK2KT",
			// this isn't allowed
			AsUser: false,
		},
	}
	expectedPostMessageParameters := slack.PostMessageParameters{
		Channel: "CDNPXK2KT",
		// should get changed to this
		AsUser: true,
	}
	s := new(slackMessageMock)
	s.On("PostMessage", req.Channel, req.Text, expectedPostMessageParameters).
		Return("CDNPXK2KT", "1541493485.000300", nil)
	expectedResp := PostMessageResp{
		Channel:   "CDNPXK2KT",
		Timestamp: "1541493485.000300",
	}
	n := new(natsPubMock)
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	postMessageHandler(n, s)("slack.chat.postMessage", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

type natsPubMock struct {
	mock.Mock
}

func (n *natsPubMock) Publish(subject string, v interface{}) error {
	args := n.Called(subject, v)
	return args.Error(0)
}

type slackMessageMock struct {
	mock.Mock
}

func (s *slackMessageMock) PostMessage(channel, text string, params slack.PostMessageParameters) (respChannel string, respTimestamp string, err error) {
	args := s.Called(channel, text, params)
	return args.String(0), args.String(1), args.Error(2)
}
