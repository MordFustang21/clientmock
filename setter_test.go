package clientmock

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClientMock_ReturnStatus(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		returnedStatus int
		expectedErr    bool
	}{
		{
			"status match",
			http.StatusOK,
			http.StatusOK,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, mock, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			mock.ReturnStatus(test.returnedStatus)

			res, err := client.Get("http://www.test.com")
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.expectedStatus {
				t.Fatalf("response codes don't match expected %d got %d", test.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestClientMock_ReturnBody(t *testing.T) {
	tests := []struct {
		name       string
		returnBody string
	}{
		{
			"good body",
			`{"test": "key"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, mock, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			mock.ReturnBody(bytes.NewBufferString(test.returnBody))

			res, err := client.Get("http://www.test.com")
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			if string(data) != test.returnBody {
				t.Fatalf("expected %s got %s", test.returnBody, string(data))
			}
		})
	}
}
