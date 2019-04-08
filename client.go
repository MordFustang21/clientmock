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
	ExpectMethod(method string)
	ExpectBody(body string)
	ExpectPath(path string)
	ExpectationsMet() error
	Expectation(exp Expectation)
	ReturnStatus(status int)
	ReturnBody(body io.Reader)
}

type clientMock struct {
	expectations []Expectation
	setters      []Setter
}

// RoundTrip implements the RoundTripper interface for use on an http.Client
func (m *clientMock) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, exp := range m.expectations {
		exp.Check(req)
	}

	resp := &http.Response{}

	for _, setter := range m.setters {
		setter.Set(resp)
	}

	if resp.Body == nil {
		resp.Body = ioutil.NopCloser(&bytes.Buffer{})
	}

	return resp, nil
}

func (m *clientMock) ExpectMethod(method string) {
	m.expectations = append(m.expectations, &ExpectedMethod{
		method: method,
	})
}

func (m *clientMock) ExpectBody(body string) {
	m.expectations = append(m.expectations, &ExpectedBody{
		body: body,
	})
}

func (m *clientMock) ExpectPath(path string) {
	m.expectations = append(m.expectations, &ExpectedPath{
		path: path,
	})
}

func (m *clientMock) Expectation(exp Expectation) {
	m.expectations = append(m.expectations, exp)
}

func (m *clientMock) ReturnStatus(status int) {
	m.setters = append(m.setters, &SetStatusCode{
		code: status,
	})
}

func (m *clientMock) ReturnBody(body io.Reader) {
	m.setters = append(m.setters, &SetBody{
		body: body,
	})
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
