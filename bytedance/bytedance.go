package bytedance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://open.microapp.bytedance.com/openapi/"
	userAgent      = "go-bytedance"
)

// Client manages communication with the bytedance API.
type Client struct {
	client *http.Client

	// Base URL for API.requests.
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communication with bytedance API.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of bytedance API.
	ThirdParty *ThirdPartyService
	MicroApp   *MicroAppService
}

type service struct {
	client *Client
}

// NewClient returns a new bytedance API client.
func NewClient() *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := Client{BaseURL: baseURL, UserAgent: userAgent, client: &http.Client{}}
	c.common.client = &c
	c.ThirdParty = (*ThirdPartyService)(&c.common)
	c.MicroApp = (*MicroAppService)(&c.common)
	return &c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var contentType string
	var buf bytes.Buffer
	if body != nil {

		if render, ok := body.(FormRender); ok {
			w := multipart.NewWriter(&buf)
			for key, r := range render.MultipartParams() {
				var fw io.Writer
				if x, ok := r.(*File); ok {
					if fw, err = w.CreateFormFile(key, x.Name); err != nil {
						return nil, err
					}
				}
				if _, err = io.Copy(fw, r); err != nil {
					return nil, err
				}
			}
			for k, v := range render.Params() {
				_ = w.WriteField(k, v)
			}
			if err = w.Close(); err != nil {
				return nil, err
			}
			contentType = w.FormDataContentType()
		} else {
			enc := json.NewEncoder(&buf)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(body); err != nil {
				return nil, err
			}
			contentType = "application/json"

		}

	}

	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// FormRender for form with file request.
type FormRender interface {
	Params() map[string]string
	MultipartParams() map[string]io.Reader
}

// Response represents bytedance Response.
type Response struct {
	*http.Response

	ErrNo   int             `json:"errno"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// ErrorResponse represents bytedance Error Response.
type ErrorResponse struct {
	*http.Response
	requestBody []byte

	ErrNo   int    `json:"errno"`
	Message string `json:"message"`
}

// Error implements builtin.error interface.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v body %s : %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.requestBody,
		r.ErrNo, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
// API error responses are expected to have response body.
// If none error find, return bytedance.Response.
func CheckResponse(r *http.Response, reqBody []byte) error {
	errResponse := &ErrorResponse{
		Response:    r,
		requestBody: reqBody,
	}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		if err = json.Unmarshal(data, &errResponse); err != nil {
			return err
		}
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	switch errResponse.ErrNo {
	case 0:
		return nil
	default:
		return errResponse
	}
}

// Do Sends an API request and returns the API response. The API reponse is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occured. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or timeout, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}

	req = req.WithContext(ctx)

	// save body for display request error info.
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.New("read request body error")
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error id probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err

	}

	if err := CheckResponse(resp, body); err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore io.EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}
	return resp, err

}
