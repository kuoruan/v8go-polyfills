package base64

import (
	"fmt"

	"rogchap.com/v8go"
)

func InjectTo(iso *v8go.Isolate, global *v8go.ObjectTemplate) error {
	b := NewBase64()

	for _, f := range []struct {
		Name string
		Func func() v8go.FunctionCallback
	}{
		{Name: "atob", Func: b.GetAtobFunctionCallback},
		{Name: "btoa", Func: b.GetBtoaFunctionCallback},
	} {
		fn, err := v8go.NewFunctionTemplate(iso, f.Func())
		if err != nil {
			return fmt.Errorf("v8go-polyfills/fetch: %w", err)
		}

		if err := global.Set(f.Name, fn); err != nil {
			return fmt.Errorf("v8go-polyfills/fetch: %w", err)
		}
	}

	return nil
}
