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

func newV8ContextWithFetch(opt ...Option) (*v8go.Context, error) {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)

	if err := InjectTo(iso, global, opt...); err != nil {
		return nil, err
	}

	return v8go.NewContext(iso, global)
}
