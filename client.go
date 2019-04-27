package clientmock

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

// Mock interface serves to create expectations
// for any kind of http.Client action in order to mock
// and test real http.Client behavior.
type Mock interface {
	// ExpectMethod sets the method expected on the incoming request
	ExpectMethod(method string) Mock
	// ExpectBody sets the body expected on the incoming request
	ExpectBody(body string) Mock
	// ExpectHeader sets the expected headers on the incoming request
	ExpectHeader(header http.Header) Mock
	// Expectation allows for any custom expectation that follows the Expectation interface
	Expectation(exp Expectation) Mock
	// ExpectationsMet will return an error if expectations weren't met
	ExpectationsMet() error
	// ReturnStatus will set the status on returned http.Response
	ReturnStatus(status int) Mock
	// ReturnBody will set the body on the returned http.Response to the given io.Reader
	ReturnBody(body io.Reader) Mock
	// ReturnHeader will set the headers on the returned http.Response to the given headers
	ReturnHeader(header http.Header) Mock
	// ReturnError will return an error instead of a response object
	ReturnError(err error) Mock
	// Add will add the current expectations, setters, and error to the stack
	Add() Mock
}

// clientMock is the default mock returned when getting a new client
type clientMock struct {
	requestResponse
	index int
	stack []requestResponse
	sync.RWMutex
}

// requestResponse contains the setters, expectations, and error for a single request
type requestResponse struct {
	expectations []Expectation
	setters      []Setter
	err          error
}

// RoundTrip implements the RoundTripper interface for use on an http.Client
func (m *clientMock) RoundTrip(req *http.Request) (*http.Response, error) {
	m.Lock()
	defer m.Unlock()

	// if stack is empty (forgot Add()) then call and add 1 request response to stack
	if len(m.expectations) != 0 || len(m.setters) != 0 || m.err != nil {
		m.Add()
	}

	rr := m.stack[m.index]
	defer func() {
		m.index++
	}()

	// if error is set return it instead of response
	if rr.err != nil {
		return nil, rr.err
	}

	// initialize response with default values
	resp := &http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(&bytes.Buffer{}),
		ContentLength: 0,
		Request:       req,
		Header:        make(http.Header, 0),
	}

	// process all setters
	for _, setter := range rr.setters {
		setter.Set(resp)
	}

	// process all expectations
	for _, exp := range rr.expectations {
		exp.Check(req)
	}

	return resp, nil
}

func (m *clientMock) ExpectMethod(method string) Mock {
	m.expectations = append(m.expectations, &ExpectedMethod{
		method: method,
	})

	return m
}

func (m *clientMock) ExpectBody(body string) Mock {
	m.expectations = append(m.expectations, &ExpectedBody{
		body: body,
	})

	return m
}

func (m *clientMock) ExpectHeader(h http.Header) Mock {
	m.expectations = append(m.expectations, &ExpectedHeader{
		h: h,
	})

	return m
}

func (m *clientMock) Expectation(exp Expectation) Mock {
	m.expectations = append(m.expectations, exp)

	return m
}

func (m *clientMock) ReturnStatus(status int) Mock {
	m.setters = append(m.setters, &SetStatusCode{
		code: status,
	})

	return m
}

func (m *clientMock) ReturnBody(body io.Reader) Mock {
	m.setters = append(m.setters, &SetBody{
		body: body,
	})

	return m
}

func (m *clientMock) ReturnHeader(header http.Header) Mock {
	m.setters = append(m.setters, &SetHeader{
		h: header,
	})

	return m
}

func (m *clientMock) ReturnError(err error) Mock {
	m.err = err
	return m
}

func (m *clientMock) ExpectationsMet() error {
	for i, rr := range m.stack {
		for _, exp := range rr.expectations {
			if !exp.Met() {
				return errors.New(fmt.Sprintf("error request [%d] message: %s", i+1, exp.Message()))
			}
		}
	}

	return nil
}

func (m *clientMock) Add() Mock {
	rr := requestResponse(m.requestResponse)

	m.stack = append(m.stack, rr)
	m.expectations = make([]Expectation, 0)
	m.setters = make([]Setter, 0)
	m.err = nil

	return m
}

// NewClient create new mockable http.Client
func NewClient() (*http.Client, Mock, error) {
	// setup mock
	m := &clientMock{
		index: 0,
	}

	// return mock and client to caller
	return &http.Client{
		Transport: m,
	}, m, nil
}
