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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"rogchap.com/v8go"
)

func TestNewFetcher(t *testing.T) {
	t.Parallel()

	f1 := NewFetcher()
	if f1 == nil {
		t.Error("create fetcher failed")
		return
	}

	if h := f1.GetLocalHandler(); h == nil {
		t.Error("local handler is <nil>")
	}

	f2 := NewFetcher(WithLocalHandler(nil))
	if f2 == nil {
		t.Error("create fetcher with local handler failed")
		return
	}

	if h := f2.GetLocalHandler(); h != nil {
		t.Error("set fetcher local handler to <nil> failed")
		return
	}
}

func TestFetchJSON(t *testing.T) {
	t.Parallel()

	ctx, err := newV8ContextWithFetch()
	if err != nil {
		t.Errorf("create v8: %s", err)
		return
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; utf-8")
		_, _ = w.Write([]byte(`{"status": true}`))
	}))

	val, err := ctx.RunScript(fmt.Sprintf("fetch('%s').then(res => res.json())", srv.URL), "fetch_json.js")
	if err != nil {
		t.Error(err)
		return
	}

	proms, err := val.AsPromise()
	if err != nil {
		t.Error(err)
		return
	}

	for proms.State() == v8go.Pending {
		continue
	}

	res, err := proms.Result().AsObject()
	if err != nil {
		t.Error(err)
		return
	}

	status, err := res.Get("status")
	if err != nil {
		t.Error(err)
		return
	}

	if !status.Boolean() {
		t.Error("status should be true")
	}
}

func TestHeaders(t *testing.T) {
	t.Parallel()

	iso := v8go.NewIsolate()

	ctx := v8go.NewContext(iso)

	obj, err := newHeadersObject(ctx, http.Header{
		"AA": []string{"aa"},
		"BB": []string{"bb"},
	})
	if err != nil {
		t.Error(err)
		return
	}

	aa, err := obj.Get("AA")
	if err != nil {
		t.Error(err)
		return
	}

	if aa.String() != "aa" {
		t.Errorf("should be 'aa' but is '%s'", aa.String())
		return
	}

	fn, err := obj.Get("get")
	if err != nil {
		t.Error(err)
		return
	}

	if !fn.IsFunction() {
		t.Error("should be function")
		return
	}
}

func newV8ContextWithFetch(opt ...Option) (*v8go.Context, error) {
	iso := v8go.NewIsolate()
	global := v8go.NewObjectTemplate(iso)

	if err := InjectTo(iso, global, opt...); err != nil {
		return nil, err
	}

	return v8go.NewContext(iso, global), nil
}
