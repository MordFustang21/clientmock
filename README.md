# clientmock
[![GoDoc](https://godoc.org/github.com/MordFustang21/clientmock?status.svg)](https://godoc.org/github.com/MordFustang21/clientmock)
[![Go Report Card](https://goreportcard.com/badge/github.com/MordFustang21/clientmock)](https://goreportcard.com/report/github.com/MordFustang21/clientmock)
[![codecov](https://codecov.io/gh/MordFustang21/clientmock/branch/master/graph/badge.svg)](https://codecov.io/gh/MordFustang21/clientmock)

A mockable http.Client used to validate request data and return fake response data

## Install

    go get github.com/DATA-DOG/go-sqlmock

## Documentation and Examples

### Example returning status code and body
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

## Contributions

Feel free to open a pull request. Note, if you wish to contribute an extension to public (exported methods or types) -
please open an issue before, to discuss whether these changes can be accepted. All backward incompatible changes are
and will be treated cautiously

## License
The [GNU General Public License V3.0](https://www.gnu.org/licenses/gpl-3.0.en.html)