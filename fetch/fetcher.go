package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.kuoruan.net/v8go-polyfills/fetch/internal"
	. "go.kuoruan.net/v8go-polyfills/internal"

	"rogchap.com/v8go"
)

const (
	UserAgentLocal  = "<local>"
	RemoteAddrLocal = "0.0.0.0"
)

var defaultLocalHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
})

type Fetcher interface {
	GetLocalHandler() http.Handler

	GetFetchFunctionCallback() v8go.FunctionCallback
}

type fetcher struct {
	// Use local handler to handle the absolute path request
	LocalHandler http.Handler
}

func NewFetcher(opt ...Option) Fetcher {
	ft := &fetcher{
		LocalHandler: defaultLocalHandler,
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
		iso, _ := ctx.Isolate()

		args := info.Args()

		resolver, _ := v8go.NewPromiseResolver(ctx)

		go func() {
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
				res, err = internal.FetchHandlerFunc(f.LocalHandler, r)
			} else {
				res, err = internal.FetchRemote(r)
			}
			if err != nil {
				resolver.Reject(wrapError(iso, err))
				return
			}

			objTmp, err := newResObjTemplate(iso, res)
			if err != nil {
				resolver.Reject(wrapError(iso, err))
				return
			}
			obj, err := objTmp.NewInstance(ctx)
			if err != nil {
				resolver.Reject(wrapError(iso, err))
				return
			}

			resolver.Resolve(obj)
		}()

		return resolver.GetPromise().Value
	}
}

func initRequest(reqUrl string, reqInit internal.RequestInit) (*internal.Request, error) {
	u, err := internal.ParseURL(reqUrl)
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

	if u.IsAbs() {
		req.Header.Set("User-Agent", UserAgentLocal)
		req.RemoteAddr = RemoteAddrLocal
	} else {
		req.Header.Set("User-Agent", UserAgent())
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

func newResObjTemplate(iso *v8go.Isolate, res *internal.Response) (*v8go.ObjectTemplate, error) {
	headers, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}

	for k, v := range res.Headers {
		if err := headers.Set(k, v); err != nil {
			return nil, err
		}
	}

	textFn, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
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

	jsonFn, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
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

	for _, e := range []struct {
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
		{Key: "text", Val: textFn},
		{Key: "json", Val: jsonFn},
	} {
		if err := resTmp.Set(e.Key, e.Val, v8go.ReadOnly); err != nil {
			return nil, err
		}
	}

	return resTmp, nil
}

func wrapError(iso *v8go.Isolate, err error) *v8go.Value {
	e, _ := v8go.NewValue(iso, fmt.Sprintf("fetch: %v", err))
	return e
}

func UserAgent() string {
	return fmt.Sprintf("v8go-polyfills/%s (v8go/%s)", Version, v8go.Version())
}
