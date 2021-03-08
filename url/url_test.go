package url

import (
	"testing"

	"rogchap.com/v8go"
)

func TestInject(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()

	if err := Inject(ctx); err != nil {
		t.Errorf("inject url polyfill: %v", err)
	}

	if val, _ := ctx.RunScript("typeof URL", ""); val.String() != "function" {
		t.Error("inject URL failed")
	}

	if val, _ := ctx.RunScript("typeof URLSearchParams", ""); val.String() != "function" {
		t.Error("inject URLSearchParams failed")
	}

	if val, _ := ctx.RunScript("new URLSearchParams('?a=1').get('a')", ""); val.String() != "1" {
		t.Error("test URLSearchParams failed")
	}
}
