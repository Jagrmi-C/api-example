go-http-client
==============

The `gohttp` package provides a familiar HTTP client interface with
automatic retries and exponential backoff. It is a thin wrapper over the
standard `net/http` client library and exposes nearly the same public API.

`gohttp` performs automatic retries under certain conditions. Mainly, if
an error is returned by the client (connection errors, etc.), or if a 500-range
response code is received (except 501), then a retry is invoked after a wait
period.  Otherwise, the response is returned and left to the caller to
interpret.

===========
Example Use
===========

This package provides a functionality that allows you to add custom retries to the client, headers, timeouts.
The simple example of a request is shown below:

```go
...
preparedJSONHeaders := make(http.Header, 1)
preparedJSONHeaders.Add(gohttp.HeaderContentType, gohttp, ContentTypeJSON)

defaultJsonClient := gohttp.NewClientBuilder().
    SetHeaders(preparedJSONHeaders).
    SetRetryMax(4).
    Build()
resp, err := defaultJsonClient.DoReq(ctx, attrs)
if err != nil {
    return err
}
...
```

The returned response object is an `*gohttp.Response`, the same thing you would
usually get from `net/http`, but this structure contains as field a response body. Had the request failed one or more times, the above
call would block and retry with exponential backoff.

## Getting a stdlib `*http.Client` with retries

It's possible to convert a `*gohttp.Client` directly to a `*http.Client`.
Simply configure a `*gohttp.Client` as you wish, and then call `StandardClient()`:

```go
defaultJsonClient := gohttp.NewClientBuilder().
 SetHeaders(preparedJSONHeaders).
 Build()

standardClient := defaultJsonClient.Client() // *http.Client
```
