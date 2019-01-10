package slack

import (
	"bytes"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/nlopes/slack"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func EventStream(n *nats.EncodedConn, slackToken string) {
	eventHandler(n, newRtmClient(slackToken))
}

func newRtmClient(token string) *slack.RTM {
	s := slack.New(token).NewRTM()
	if _, err := s.AuthTest(); err != nil {
		log.Fatal().Err(err).Msg("failed to authenticate with slack")
	}
	go s.ManageConnection()
	return s
}

// TODO publish other data types as well?
func eventHandler(n NatsPublisher, sc *slack.RTM) {
	for ev := range sc.IncomingEvents {
		switch d := ev.Data.(type) {
		case *slack.MessageEvent:
			// e.g. slack.event.message
			subject := fmt.Sprintf("slack.event.%s", ev.Type)
			if err := n.Publish(subject, d); err != nil {
				log.Error().Err(err).Str("subject", subject).Msg("could not publish to NATS")
			}
		}
	}
}

type NatsPublisher interface {
	Publish(subject string, v interface{}) error
}

func ReqHandler(n *nats.Conn, slackToken string) {
	c := newSlack(slackToken)
	if _, err := n.QueueSubscribe("slack.>", "slack-nats", func(m *nats.Msg) {
		if strings.HasPrefix(m.Subject, "slack.event.") {
			return // these are events we've raised & aren't requests, so dump them
		}
		respMsg := c.Do(toPath(m.Subject), m.Data)
		if err := n.Publish(m.Reply, respMsg); err != nil {
			log.Error().
				Err(err).
				Str("subject", m.Subject).
				Str("reply", m.Reply).
				Msg("could not publish to NATS")
		}
	}); err != nil {
		log.Fatal().Err(err).Str("subject", "slack.>").Msg("could not subscribe to NATS")
	}
}

// e.g `slack.channels.leave` -> `channels.leave`
func toPath(subj string) string {
	return strings.TrimPrefix(subj, "slack.")
}

type Slack struct {
	client HttpDoer
	token  string
	url    string
}

type HttpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

func newSlack(token string) Slack {
	return Slack{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: token,
		url:   "https://slack.com/api/",
	}
}

func (s Slack) Do(path string, body []byte) []byte {
	req, err := http.NewRequest(http.MethodPost, s.url+path, bytes.NewReader(body))
	if err != nil {
		return errorResp(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	// TODO set timeout on request context?
	//req.WithContext()

	resp, err := s.client.Do(req)
	if err != nil {
		return errorResp(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errorResp(err)
	}
	return respBody
}

func errorResp(err error) []byte {
	return []byte(fmt.Sprintf(`{"error" : "%s"}`, err.Error()))
}
