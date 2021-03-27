package fetch

import "testing"

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
