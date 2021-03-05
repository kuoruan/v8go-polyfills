package fetch

import (
	_ "embed"

	polyfills "go.kuoruan.net/v8go-polyfills"

	"rogchap.com/v8go"
)

//go:embed fetch.js
var fetchPolyfill string

func Inject(ctx polyfills.Context, opt ...Option) error {
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
	global := ctx.Global()

	fetchFn, _ := v8go.NewFunctionTemplate(iso, f.goFetchSync)
	if err := global.Set("_goFetchSync", fetchFn); err != nil {
		return err
	}

	_, err = ctx.RunScript(fetchPolyfill, "fetch.js")
	return err
}
