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

package timers

import (
	"testing"
	"time"

	"rogchap.com/v8go"
)

func Test_SetTimeout(t *testing.T) {
	ctx, err := newV8ContextWithTimers()
	if err != nil {
		t.Error(err)
		return
	}

	val, err := ctx.RunScript("setTimeout(function() {}, 2000)", "set_timeout.js")
	if err != nil {
		t.Error(err)
		return
	}

	if !val.IsInt32() {
		t.Errorf("except 1 but got %v", val)
		return
	}

	if id := val.Int32(); id != 1 {
		t.Errorf("except 1 but got %d", id)
	}

	time.Sleep(time.Second * 6)
}

func newV8ContextWithTimers() (*v8go.Context, error) {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)

	if err := InjectTo(iso, global); err != nil {
		return nil, err
	}

	return v8go.NewContext(iso, global)
}
