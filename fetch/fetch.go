package fetch

import (
	_ "embed"

	"go.kuoruan.net/v8go-polyfills/internal"

	"rogchap.com/v8go"
)

//go:embed bundle.js
var fetchPolyfill string

func Inject(ctx internal.Context, opt ...Option) error {
	if ctx == nil {
		panic("ctx is required")
	}

	f := &fetcher{}

	for _, o := range opt {
		o.apply(f)
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return err
	}
	obj := ctx.Global()

	fetchFn, _ := v8go.NewFunctionTemplate(iso, f.goFetchSync)
	if err := obj.Set("_goFetchSync", fetchFn); err != nil {
		return err
	}

	_, err = ctx.RunScript(fetchPolyfill, "fetch.js")
	return err
}
