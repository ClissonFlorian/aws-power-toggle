package backend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// TestCreateHTTPClient should test the behaviour of createHTTPClient
func TestCreateHTTPClient(t *testing.T) {
	httpClient := createHTTPClient()

	if httpClient.Timeout != 5*time.Second {
		t.Errorf("http client timeout not set")
	}
}

// TestSlackSendMessage tests the slack client behaviour
func TestSlackSendMessage(t *testing.T) {

	slackEnabled = true
	slackWebHooks = []string{"https://slack.linuxctl.com"}

	// replace slackClient to avoid making a real http request
	slackClient = NewTestClient(func(req *http.Request) *http.Response {
		// Check the request params
		if req.Method != "POST" {
			t.Errorf("expected POST but sent %v", req.Method)
		}
		if req.URL.String() != slackWebHooks[0] {
			t.Errorf("slack client sent to unexpected endpoint: %v", req.URL.String())
		}

		// response to slack client
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`ok`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	slackSendMessage("test data")
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
