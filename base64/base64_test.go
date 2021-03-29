package base64

import (
	"testing"

	"rogchap.com/v8go"
)

func TestAtob(t *testing.T) {
	ctx, err := newV8goContext()
	if err != nil {
		t.Error(err)
		return
	}

	val, err := ctx.RunScript("atob()", "atob_undefined.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "" {
		t.Errorf("assert '' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("atob('')", "atob_empty.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "" {
		t.Errorf("assert '' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("atob('MTIzNA==')", "atob_1234.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "1234" {
		t.Errorf("assert '1234' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("atob('5rGJ5a2X')", "atob_unicode.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "汉字" {
		t.Errorf("assert '汉字' but got '%s'", s)
		return
	}
}

func TestBtoa(t *testing.T) {
	ctx, err := newV8goContext()
	if err != nil {
		t.Error(err)
		return
	}

	val, err := ctx.RunScript("btoa()", "btoa_undefined.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "" {
		t.Errorf("assert '' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("btoa('')", "atob_empty.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "" {
		t.Errorf("assert '' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("btoa('1234')", "btoa_1234.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "MTIzNA==" {
		t.Errorf("assert 'MTIzNA==' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("btoa('汉字')", "btoa_unicode.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "5rGJ5a2X" {
		t.Errorf("assert '5rGJ5a2X' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("btoa({})", "btoa_object.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "W29iamVjdCBPYmplY3Rd" {
		t.Errorf("assert 'W29iamVjdCBPYmplY3Rd' but got '%s'", s)
		return
	}
}

func newV8goContext() (*v8go.Context, error) {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)

	if err := InjectTo(iso, global); err != nil {
		return nil, err
	}

	return v8go.NewContext(iso, global)
}
