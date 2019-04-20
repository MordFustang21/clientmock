package clientmock

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Setter is an interface that sets a value on an *http.Response
type Setter interface {
	Set(response *http.Response)
}

// SetStatusCode implements the Setter interface and sets the status code
type SetStatusCode struct {
	code int
}

// Set will set the response status code with the returning status
func (s *SetStatusCode) Set(resp *http.Response) {
	resp.StatusCode = s.code
	resp.Status = http.StatusText(s.code)
}

// SetBody implements the Setter interface and sets the status code
type SetBody struct {
	body io.Reader
}

// Set will set the response body with the returning body
func (s *SetBody) Set(response *http.Response) {
	response.Body = ioutil.NopCloser(s.body)
}

// SetHeader implements the Setter interface and sets the response headers
type SetHeader struct {
	h http.Header
}

// Set sets the response headers on the request
func (s *SetHeader) Set(response *http.Response) {
	response.Header = s.h
}