package channel

import (
	"github.com/natsflow/slack-nats/pkg/nats"
	"github.com/nlopes/slack"
)

type SlackChannel interface {
	// these signatures are fixed - they are implemented by the 3rd party slack library
	GetChannelInfo(channelID string) (*slack.Channel, error)
	JoinChannel(channelName string) (*slack.Channel, error)
	LeaveChannel(channelID string) (bool, error)
}

// https://api.slack.com/methods/channels.join
// Provide channel id OR name.
func JoinHandler(n nats.PubSub, slack SlackChannel) {
	nats.Subscribe(n, "slack.channel.join", joinHandler(n, slack))
}

func joinHandler(n nats.Publisher, slack SlackChannel) func(subject, reply string, req JoinReq) {
	return func(subject, reply string, req JoinReq) {
		chName := req.Name
		if chName == "" {
			ch, err := slack.GetChannelInfo(req.Id)
			if err != nil {
				nats.PublishReply(n, subject, reply, JoinResp{Err: err.Error()})
				return
			}
			chName = ch.Name
		}
		ch, err := slack.JoinChannel(chName)
		resp := JoinResp{
			Channel: ch,
		}
		if err != nil {
			resp.Err = err.Error()
		}
		nats.PublishReply(n, subject, reply, resp)
	}
}

type JoinReq struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type JoinResp struct {
	Channel *slack.Channel `json:"channel"`
	Err     string         `json:"err"`
}

// https://api.slack.com/methods/channels.leave
func LeaveHandler(n nats.PubSub, slack SlackChannel) {
	nats.Subscribe(n, "slack.channel.leave", leaveHandler(n, slack))
}

func leaveHandler(n nats.Publisher, slack SlackChannel) func(subject, reply string, req LeaveReq) {
	return func(subject, reply string, req LeaveReq) {
		notInChannel, err := slack.LeaveChannel(req.Id)
		resp := LeaveResp{
			NotInChannel: notInChannel,
		}
		if err != nil {
			resp.Err = err.Error()
		}
		nats.PublishReply(n, subject, reply, resp)
	}
}

type LeaveReq struct {
	Id string `json:"id"`
}

type LeaveResp struct {
	NotInChannel bool   `json:"not_in_channel"`
	Err          string `json:"err"`
}
