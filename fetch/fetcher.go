package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"go.kuoruan.net/v8go-polyfills/fetch/internal"

	"rogchap.com/v8go"
)

type fetcher struct {
	// Use local handler to handle the absolute path request
	LocalHandler http.Handler
}

func (f *fetcher) goFetchSync(info *v8go.FunctionCallbackInfo) *v8go.Value {
	ctx := info.Context()
	iso, _ := ctx.Isolate()

	resolver, _ := v8go.NewPromiseResolver(ctx)

	go func() {
		args := info.Args()

		if len(args) <= 0 {
			err := errors.New("1 argument required, but only 0 present")
			resolver.Reject(wrapError(iso, err))
			return
		}

		if !args[0].IsString() {
			err := errors.New("first argument should be string")
			resolver.Reject(wrapError(iso, err))
			return
		}

		var reqInit internal.RequestInit
		if len(args) > 1 {
			str, err := v8go.JSONStringify(ctx, args[1])
			if err != nil {
				resolver.Reject(wrapError(iso, err))
				return
			}

			reader := strings.NewReader(str)
			if err := json.NewDecoder(reader).Decode(&reqInit); err != nil {
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
		if r.URL.IsAbs() {
			res, err = fetchHandlerFunc(f.LocalHandler, r)
		} else {
			res, err = fetchHttp(r)
		}
		if err != nil {
			resolver.Reject(wrapError(iso, err))
			return
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(res); err != nil {
			resolver.Reject(wrapError(iso, err))
			return
		}

		val, err := v8go.JSONParse(ctx, buf.String())
		if err != nil {
			resolver.Reject(wrapError(iso, err))
			return
		}
		resolver.Resolve(val)
	}()

	return resolver.GetPromise().Value
}

func initRequest(reqUrl string, reqInit internal.RequestInit) (*internal.Request, error) {
	u, err := internal.ParseURL(reqUrl)
	if err != nil {
		return nil, err
	}

	headers := http.Header{
		"Accept":     []string{"*/*"},
		"Connection": []string{"close"},
	}

	if u.IsAbs() {
		headers.Set("User-Agent", UserAgentLocal)
	} else {
		headers.Set("User-Agent", UserAgent())
	}

	for h, v := range reqInit.Headers {
		headerName := http.CanonicalHeaderKey(h)
		headers.Set(headerName, v)
	}

	req := &internal.Request{
		URL:     u,
		Headers: headers,
	}

	if reqInit.Method != "" {
		req.Method = strings.ToUpper(reqInit.Method)
	} else {
		req.Method = "GET"
	}

	req.Body = reqInit.Body

	switch r := strings.ToLower(reqInit.Redirect); r {
	case "error", "follow", "manual":
		req.Redirect = r
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

	req, err := http.NewRequest(r.Method, r.URL.String(), body)
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

	return handleHttpResponse(res, r.URL.String(), redirected)
}

func fetchHandlerFunc(handler http.Handler, r *internal.Request) (*internal.Response, error) {
	if handler == nil {
		return nil, errors.New("fetch: no server handler present")
	}

	var body io.Reader
	if r.Method != "GET" {
		body = strings.NewReader(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.URL.String(), body)
	if err != nil {
		return nil, err
	}
	req.RemoteAddr = RemoteAddrLocal
	req.Header = r.Headers

	rcd := httptest.NewRecorder()

	handler.ServeHTTP(rcd, req)

	return handleHttpResponse(rcd.Result(), r.URL.String(), false)
}

func handleHttpResponse(res *http.Response, url string, redirected bool) (*internal.Response, error) {
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
	e, _ := v8go.NewValue(iso, fmt.Sprintf("fetch: %v", err))
	return e
}
