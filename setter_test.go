package clientmock

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
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

func TestClientMock_ReturnHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
	}{
		{
			"empty headers",
			map[string]string{},
		},
		{
			"basic header",
			map[string]string{"Authorization": "Bearer"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, mock, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			returnHeader := http.Header{}
			for key, val := range test.headers {
				returnHeader.Add(key, val)
			}

			mock.ReturnHeader(returnHeader)

			res, err := client.Get("http://www.test.com")
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			for key := range returnHeader {
				exp := returnHeader.Get(key)
				got := res.Header.Get(key)
				if exp != "" && exp != got {
					t.Errorf("Expected %s got %s for %s\n", exp, got, key)
				}
			}
		})
	}
}

func TestClientMock_ReadUnsetBody(t *testing.T) {
	client, mock, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectMethod(http.MethodPost)

	res, err := client.Get("http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if data == nil {
		t.Fatal("nil data read from body")
	}
}

func TestClientMock_ReturnError(t *testing.T) {
	client, mock, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	mock.ReturnError(&url.Error{})

	_, err = client.Get("http://test.com")
	if err == nil {
		t.Fatal("expected error got nil")
	}
}