package clientmock

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
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
}

// clientMock is the default mock returned when getting a new client
type clientMock struct {
	expectations []Expectation
	setters      []Setter
}

// RoundTrip implements the RoundTripper interface for use on an http.Client
func (m *clientMock) RoundTrip(req *http.Request) (*http.Response, error) {

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
	for _, setter := range m.setters {
		setter.Set(resp)
	}

	// process all expectations
	for _, exp := range m.expectations {
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

func (m *clientMock) ExpectationsMet() error {
	for _, exp := range m.expectations {
		if !exp.Met() {
			return errors.New(exp.Message())
		}
	}

	return nil
}

// NewClient create new mockable http.Client
func NewClient() (*http.Client, Mock, error) {
	// setup mock
	m := &clientMock{}

	// return mock and client to caller
	return &http.Client{
		Transport: m,
	}, m, nil
}
