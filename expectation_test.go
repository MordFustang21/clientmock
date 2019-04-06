package clientmock

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func TestClientMock_ExpectMethod(t *testing.T) {

	tests := []struct {
		name           string
		reqMethod      string
		expectedMethod string
		wantErr        bool
	}{
		{
			"expected correct",
			http.MethodGet,
			http.MethodGet,
			false,
		},
		{
			"expected error",
			http.MethodPost,
			http.MethodGet,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, mock, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			mock.ExpectMethod(test.expectedMethod)

			req, err := http.NewRequest(test.reqMethod, "http://test.com", nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			err = mock.ExpectationsMet()
			if (err != nil) != test.wantErr {
				t.Fatalf("error check failed wantError: %t got %v", test.wantErr, err)
			}
		})
	}
}

func TestClientMock_ExpectBody(t *testing.T) {

	tests := []struct {
		name         string
		reqBody      string
		expectedBody string
		wantErr      bool
	}{
		{
			"expected correct",
			"test body",
			"test body",
			false,
		},
		{
			"expected error",
			"test body",
			"test body different",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, mock, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			mock.ExpectBody(test.expectedBody)

			req, err := http.NewRequest(http.MethodPost, "http://test.com", bytes.NewBufferString(test.reqBody))
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			err = mock.ExpectationsMet()
			if (err != nil) != test.wantErr {
				t.Fatalf("error check failed wantError: %t got %v", test.wantErr, err)
			}
		})
	}
}

type localhostExpectation struct {
	met bool
	msg string
}

func (l *localhostExpectation) Check(req *http.Request) {
	if strings.Contains(req.Host, "localhost") {
		l.met = true
		return
	}
}

func (l *localhostExpectation) Met() bool {
	return l.met
}

func (l *localhostExpectation) Message() string {
	return l.msg
}

func TestClientMock_Expection(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			"expected correct",
			"http://localhost",
			false,
		},
		{
			"expected error",
			"http://www.test.com",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, mock, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			mock.Expectation(&localhostExpectation{})

			res, err := client.Get(test.url)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			err = mock.ExpectationsMet()
			if (err != nil) != test.wantErr {
				t.Fatalf("error check failed wantError: %t got %v", test.wantErr, err)
			}
		})
	}
}
