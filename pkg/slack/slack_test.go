package slack

import (
	"errors"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
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

func TestSlackClient(t *testing.T) {
	resp := []byte(`{ok: true,channel: 'CDNPXK2KT',ts:'1545143890.003100',message:{bot_id:'B9L0ACSUA',type:'message',text:'Hello there',user:'U7KMBRAVB',ts:'1545143890.003100'}}`)
	var actualReq *http.Request
	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualReq = r
		w.Write(resp)
	}))
	defer slackServer.Close()

	s := newSlack("MY_SLACK_TOKEN")
	s.url = slackServer.URL + "/"
	message := `{text: 'Hello there', channel: 'CDNPXK2KT', as_user: 'true'}`
	respBody := s.Do("chat.PostMessage", []byte(message))

	assert.Equal(t, "application/json; charset=utf-8", actualReq.Header.Get("Content-Type"))
	assert.Equal(t, "Bearer MY_SLACK_TOKEN", actualReq.Header.Get("Authorization"))
	assert.Equal(t, resp, respBody)
}

func TestSlackClientShouldReturnErrorResponse(t *testing.T) {
	httpCli := new(HttpClientMock)
	httpCli.On("Do", mock.Anything).Return(&http.Response{}, errors.New("unknown Host"))

	s := newSlack("MY_SLACK_TOKEN")
	s.client = httpCli

	respBody := s.Do("chat.PostMessage", []byte(`{text: 'Hello there', channel: 'CDNPXK2KT', as_user: 'true'}`))

	assert.Equal(t, errorResp(errors.New("unknown Host")), respBody)
}

type HttpClientMock struct {
	mock.Mock
}

func (h *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	args := h.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
