package chat

import (
	"github.com/natsflow/slack-nats/pkg/nats"
	"github.com/nlopes/slack"
)

type SlackChat interface {
	// these signatures are fixed - they are implemented by the 3rd party slack library
	PostMessage(channel, text string, params slack.PostMessageParameters) (respChannel string, respTimestamp string, err error)
}

// https://api.slack.com/methods/chat.postMessage
func PostMessageHandler(n nats.PubSub, slack SlackChat) {
	nats.Subscribe(n, "slack.chat.postMessage", postMessageHandler(n, slack))
}

func postMessageHandler(n nats.Publisher, slackChat SlackChat) func(subject, reply string, req PostMessageReq) {
	return func(subject, reply string, req PostMessageReq) {
		req.PostMessageParameters.AsUser = true
		respChannel, respTimestamp, err := slackChat.PostMessage(req.Channel, req.Text, req.PostMessageParameters)
		resp := PostMessageResp{
			Channel:   respChannel,
			Timestamp: respTimestamp,
		}
		if err != nil {
			resp.Err = err.Error()
		}
		nats.PublishReply(n, subject, reply, resp)
	}
}

type PostMessageReq struct {
	Text string `json:"text"`
	slack.PostMessageParameters
}

type PostMessageResp struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"`
	Err       string `json:"err"`
}
