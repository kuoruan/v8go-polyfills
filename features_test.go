package polyfills

import (
	"go.kuoruan.net/v8go-polyfills/console"
	"rogchap.com/v8go"
	"testing"

	_ "embed"
)

//go:embed features-test.js
var test string

func TestFeatures(t *testing.T) {
	iso := v8go.NewIsolate()
	ctx := v8go.NewContext(iso)

	if err := console.InjectTo(ctx); err != nil {
		t.Error(err)
		return
	}

	if _, err := ctx.RunScript(test, "features-test.js"); err != nil {
		t.Error(err)
	}
}
