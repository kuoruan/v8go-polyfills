/*
 * Copyright (c) 2021 Xingwang Liao
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package fetch

import (
	"fmt"
	"testing"
	"time"

	"rogchap.com/v8go"
)

func TestInjectTo(t *testing.T) {
	t.Parallel()

	iso := v8go.NewIsolate()
	global := v8go.NewObjectTemplate(iso)

	if err := InjectTo(iso, global); err != nil {
		t.Errorf("error when inject fetch polyfill, %s", err)
		return
	}

	ctx := v8go.NewContext(iso, global)

	val, err := ctx.RunScript("fetch('https://www.example.com')", "fetch_example.js")
	if err != nil {
		t.Errorf("failed to do fetch test: %s", err)
		return
	}

	pro, err := val.AsPromise()
	if err != nil {
		t.Errorf("can't convert to promise object: %s", err)
		return
	}

	done := make(chan bool, 1)
	go func() {
		for pro.State() == v8go.Pending {
			continue
		}

		done <- true
	}()

	select {
	case <-time.After(time.Second * 10):
		t.Errorf("request timeout")
		return
	case <-done:
		stat := pro.State()
		if stat == v8go.Rejected {
			fmt.Printf("reject with error: %s\n", pro.Result().String())
		}

		if pro.State() != v8go.Fulfilled {
			t.Errorf("should fetch success, but not")
			return
		}
	}

	obj, err := pro.Result().AsObject()
	if err != nil {
		t.Errorf("can't convert fetch result to object, %s", err)
		return
	}

	ok, err := obj.Get("ok")
	if err != nil {
		t.Errorf("get object 'ok' failed: %s", err)
		return
	}

	if !ok.Boolean() {
		t.Error("should be ok, but not")
	}
}
