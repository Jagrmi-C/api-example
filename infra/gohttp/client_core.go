package gohttp

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/idtoken"
)

// IAMClientError is an error that occur when interacting with an IAM (Identity and Access Management) service client.
// It is a wrapper around the standard Go `error` interface,
// with additional information or context related to the specific error.
type IAMClientError struct {
	error
}

func (e IAMClientError) Error() string {
	return fmt.Sprintf("creation an iam client: %v", e.error)
}

var (
	// A regular expression to match the error returned by net/http when the
	// configured number of redirects is exhausted. This error isn't typed
	// specifically so we resort to matching on the error string.
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

	// A regular expression to match the error returned by net/http when the
	// scheme specified in the URL is invalid. This error isn't typed
	// specifically so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

	// A regular expression to match the error returned by net/http when the
	// TLS certificate is not trusted. This error isn't typed
	// specifically so we resort to matching on the error string.
	notTrustedErrorRe = regexp.MustCompile(`certificate is not trusted`)
)

// headers.
const (
	HeaderContentType = "Content-Type"
	HeaderUserAgent   = "User-Agent"
	HeaderAuth        = "Authorization"

	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
)

// default timeouts.
const (
	// DefaultTLSHandshakeTimeout for TLS handshake.
	DefaultTLSHandshakeTimeout = 10 * time.Second

	// DefaultResponseHeaderTimeout  for waiting to read a response header.
	DefaultResponseHeaderTimeout = 5 * time.Second

	// DefaultConnectionTimeout ...
	DefaultConnectionTimeout = 2 * time.Second

	// DefaultMaxIdleConnections to keep in pool.
	DefaultMaxIdleConnections = 15
)

// Logger interface allows to use other loggers than standard log.Logger.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Client is a getter that returns a http client,
// that contains under the hood logic to create a client only once.
func (c *httpClient) Client() *http.Client {
	c.clientOnce.Do(func() {
		if c.b.client != nil {
			c.client = c.b.client
			return
		}

		c.client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig(c.b),
				DialContext: (&net.Dialer{
					Timeout: c.ConnectionTimeout(),
				}).DialContext,
				ResponseHeaderTimeout: c.ResponseTimeout(),
			},
			Timeout: c.ConnectionTimeout() + c.ResponseTimeout(),
		}
	})

	return c.client
}

// ResponseTimeout is a getter for the response timeout
// if exists in builder config retrieve this one, if not
// set it defaults to.
func (c *httpClient) ResponseTimeout() time.Duration {
	if c.b.responseHeaderTimeOut > 0 {
		return c.b.responseHeaderTimeOut
	}

	if c.b.disableTimeouts {
		return 0
	}

	return DefaultResponseHeaderTimeout
}

// ConnectionTimeout is a getter for the connection timeout
// if exists in builder config retrieve this one, if not
// set it defaults to.
func (c *httpClient) ConnectionTimeout() time.Duration {
	if c.b.connectionTimeout > 0 {
		return c.b.connectionTimeout
	}

	if c.b.disableTimeouts {
		return 0
	}

	return DefaultConnectionTimeout
}

// RetryMax is a getter for max retries
// if exists in builder config retrieve this one, if not
// set it defaults to.
func (c *httpClient) RetryMax() int {
	if c.b.RetryMax != 0 {
		return c.b.RetryMax
	}

	return defaultRetryMax
}

// RequestPayload is a getter that prepare payload body from the request structure
// by the content type.
func (c *httpClient) RequestPayload(contentType string, reqBody interface{}) ([]byte, error) {
	if reqBody == nil {
		return nil, nil
	}

	switch {
	case contentType == ContentTypeJSON:
		return json.Marshal(reqBody)
	case contentType == ContentTypeXML:
		return xml.Marshal(reqBody)
	case strings.Contains(contentType, "multipart/form-data"):
		buf, ok := reqBody.(*bytes.Buffer)
		if !ok {
			return nil, fmt.Errorf("some issue")
		}

		return buf.Bytes(), nil
	default:
		return json.Marshal(reqBody)
	}
}

