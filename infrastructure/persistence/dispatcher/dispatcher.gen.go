// Package internal provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package dispatcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// ApiInternalRunsCreate request  with any body
	ApiInternalRunsCreateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ApiInternalRunsCreate(ctx context.Context, body ApiInternalRunsCreateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ApiInternalRunsCreateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewApiInternalRunsCreateRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ApiInternalRunsCreate(ctx context.Context, body ApiInternalRunsCreateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewApiInternalRunsCreateRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewApiInternalRunsCreateRequest calls the generic ApiInternalRunsCreate builder with application/json body
func NewApiInternalRunsCreateRequest(server string, body ApiInternalRunsCreateJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewApiInternalRunsCreateRequestWithBody(server, "application/json", bodyReader)
}

// NewApiInternalRunsCreateRequestWithBody generates requests for ApiInternalRunsCreate with any type of body
func NewApiInternalRunsCreateRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/internal/dispatch")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	req = req.WithContext(ctx)
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// ApiInternalRunsCreate request  with any body
	ApiInternalRunsCreateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*ApiInternalRunsCreateResponse, error)

	ApiInternalRunsCreateWithResponse(ctx context.Context, body ApiInternalRunsCreateJSONRequestBody) (*ApiInternalRunsCreateResponse, error)
}

type ApiInternalRunsCreateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON207      *RunsCreated
}

// Status returns HTTPResponse.Status
func (r ApiInternalRunsCreateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ApiInternalRunsCreateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ApiInternalRunsCreateWithBodyWithResponse request with arbitrary body returning *ApiInternalRunsCreateResponse
func (c *ClientWithResponses) ApiInternalRunsCreateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*ApiInternalRunsCreateResponse, error) {
	rsp, err := c.ApiInternalRunsCreateWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseApiInternalRunsCreateResponse(rsp)
}

func (c *ClientWithResponses) ApiInternalRunsCreateWithResponse(ctx context.Context, body ApiInternalRunsCreateJSONRequestBody) (*ApiInternalRunsCreateResponse, error) {
	rsp, err := c.ApiInternalRunsCreate(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseApiInternalRunsCreateResponse(rsp)
}

// ParseApiInternalRunsCreateResponse parses an HTTP response from a ApiInternalRunsCreateWithResponse call
func ParseApiInternalRunsCreateResponse(rsp *http.Response) (*ApiInternalRunsCreateResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ApiInternalRunsCreateResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 207:
		var dest RunsCreated
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON207 = &dest

	}

	return response, nil
}
