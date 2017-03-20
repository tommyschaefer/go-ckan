package ckan

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "1"
	userAgent      = "go-ckan/" + libraryVersion
	mediaType      = "application/json"
)

// Client manages communication with a CKAN V3 API.
type Client struct {
	// HTTP client used to communicate with the CKAN API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string
}

// // An ErrorResponse reports the error caused by an API request
// type ErrorResponse struct {
// 	// HTTP response that caused this error
// 	Response *http.Response

// 	// Error message
// 	Message string `jsonon:"message"`

// 	// RequestID returned from the API, useful to contact stringupport.
// 	RequestID string `json:"request_id"`
// }

// NewClient returns a new CKAN API client.
func NewClient(baseURL *url.URL, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the containsext's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	// If the status code is non-200, a formatting error has occurred, and
	// JSON decoding should not occur.
	if c := resp.StatusCode; c < 200 || c > 299 {
		return resp, &FormattingError{}
	}

	r := new(response)
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return resp, err
	}

	if !r.Success {
		return resp, r.Error
	}

	if v != nil {
		err = json.Unmarshal(r.Result, v)
	}

	return resp, err
}

// response is used to parse the returned JSON from the CKAN API. The result
// is saved in a RawMessage to allow marshalling to a provided type in Do.
type response struct {
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result"`
	Error   *Error          `json:"error"`
	Help    string          `json:"help"`
}

// FormattingError occurs when a CKAN API returns a non-200 status code. This
// indicates that there was a formatting error with the request, and the
// request could not be processed by the CKAN server.
//
//
// From the docs:
//
//   If there are major formatting problems with a request to the API, CKAN
//   may still return an HTTP response with a 409, 400 or 500 status code (in
//   increasing order of severity). In future CKAN versions we intend to
//   remove these responses, and instead send a 200 OK response and use the
//   "success" and "error" items.
type FormattingError struct{}

func (*FormattingError) Error() string {
	return "request malformed and no body was returned"
}

// An Error reports the error returned by an API request
type Error struct {
	Message string `json:"message"`
	Type    string `json:"__type"`
}

func (e *Error) Error() string {
	return e.Message
}