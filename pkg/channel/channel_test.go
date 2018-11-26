package channel

import (
	"errors"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestJoinUsingChannelName(t *testing.T) {
	req := JoinReq{
		Name: "hcom-nats-test",
	}
	ch := &slack.Channel{
		IsChannel: true,
		IsMember:  true,
	}
	s := new(slackChannelMock)
	s.On("JoinChannel", req.Name).Return(ch, nil)

	expectedResp := JoinResp{
		Channel: ch,
	}
	n := new(natsPubMock)
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	joinHandler(n, s)("slack.channel.join", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

func TestJoinHandlerShouldPublishErrorIfJoinFails(t *testing.T) {
	req := JoinReq{
		Name: "hcom-nats-test",
	}
	ch := &slack.Channel{}
	err := errors.New("failed to join channel")
	s := new(slackChannelMock)
	s.On("JoinChannel", req.Name).Return(ch, err)
	expectedResp := JoinResp{
		Channel: ch,
		Err:     "failed to join channel",
	}
	n := new(natsPubMock)
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	joinHandler(n, s)("slack.channel.join", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

func TestJoinUsingChannelId(t *testing.T) {
	req := JoinReq{
		Id: "CDNPXK2KT",
	}
	ch := &slack.Channel{
		IsChannel: true,
		IsMember:  true,
	}
	ch.Name = "hcom-nats-test"
	s := new(slackChannelMock)
	s.On("GetChannelInfo", req.Id).Return(ch, nil)
	s.On("JoinChannel", "hcom-nats-test").Return(ch, nil)
	expectedResp := JoinResp{
		Channel: ch,
	}
	n := new(natsPubMock)
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	joinHandler(n, s)("slack.channel.join", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

func TestJoinHandlerShouldPublishErrorsIfIdNotValidAndNameNotProvided(t *testing.T) {
	req := JoinReq{
		Id: "XXXXXXXX",
	}
	ch := &slack.Channel{}
	err := errors.New("id not a valid channel")
	s := new(slackChannelMock)
	s.On("GetChannelInfo", req.Id).Return(ch, err)
	expectedResp := JoinResp{
		Err: err.Error(),
	}
	n := new(natsPubMock)
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	joinHandler(n, s)("slack.channel.join", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

func TestLeave(t *testing.T) {
	req := LeaveReq{
		Id: "CDNPXK2KT",
	}
	s := new(slackChannelMock)
	s.On("LeaveChannel", req.Id).Return(false, nil)
	expectedResp := LeaveResp{}
	n := new(natsPubMock)
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	leaveHandler(n, s)("slack.channel.leave", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

func TestLeaveWhenUserWasNotInChannel(t *testing.T) {
	req := LeaveReq{
		Id: "CDNPXK2KT",
	}
	s := new(slackChannelMock)
	// user was not in channel
	s.On("LeaveChannel", req.Id).Return(true, nil)
	n := new(natsPubMock)
	expectedResp := LeaveResp{
		NotInChannel: true,
	}
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	leaveHandler(n, s)("slack.channel.leave", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

	s.AssertExpectations(t)
	n.AssertExpectations(t)
}

func TestLeaveHandlerShouldPublishErrorIfLeaveFails(t *testing.T) {
	req := LeaveReq{
		Id: "CDNPXK2KT",
	}
	err := errors.New("failed to leave channel")
	s := new(slackChannelMock)
	s.On("LeaveChannel", req.Id).Return(false, err)
	n := new(natsPubMock)
	// error returned in resp
	expectedResp := LeaveResp{Err: "failed to leave channel"}
	n.On("Publish", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", expectedResp).Return(nil)

	leaveHandler(n, s)("slack.channel.leave", "_INBOX.OPYB6GUMJ4FYBAYWTKB6WP.OPYB6GUMJ4FYBAYWTKB6ZD", req)

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

type slackChannelMock struct {
	mock.Mock
}

func (s *slackChannelMock) GetChannelInfo(channelID string) (*slack.Channel, error) {
	args := s.Called(channelID)
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (s *slackChannelMock) JoinChannel(channelName string) (*slack.Channel, error) {
	args := s.Called(channelName)
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (s *slackChannelMock) LeaveChannel(channelID string) (bool, error) {
	args := s.Called(channelID)
	return args.Bool(0), args.Error(1)
}
