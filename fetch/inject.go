package fetch

import (
	"fmt"

	"rogchap.com/v8go"
)

func InjectTo(iso *v8go.Isolate, global *v8go.ObjectTemplate, opt ...Option) error {
	f := NewFetcher(opt...)

	fetchFn, err := v8go.NewFunctionTemplate(iso, f.GetFetchFunctionCallback())
	if err != nil {
		return fmt.Errorf("v8go-polyfills/fetch: %w", err)
	}

	if err := global.Set("fetch", fetchFn); err != nil {
		return fmt.Errorf("v8go-polyfills/fetch: %w", err)
	}

	return nil
}
