package fetch

import (
	"fmt"
	"testing"
	"time"

	"rogchap.com/v8go"
)

func TestInjectTo(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)

	if err := InjectTo(iso, global); err != nil {
		t.Errorf("error when inject fetch polyfill, %s", err)
		return
	}

	ctx, err := v8go.NewContext(iso, global)
	if err != nil {
		t.Errorf("create context failed: %s", err)
		return
	}

	val, err := ctx.RunScript("fetch('https://www.example.com')", "fetch_example.js")
	if err != nil {
		t.Errorf("failed to do fetch test: %s", err)
		return
	}

	pro, err := val.AsPromise()
	if err != nil {
		t.Errorf("can't convert to promise object: %s", err)
		return
	}

	done := make(chan bool, 1)
	go func() {
		for pro.State() == v8go.Pending {
			continue
		}

		done <- true
	}()

	select {
	case <-time.After(time.Second * 10):
		t.Errorf("request timeout")
		return
	case <-done:
		stat := pro.State()
		if stat == v8go.Rejected {
			fmt.Printf("reject with error: %s\n", pro.Result().String())
		}

		if pro.State() != v8go.Fulfilled {
			t.Errorf("should fetch success, but not")
			return
		}
	}

	obj, err := pro.Result().AsObject()
	if err != nil {
		t.Errorf("can't convert fetch result to object, %s", err)
		return
	}

	ok, err := obj.Get("ok")
	if err != nil {
		t.Errorf("get object 'ok' failed: %s", err)
		return
	}

	if !ok.Boolean() {
		t.Error("should be ok, but not")
	}
}
