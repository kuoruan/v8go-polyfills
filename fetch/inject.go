package fetch

import (
	"fmt"

	"rogchap.com/v8go"
)

func InjectTo(iso *v8go.Isolate, global *v8go.ObjectTemplate, opt ...Option) error {
	f := NewFetcher(opt...)
	if err := global.Set("fetch", f.GetFetchFunctionCallback()); err != nil {
		return fmt.Errorf("v8go-polyfills/fetch: %w", err)
	}

	return nil
}
