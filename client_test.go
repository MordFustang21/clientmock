package clientmock

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client, mock, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal(err)
	}

	if mock == nil {
		t.Fatal(err)
	}
}
