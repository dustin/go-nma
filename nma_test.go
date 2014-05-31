package nma

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var verifySampleError = `<?xml version="1.0" encoding="UTF-8"?>` +
	`<nma><error code="400">Parameter 'apikey' not provided.</error></nma>`

var verifySampleSuccess = `<?xml version="1.0" encoding="UTF-8"?>` +
	`<nma><success code="200" remaining="795" resettimer="52"/></nma>`

func TestErrorParsing(t *testing.T) {
	expected := "Parameter 'apikey' not provided."

	_, err := decodeResponse(strings.NewReader(verifySampleError))
	if err == nil || err.Error() != expected {
		t.Errorf("Expected ``%s'', got ``%v''", expected, err)
	}

	_, err = decodeResponse(strings.NewReader("<3"))
	if err == nil {
		t.Errorf("Expected error parsing invalid xml")
	}
}

func hres(status int, s string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Body:       ioutil.NopCloser(strings.NewReader(s)),
	}
}

func TestErrorResponseParsing(t *testing.T) {
	n := &NMA{}
	expected := "Parameter 'apikey' not provided."

	err := n.handleResponse(hres(200, verifySampleError))
	if err == nil || err.Error() != expected {
		t.Errorf("Expected ``%s'', got ``%v''", expected, err)
	}

	expected = "HTTP Error Bad Request - you wrong"
	err = n.handleResponse(hres(400, "you wrong"))
	if err == nil || err.Error() != expected {
		t.Errorf("Expected ``%s'', got ``%v''", expected, err)
	}

	err = n.handleResponse(hres(200, "<3"))
	if err == nil {
		t.Errorf("Expected error parsing invalid xml")
	}
}

func TestSuccessParsing(t *testing.T) {
	succ, err := decodeResponse(strings.NewReader(verifySampleSuccess))
	if err != nil {
		t.Fatalf("Failed to decode success: %v", err)
	}
	if succ.Succ.Code != 200 {
		t.Fatalf("Expected status 200, got %v", succ.Succ.Code)
	}
	if succ.Succ.Remaining != 795 {
		t.Fatalf("Expected 795 remaining, got %v", succ.Succ.Remaining)
	}
	if succ.Succ.Resettimer != 52 {
		t.Fatalf("Expected 52 to reset, got %v", succ.Succ.Resettimer)
	}
}

func TestSetup(t *testing.T) {
	n := New("aqui")
	n.AddKey("k2")
	n.SetDeveloperKey("dk")
	if n.apiKey[0] != "aqui" || n.apiKey[1] != "k2" {
		t.Errorf("Incorrect API keys: %v", n.apiKey)
	}
	if n.developerKey != "dk" {
		t.Errorf("Incorrect developer key: %v", n.developerKey)
	}
	if n.client != http.DefaultClient {
		t.Errorf("Incorrect client: %v", http.DefaultClient)
	}
}

type staticRoundTripper struct {
	res *http.Response
}

func (s staticRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return s.res, nil
}

func testClient(res *http.Response) *http.Client {
	return &http.Client{Transport: staticRoundTripper{res}}
}

type failingRoundTripper struct{}

func (failingRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("nope")
}

func TestVerify(t *testing.T) {
	n := New("aqui")
	n.SetDeveloperKey("dk")

	// Success case
	n.client = testClient(hres(200, verifySampleSuccess))
	if err := n.Verify("k"); err != nil {
		t.Errorf("Failed verifying: %v", err)
	}

	// Error case
	n.client = testClient(hres(400, verifySampleError))
	if err := n.Verify("k"); err == nil {
		t.Errorf("Expected error, but passed :(")
	}

	// Hard error case
	n.client = &http.Client{Transport: failingRoundTripper{}}
	if err := n.Verify("k"); err == nil {
		t.Errorf("Expected error, but passed :(")
	}
}

func TestNotify(t *testing.T) {
	n := New("aqui")
	n.SetDeveloperKey("dk")

	// Success case
	n.client = testClient(hres(200, verifySampleSuccess))
	if err := n.Notify(&Notification{
		Application: "test",
		Description: "A thing",
		Event:       "tested",
		Priority:    High,
		URL:         "http://www.spy.net/",
		ContentType: ContentTypeText,
	}); err != nil {
		t.Errorf("Failed a notification: %v", err)
	}

	// Soft failure case
	n.client = testClient(hres(400, verifySampleError))
	if err := n.Notify(&Notification{}); err == nil {
		t.Errorf("Expected to fail a notification, but succeed")
	}

	// Hard failure case
	n.client = &http.Client{Transport: failingRoundTripper{}}
	if err := n.Notify(&Notification{
		Application: "test",
		Description: "A thing",
		Event:       "tested",
		Priority:    High,
		URL:         "http://www.spy.net/",
		ContentType: ContentTypeText,
	}); err == nil {
		t.Errorf("Expected to fail a notification, but succeed")
	}
}
