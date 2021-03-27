package console

import (
	"os"
	"testing"

	"rogchap.com/v8go"
)

func TestInject(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)

	if err := InjectTo(iso, global, WithOutput(os.Stdout)); err != nil {
		t.Error(err)
	}

	ctx, _ := v8go.NewContext()

	if _, err := ctx.RunScript("console.log(1111)", ""); err != nil {
		t.Error(err)
	}
}
