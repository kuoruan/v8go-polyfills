/*
 * Copyright (c) 2021 Xingwang Liao
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"go.kuoruan.net/v8go-polyfills/fetch/internal"
	. "go.kuoruan.net/v8go-polyfills/internal"

	"rogchap.com/v8go"
)

const (
	UserAgentLocal = "<local>"
	AddrLocal      = "0.0.0.0:0"
)

var defaultLocalHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
})

var defaultUserAgentProvider = UserAgentProviderFunc(func(u *url.URL) string {
	if !u.IsAbs() {
		return UserAgentLocal
	}

	return UserAgent()
})

type Fetcher interface {
	GetLocalHandler() http.Handler

	GetFetchFunctionCallback() v8go.FunctionCallback
}

type fetcher struct {
	// Use local handler to handle the relative path (starts with "/") request
	LocalHandler http.Handler

	UserAgentProvider UserAgentProvider
	AddrLocal         string
}

func NewFetcher(opt ...Option) Fetcher {
	ft := &fetcher{
		LocalHandler:      defaultLocalHandler,
		UserAgentProvider: defaultUserAgentProvider,
		AddrLocal:         AddrLocal,
	}

	for _, o := range opt {
		o.apply(ft)
	}

	return ft
}

func (f *fetcher) GetLocalHandler() http.Handler {
	return f.LocalHandler
}

func (f *fetcher) GetFetchFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		args := info.Args()

		resolver, _ := v8go.NewPromiseResolver(ctx)

		go func() {
			if len(args) <= 0 {
				err := errors.New("1 argument required, but only 0 present")
				resolver.Reject(newErrorValue(ctx, err))
				return
			}

			var reqInit internal.RequestInit
			if len(args) > 1 {
				str, err := v8go.JSONStringify(ctx, args[1])
				if err != nil {
					resolver.Reject(newErrorValue(ctx, err))
					return
				}

				reader := strings.NewReader(str)
				if err := json.NewDecoder(reader).Decode(&reqInit); err != nil {
					resolver.Reject(newErrorValue(ctx, err))
					return
				}
			}

			r, err := f.initRequest(args[0].String(), reqInit)
			if err != nil {
				resolver.Reject(newErrorValue(ctx, err))
				return
			}

			var res *internal.Response

			// do local request
			if !r.URL.IsAbs() {
				res, err = f.fetchLocal(r)
			} else {
				res, err = f.fetchRemote(r)
			}
			if err != nil {
				resolver.Reject(newErrorValue(ctx, err))
				return
			}

			resObj, err := newResponseObject(ctx, res)
			if err != nil {
				resolver.Reject(newErrorValue(ctx, err))
				return
			}

			resolver.Resolve(resObj)
		}()

		return resolver.GetPromise().Value
	}
}

func (f *fetcher) initRequest(reqUrl string, reqInit internal.RequestInit) (*internal.Request, error) {
	u, err := internal.ParseRequestURL(reqUrl)
	if err != nil {
		return nil, err
	}

	req := &internal.Request{
		URL:  u,
		Body: reqInit.Body,
		Header: http.Header{
			"Accept":     []string{"*/*"},
			"Connection": []string{"close"},
		},
	}

	var ua string
	if f.UserAgentProvider != nil {
		ua = f.UserAgentProvider.GetUserAgent(u)
	} else {
		ua = defaultUserAgentProvider(u)
	}

	req.Header.Set("User-Agent", ua)

	// url has no scheme, its a local request
	if !u.IsAbs() {
		req.RemoteAddr = f.AddrLocal
	}

	for h, v := range reqInit.Headers {
		headerName := http.CanonicalHeaderKey(h)
		req.Header.Set(headerName, v)
	}

	if reqInit.Method != "" {
		req.Method = strings.ToUpper(reqInit.Method)
	} else {
		req.Method = "GET"
	}

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

func (f *fetcher) fetchLocal(r *internal.Request) (*internal.Response, error) {
	if f.LocalHandler == nil {
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

	f.LocalHandler.ServeHTTP(rcd, req)

	return internal.HandleHttpResponse(rcd.Result(), r.URL.String(), false)
}

func (f *fetcher) fetchRemote(r *internal.Request) (*internal.Response, error) {
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
	client := &http.Client{
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
		Timeout: 20 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return internal.HandleHttpResponse(res, r.URL.String(), redirected)
}

func newResponseObject(ctx *v8go.Context, res *internal.Response) (*v8go.Object, error) {
	iso, _ := ctx.Isolate()

	// create a header template,
	// if v8go supports Map, change this to a Map Object
	headersTmp, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}

	headers, err := headersTmp.NewInstance(ctx)
	if err != nil {
		return nil, err
	}

	for k, v := range res.Headers {
		if err := headers.Set(k, v); err != nil {
			return nil, err
		}
	}

	textFnTmp, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		resolver, _ := v8go.NewPromiseResolver(ctx)

		go func() {
			v, _ := v8go.NewValue(iso, res.Body)
			resolver.Resolve(v)
		}()

		return resolver.GetPromise().Value
	})
	if err != nil {
		return nil, err
	}

	jsonFnTmp, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()

		resolver, _ := v8go.NewPromiseResolver(ctx)

		go func() {
			val, err := v8go.JSONParse(ctx, res.Body)
			if err != nil {
				rejectVal, _ := v8go.NewValue(iso, err.Error())
				resolver.Reject(rejectVal)
				return
			}

			resolver.Resolve(val)
		}()

		return resolver.GetPromise().Value
	})
	if err != nil {
		return nil, err
	}

	resTmp, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}

	for _, f := range []struct {
		Name string
		Tmp  interface{}
	}{
		{Name: "text", Tmp: textFnTmp},
		{Name: "json", Tmp: jsonFnTmp},
	} {
		if err := resTmp.Set(f.Name, f.Tmp, v8go.ReadOnly); err != nil {
			return nil, err
		}
	}

	resObj, err := resTmp.NewInstance(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range []struct {
		Key string
		Val interface{}
	}{
		{Key: "headers", Val: headers},
		{Key: "ok", Val: res.OK},
		{Key: "redirected", Val: res.Redirected},
		{Key: "status", Val: res.Status},
		{Key: "statusText", Val: res.StatusText},
		{Key: "url", Val: res.URL},
		{Key: "body", Val: res.Body},
	} {
		if err := resObj.Set(v.Key, v.Val); err != nil {
			return nil, err
		}
	}

	return resObj, nil
}

// v8go currently not support reject a *v8go.Object,
// so we should new *v8go.Value here
func newErrorValue(ctx *v8go.Context, err error) *v8go.Value {
	iso, _ := ctx.Isolate()
	e, _ := v8go.NewValue(iso, fmt.Sprintf("fetch: %v", err))
	return e
}

func UserAgent() string {
	return fmt.Sprintf("v8go-polyfills/%s (v8go/%s)", Version, v8go.Version())
}
