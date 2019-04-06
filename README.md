# clientmock
[![GoDoc](https://godoc.org/github.com/MordFustang21/clientmock?status.svg)](https://godoc.org/github.com/MordFustang21/supernova)
[![Go Report Card](https://goreportcard.com/badge/github.com/MordFustang21/clientmock)](https://goreportcard.com/report/github.com/mordfustang21/supernova)

A mockable http.Client used to validate request data and return fake response data

Example returning status code and body
```go
package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	
	"github.com/MordFustang21/clientmock"
)

func TestGet(t *testing.T) {
	client, mock, err := clientmock.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	
	expectedStatus := http.StatusBadRequest
	expectedBody := "empty request body"
	
	mock.ReturnStatus(expectedStatus)
	mock.ReturnBody(bytes.NewBufferString(expectedBody))
	
	req, err := http.NewRequest(http.MethodPost, "http://www.github.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	
	// Status code will be 400 as defined in mock
	if res.StatusCode != expectedStatus {
		t.Fatal("got wrong status")
	}
	
	// body will contain "empty request body" as defined in mock
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	
	if string(data) != expectedBody {
		t.Fatal("returned body doesn't match expected")
	}
}

```