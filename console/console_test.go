package console

import (
	"testing"

	"rogchap.com/v8go"
)

func TestInject(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()

	if err := Inject(ctx); err != nil {
		t.Error(err)
	}

	if _, err := ctx.RunScript("console.log(1111)", ""); err != nil {
		t.Error(err)
	}
}
