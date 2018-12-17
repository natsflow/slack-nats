package slack

import (
	"bytes"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var logger *zap.SugaredLogger

func init() {
	l, _ := zap.NewProduction()
	logger = l.Sugar()
}

func Client(token string) *slack.RTM {
	s := slack.New(token).NewRTM()
	if _, err := s.AuthTest(); err != nil {
		logger.Fatalf("failed to authenticate with slack: %v", err)
	}
	go s.ManageConnection()
	return s
}


func NewSlack() {

}

type Slack struct {

}



//type NatsPublisher interface {
//	Publish(subject string, v interface{}) error
//}

// TODO publish other data types as well?
func EventStream(n *nats.EncodedConn, slackToken string) {
	s := Client(slackToken)

	for ev := range s.IncomingEvents {
		switch d := ev.Data.(type) {
		case *slack.MessageEvent:
			// e.g. slack.event.message
			subject := fmt.Sprintf("slack.event.%s", ev.Type)
			if err := n.Publish(subject, d); err != nil {
				logger.Errorf("could not publish to nats subject=%s: %v", subject, err)
			}
		}
	}
}

func ReqHandler(n *nats.Conn, slackToken string) {
	subj := "slack.>"
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	if _, err := n.Subscribe(subj, func(m *nats.Msg) {
		logger.Infof("got subject=%s with reply %s", m.Subject, m.Reply)

		if strings.HasPrefix(m.Subject, "slack.event.") {
			return // these are events we've raised & aren't requests
		}

		respMsg := reqHandler(c, slackToken, m.Subject, m.Data)
		if err := n.Publish(m.Reply, respMsg); err != nil {
			logger.Errorf("could not publish to nats subject=%s reply=%s: %v", m.Subject, m.Reply, err)
		}
	}); err != nil {
		logger.Fatalf("failed to subscribe to %q: %v", subj, err)
	}
}


func reqHandler(c *http.Client, slackToken string, subj string, reqMsg []byte) []byte {
	u := subjectToUrl(subj)

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(reqMsg))
	if err != nil {
		return errorResp(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slackToken))
	// TODO set on request context?
	//req.WithContext()

	resp, err := c.Do(req)
	if err != nil {
		return errorResp(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errorResp(err)
	}
	return body
}



//func testConnection() {
//
//}

// converts e.g `slack.channels.leave` to `https://slack.com/api/channels.leave`
func subjectToUrl(subj string) string {
	path := strings.TrimPrefix(subj, "slack.")
	return fmt.Sprintf("https://slack.com/api/%s", path)
}

func errorResp(err error) []byte {
	return []byte(fmt.Sprintf(`{"error" : "%s"}`, err.Error()))
}
