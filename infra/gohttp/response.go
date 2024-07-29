package gohttp

import (
	"encoding/json"
	"net/http"
)

// ResponseAttributes represents the attributes of an HTTP response that a client can use to process the response.
type ResponseAttributes struct {
	status     string
	statusCode int
	headers    http.Header
	body       []byte
}

// NewResponseAttributes is a factory function that returns a pointer to a new ResponseAttributes struct
// with the given status code, response body, and headers.
func NewResponseAttributes(statusCode int, body []byte, headers http.Header) *ResponseAttributes {
	return &ResponseAttributes{
		statusCode: statusCode,
		body:       body,
		headers:    headers,
	}
}

// StatusCode is a getter for the status code.
func (r *ResponseAttributes) StatusCode() int {
	return r.statusCode
}

// Headers is a getter for the headers.
func (r *ResponseAttributes) Headers() http.Header {
	return r.headers
}

// Bytes is a getter for the bytes of the response body.
func (r *ResponseAttributes) Bytes() []byte {
	return r.body
}

// String is a getter for the string representation of the response body.
func (r *ResponseAttributes) String() string {
	return string(r.body)
}

// UnMarshalJson is a method that unmarshals the response body into a target struct.
// The target argument should be a pointer to the struct that the response body
// will be unmarshaled into.
func (r *ResponseAttributes) UnMarshalJson(target interface{}) error {
	return json.Unmarshal(r.body, target)
}
