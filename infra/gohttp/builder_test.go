package gohttp_test

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/infra/gohttp"
)

const (
	DefaultCustomTimeout = 10 * time.Second
)

type FakeRequestDTO struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func Test_clientBuilder_Build(t *testing.T) {
	t.Run(`The builder was builded without options`,
		func(t *testing.T) {
			myClient := gohttp.NewClientBuilder().Build()

			httpClient := myClient.Client()

			assert.NotNil(t, httpClient, "clientBuilder.Build() have to return not nil")

			assert.Equal(
				t,
				gohttp.DefaultResponseHeaderTimeout,
				myClient.ResponseTimeout(),
			)

			assert.Equal(
				t,
				gohttp.DefaultConnectionTimeout,
				myClient.ConnectionTimeout(),
			)

			h := myClient.RequestHeaders(nil)
			assert.Empty(t, h)

			ct := h.Get(gohttp.HeaderContentType)
			assert.Empty(t, ct, "headerContentType should be empty")

			payload, err := myClient.RequestPayload(ct, FakeRequestDTO{
				Field1: "1",
				Field2: "2",
			})

			assert.NoError(t, err)
			assert.IsType(t, []byte{}, payload)
		},
	)

	t.Run(`The builder was builded as application/json client
	with custom timeouts and headers`,
		func(t *testing.T) {
			customHeaders := http.Header{}
			customHeaders.Add(
				gohttp.HeaderContentType, gohttp.ContentTypeJSON,
			)

			customTimeOut := 5 * time.Second

			myClient := gohttp.NewClientBuilder().
				SetConnectionTimeOut(customTimeOut).
				SetResponseHeaderTimeout(customTimeOut).
				SetHeaders(customHeaders).
				Build()

			httpClient := myClient.Client()

			assert.NotNil(t, httpClient, "clientBuilder.Build() have to return not nil")

			assert.Equal(
				t,
				customTimeOut,
				myClient.ResponseTimeout(),
			)

			assert.Equal(
				t,
				customTimeOut,
				myClient.ConnectionTimeout(),
			)

			requestHeaders := http.Header{}
			requestHeaders.Add(
				gohttp.HeaderAuth, "Bearer ...",
			)

			h := myClient.RequestHeaders(requestHeaders)
			assert.Len(t, h, 2)

			ct := h.Get(gohttp.HeaderContentType)
			assert.Equal(
				t,
				gohttp.ContentTypeJSON,
				ct,
				"headerContentType should return application/json",
			)

			payload, err := myClient.RequestPayload(ct, FakeRequestDTO{
				Field1: "1",
				Field2: "2",
			})

			assert.NoError(t, err)
			assert.IsType(t, []byte{}, payload)
		},
	)

	t.Run(`The builder was builded  with custom IAM http client as application/xml client
	with custom timeouts and headers`,
		func(t *testing.T) {
			myHTTPClient := &http.Client{
				Timeout: DefaultCustomTimeout,
				Transport: &http.Transport{
					DialContext: (&net.Dialer{
						Timeout: DefaultCustomTimeout,
					}).DialContext,
					ResponseHeaderTimeout: DefaultCustomTimeout,
				},
			}

			customHeaders := http.Header{}
			customHeaders.Add(
				gohttp.HeaderContentType, gohttp.ContentTypeXML,
			)
			customHeaders.Add(
				"X-Amz-Target", "test",
			)

			customTimeOut := 12 * time.Second

			myClient := gohttp.NewClientBuilder().
				SetHTTPClient(myHTTPClient).
				SetConnectionTimeOut(customTimeOut).
				SetResponseHeaderTimeout(customTimeOut).
				SetHeaders(customHeaders).
				Build()

			httpClient := myClient.Client()

			q3 := httpClient.Timeout.Seconds()
			_ = q3
			assert.Equal(t, DefaultCustomTimeout, httpClient.Timeout)
			assert.NotNil(t, httpClient, "clientBuilder.Build() have to return not nil")

			transport, ok := httpClient.Transport.(*http.Transport)
			assert.True(t, ok, "Failed to get DialContext from Transport")
			assert.Equal(t, DefaultCustomTimeout, transport.ResponseHeaderTimeout)

			requestHeaders := http.Header{}
			requestHeaders.Add(
				gohttp.HeaderAuth, "Bearer ...",
			)

			h := myClient.RequestHeaders(requestHeaders)
			assert.Len(t, h, 3)

			ct := h.Get(gohttp.HeaderContentType)
			assert.Equal(
				t,
				gohttp.ContentTypeXML,
				ct,
				"headerContentType should return application/xml",
			)

			payload, err := myClient.RequestPayload(ct, FakeRequestDTO{
				Field1: "1",
				Field2: "2",
			})

			assert.NoError(t, err)
			assert.IsType(t, []byte{}, payload)
		},
	)
}
