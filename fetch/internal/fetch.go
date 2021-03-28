package internal

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func FetchHandlerFunc(handler http.Handler, r *Request) (*Response, error) {
	if handler == nil {
		return nil, errors.New("no local handler present")
	}

	var body io.Reader
	if r.Method != "GET" {
		body = strings.NewReader(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.URL.String(), body)
	if err != nil {
		return nil, err
	}
	req.RemoteAddr = r.RemoteAddr
	req.Header = r.Header

	rcd := httptest.NewRecorder()

	handler.ServeHTTP(rcd, req)

	return handleHttpResponse(rcd.Result(), r.URL.String(), false)
}

func FetchRemote(r *Request) (*Response, error) {
	var body io.Reader
	if r.Method != "GET" {
		body = strings.NewReader(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.URL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header

	redirected := false
	client := http.Client{
		Transport: http.DefaultTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			switch r.Redirect {
			case RequestRedirectError:
				return errors.New("redirects are not allowed")
			default:
				if len(via) >= 10 {
					return errors.New("stopped after 10 redirects")
				}
			}

			redirected = true
			return nil
		},
		Timeout: 20,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return handleHttpResponse(res, r.URL.String(), redirected)
}

func handleHttpResponse(res *http.Response, url string, redirected bool) (*Response, error) {
	resHeaders := make(map[string]string)
	for k, v := range res.Header {
		for _, vv := range v {
			resHeaders[k] = vv
			break
		}
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Headers:    resHeaders,
		Status:     int32(res.StatusCode), // int type is not support by v8go
		StatusText: res.Status,
		OK:         res.StatusCode >= 200 && res.StatusCode < 300,
		Redirected: redirected,
		URL:        url,
		Body:       string(resBody),
	}, nil
}