// RequestHeaders is a getter that returns all headers
// [client and request].
func (c *httpClient) RequestHeaders(requestHeaders http.Header) http.Header {
	result := make(http.Header)

	// add common headers
	for k, v := range c.b.headers {
		if len(v) > 0 {
			result.Set(k, v[0])
		}
	}

	// add custom headers
	for k, v := range requestHeaders {
		if len(v) > 0 {
			result.Set(k, v[0])
		}
	}

	return result
}

// StandardClient returns a stdlib *http.Client.
func (c *httpClient) StandardClient() *http.Client {
	return c.client
}

// Logger provides an access to the logger it it exists for the client.
func (c *httpClient) Logger() interface{} {
	return c.b.Logger
}

// AddTraceIDToLogger sets  the trace identifier.
func (c *httpClient) AddTraceIDToLogger(traceID string) {
	if c.b.Logger == nil {
		return
	}
}

// NewClientWhitIAMAuth returns a new HTTP client that is authorized to make requests to the specified audience using an ID token.
func NewClientWhitIAMAuth(ctx context.Context, audience string) (*http.Client, error) {
	initializedClient, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return nil, &IAMClientError{err}
	}

	initializedClient.Timeout = DefaultResponseHeaderTimeout + DefaultConnectionTimeout

	return initializedClient, nil
}

// DefaultBackoff provides a default callback for Client.Backoff which
// will perform exponential backoff based on the attempt number and limited
// by the provided minimum and maximum durations.
//
// It also tries to parse Retry-After response header when a http.StatusTooManyRequests
// (HTTP Code 429) is found in the resp parameter. Hence it will return the number of
// seconds the server states it may be ready to process more requests from this client.
func DefaultBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
					return time.Second * time.Duration(sleep)
				}
			}
		}
	}

	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}

// LinearJitterBackoff provides a callback for Client.Backoff which will
// perform linear backoff based on the attempt number and with jitter to
// prevent a thundering herd.
//
// min and max here are *not* absolute values. The number to be multiplied by
// the attempt number will be chosen at random from between them, thus they are
// bounding the jitter.
//
// For instance:
// * To get strictly linear backoff of one second increasing each retry, set
// both to one second (1s, 2s, 3s, 4s, ...)
// * To get a small amount of jitter centered around one second increasing each
// retry, set to around one second, such as a min of 800ms and max of 1200ms
// (892ms, 2102ms, 2945ms, 4312ms, ...)
// * To get extreme jitter, set to a very wide spread, such as a min of 100ms
// and a max of 20s (15382ms, 292ms, 51321ms, 35234ms, ...)
func LinearJitterBackoff(min, max time.Duration, attemptNum int, _ *http.Response) time.Duration {
	// attemptNum always starts at zero but we want to start at 1 for multiplication
	attemptNum++

	if max <= min {
		// Unclear what to do here, or they are the same, so return min *
		// attemptNum
		return min * time.Duration(attemptNum)
	}

	// Seed randVal; doing this every time is fine
	randVal := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	// Pick a random number that lies somewhere between the min and max and
	// multiply by the attemptNum. attemptNum starts at zero so we always
	// increment here. We first get a random percentage, then apply that to the
	// difference between min and max, and add to min.
	jitter := randVal.Float64() * float64(max-min)
	jitterMin := int64(jitter) + int64(min)
	return time.Duration(jitterMin * int64(attemptNum))
}

// DefaultRetryPolicy provides a default callback for Client.CheckRetry, which
// will retry on connection errors and server errors.
func DefaultRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// don't propagate other errors
	shouldRetry, _ := baseRetryPolicy(resp, err)
	return shouldRetry, nil
}

func baseRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		if val, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(val.Error()) {
				return false, val
			}

			// Don't retry if the error was due to an invalid protocol scheme.
			if schemeErrorRe.MatchString(val.Error()) {
				return false, val
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if notTrustedErrorRe.MatchString(val.Error()) {
				return false, val
			}
			if _, ok := val.Err.(x509.UnknownAuthorityError); ok {
				return false, val
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}

	// 429 Too Many Requests is recoverable.
	if resp.StatusCode == http.StatusTooManyRequests {
		return true, nil
	}

	// Check the response code. We retry on 500-range responses to allow
	// the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side. This will catch
	// invalid response codes as well, like 0 and 999.
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented) {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	return false, nil
}
