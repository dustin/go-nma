package nma

import (
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
