package clientmock

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Expectation is an interface for an expectation check for the request that will be made
type Expectation interface {
	// Check will give the expectation the incoming request to be validated
	Check(req *http.Request)

	// Met will check if the check has passed
	Met() bool

	// Message is the msg passed to error.New if the expectation wasn't met
	Message() string
}

// ExpectedMethod will verify that the request method matches the expected method
type ExpectedMethod struct {
	method string
	met    bool
	msg    string
}

// Check will verify the *http.Request method matches the expected method
func (m *ExpectedMethod) Check(req *http.Request) {
	if m.method == req.Method {
		m.met = true
	}

	m.msg = fmt.Sprintf("expected method %s got %s", m.method, req.Method)
}

// Met returns if methods matched
func (m *ExpectedMethod) Met() bool {
	return m.met
}

// Message returns the error message or "" if none
func (m *ExpectedMethod) Message() string {
	return m.msg
}

// ExpectedBody verifies that the request body matches the expected body
type ExpectedBody struct {
	body string
	met  bool
	msg  string
}

// Check will verify the body in the request matches the expected body
func (e *ExpectedBody) Check(req *http.Request) {
	data, _ := ioutil.ReadAll(req.Body)

	if string(data) != e.body {
		e.msg = fmt.Sprintf("bodies don't match expected [%s] got [%s]", e.body, string(data))
		return
	}

	e.met = true
}

// Met returns whether or not the expectation has been met
func (e *ExpectedBody) Met() bool {
	return e.met
}

// Message returns the error message or "" if none
func (e *ExpectedBody) Message() string {
	return e.msg
}

// ExpectedHeader will verify the sent header matches the expected header
type ExpectedHeader struct {
	h   http.Header
	met bool
	msg string
}

// Check will verify the headers match the expected
func (e *ExpectedHeader) Check(req *http.Request) {
	e.met = true

	for key := range e.h {
		exp := e.h.Get(key)
		got := req.Header.Get(key)
		if exp != "" && exp != got {
			e.met = false
			e.msg += fmt.Sprintf("Expected %s got %s for %s\n", exp, got, key)
		}
	}
}

// Met returns whether or not the expectations were met
func (e *ExpectedHeader) Met() bool {
	return e.met
}

// Message returns the error message or "" if none
func (e *ExpectedHeader) Message() string {
	return e.msg
}
