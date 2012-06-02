package nma

import (
	"strings"
	"testing"
)

var verifySampleError = `<?xml version="1.0" encoding="UTF-8"?>` +
	`<nma><error code="400">Parameter 'apikey' not provided.</error></nma>`

var verifySampleSuccess = `<?xml version="1.0" encoding="UTF-8"?>` +
	`<nma><success code="200" remaining="795" resettimer="52"/></nma>`

func TestErrorParsing(t *testing.T) {
	expected := "Parameter 'apikey' not provided."

	_, err := decodeResponse("x", strings.NewReader(verifySampleError))
	if err.Error() != expected {
		t.Fatalf("Expected ``%s'', got ``%s''", expected, err.Error())
	}
}

func TestSuccessParsing(t *testing.T) {
	succ, err := decodeResponse("x", strings.NewReader(verifySampleSuccess))
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
