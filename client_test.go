package clientmock

import (
	"net/http"
	"sync"
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

func TestClientMock_Add(t *testing.T) {
	client, mock, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectMethod(http.MethodGet).Add().
		ExpectMethod(http.MethodHead)

	client.Get("http://localhost")
	client.Head("http://localhost")

	if err := mock.ExpectationsMet(); err != nil {
		t.Fatal(err)
	}
}

// Test concurrent requests
func TestClientMock_AddConcurrent(t *testing.T) {
	client, mock, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}

	mock.ExpectMethod(http.MethodGet).Add().
		ExpectMethod(http.MethodHead)

	wg.Add(2)

	go func() {
		client.Get("http://localhost")
		wg.Done()
	}()
	go func() {
		client.Head("http://localhost")
		wg.Done()
	}()

	wg.Wait()
	if err := mock.ExpectationsMet(); err != nil {
		t.Fatal(err)
	}
}

func TestClientMock_AddFail(t *testing.T) {
	client, mock, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectMethod(http.MethodGet).Add().
		ExpectMethod(http.MethodPost)

	client.Get("http://localhost")
	client.Head("http://localhost")

	if err := mock.ExpectationsMet(); err == nil {
		t.Fatal("expected error got nil")
	}
}