package v8

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"go.kuoruan.net/v8go-polyfills/fetch/internal"

	"rogchap.com/v8go"
)

type fetcher struct {
	LocalHandler http.Handler
}

func (f *fetcher) goFetchSync(info *v8go.FunctionCallbackInfo) *v8go.Value {
	ctx := info.Context()
	iso, _ := ctx.Isolate()

	resolver, _ := v8go.NewPromiseResolver(ctx)

	go func() {
		args := info.Args()

		if len(args) <= 0 {
			e, _ := v8go.NewValue(iso, "1 argument required, but only 0 present.")
			resolver.Reject(e)
			return
		}

		if !args[0].IsString() {
			e, _ := v8go.NewValue(iso, "first argument should be string.")
			resolver.Reject(e)
			return
		}

		var reqInit internal.RequestInit
		if len(args) > 1 {
			if err := json.Unmarshal([]byte(args[1].String()), &reqInit); err != nil {
				resolver.Reject(wrapError(iso, err))
				return
			}
		}

		r, err := initRequest(args[0].String(), reqInit)
		if err != nil {
			resolver.Reject(wrapError(iso, err))
			return
		}

		var res *internal.Response
		if r.IsLocal {
			res, err = fetchHandlerFunc(f.LocalHandler, r)
		} else {
			res, err = fetchHttp(r)
		}
		if err != nil {
			resolver.Reject(wrapError(iso, err))
			return
		}

		resByts, err := json.Marshal(res)
		if err != nil {
			resolver.Reject(wrapError(iso, err))
			return
		}

		resVal, _ := v8go.NewValue(iso, string(resByts))
		resolver.Resolve(resVal)
	}()

	return resolver.GetPromise().Value
}

func initRequest(reqUrl string, reqInit internal.RequestInit) (*internal.Request, error) {
	u, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}

	headers := http.Header{
		"Accept":     []string{"*/*"},
		"Connection": []string{"close"},
	}

	var localURL bool
	switch u.Scheme {
	case "http", "https":
		localURL = false

		headers.Set("User-Agent", "v8go-fetch/0.0")
	case "":
		if !strings.HasPrefix(u.Path, "/") {
			return nil, fmt.Errorf("unsupported relatve path: %s", u.Path)
		}
		localURL = true

		headers.Set("User-Agent", "<local>")
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	for h, v := range reqInit.Headers {
		headerName := http.CanonicalHeaderKey(h)
		headers.Set(headerName, v)
	}

	req := &internal.Request{
		URL:     u.String(),
		Headers: headers,
		IsLocal: localURL,
	}

	if reqInit.Method != "" {
		req.Method = strings.ToUpper(reqInit.Method)
	} else {
		req.Method = "GET"
	}

	req.Body = reqInit.Body

	switch r := strings.ToLower(string(reqInit.Redirect)); r {
	case "error", "follow", "manual":
		req.Redirect = internal.RequestRedirect(r)
	case "":
		req.Redirect = internal.RequestRedirectFollow
	default:
		return nil, fmt.Errorf("unsupported redirect: %s", reqInit.Redirect)
	}

	return req, nil
}

func fetchHttp(r *internal.Request) (*internal.Response, error) {
	var body io.Reader
	if r.Method != "GET" {
		body = strings.NewReader(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}

	req.Header = r.Headers

	redirected := false
	client := http.Client{
		Transport: http.DefaultTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			switch r.Redirect {
			case internal.RequestRedirectError:
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

	return handleHttpResponse(res, r.URL, redirected)
}

func fetchHandlerFunc(handler http.Handler, r *internal.Request) (*internal.Response, error) {
	if handler == nil {
		return nil, errors.New("no server handler present")
	}

	rcd := httptest.NewRecorder()

	var body io.Reader
	if r.Method != "GET" {
		body = strings.NewReader(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Headers

	handler.ServeHTTP(rcd, req)

	result := rcd.Result()

	return handleHttpResponse(result, r.URL, false)
}

func handleHttpResponse(res *http.Response, url string, redirected bool) (*internal.Response, error) {
	resHeaders := make(map[string]string, 0)
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

	return &internal.Response{
		Headers:    resHeaders,
		Status:     res.StatusCode,
		StatusText: res.Status,
		OK:         res.StatusCode >= 200 && res.StatusCode < 300,
		Redirected: redirected,
		URL:        url,
		Body:       string(resBody),
	}, nil
}

func wrapError(iso *v8go.Isolate, err error) *v8go.Value {
	e, _ := v8go.NewValue(iso, err.Error())
	return e
}
