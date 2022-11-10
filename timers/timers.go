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
	"context"
	"errors"

	"go.kuoruan.net/v8go-polyfills/timers/internal"
	"rogchap.com/v8go"
)

type Timers interface {
	GetSetTimeoutFunctionCallback() v8go.FunctionCallback
	GetSetIntervalFunctionCallback() v8go.FunctionCallback

	GetClearTimeoutFunctionCallback() v8go.FunctionCallback
	GetClearIntervalFunctionCallback() v8go.FunctionCallback
}

type timers struct {
	Items      map[int32]*internal.Item
	NextItemID int32
	ctx        context.Context
}

const initNextItemID = 1

func NewTimers(ctx context.Context) Timers {
	return &timers{
		Items:      make(map[int32]*internal.Item),
		NextItemID: initNextItemID,
		ctx:        ctx,
	}
}

func (t *timers) GetSetTimeoutFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()

		id, err := t.startNewTimer(t.ctx, info.This(), info.Args(), false)
		if err != nil {
			return newInt32Value(ctx, 0)
		}

		return newInt32Value(ctx, id)
	}
}

func (t *timers) GetSetIntervalFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()

		id, err := t.startNewTimer(t.ctx, info.This(), info.Args(), true)
		if err != nil {
			return newInt32Value(ctx, 0)
		}

		return newInt32Value(ctx, id)
	}
}

func (t *timers) GetClearTimeoutFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 && args[0].IsInt32() {
			t.clear(args[0].Int32(), false)
		}

		return nil
	}
}

func (t *timers) GetClearIntervalFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 && args[0].IsInt32() {
			t.clear(args[0].Int32(), true)
		}

		return nil
	}
}

func (t *timers) clear(id int32, interval bool) {
	if id < initNextItemID {
		return
	}

	if item, ok := t.Items[id]; ok && item.Interval == interval {
		item.Clear()
	}
}

func (t *timers) startNewTimer(ctx context.Context, this v8go.Valuer, args []*v8go.Value, interval bool) (int32, error) {
	if len(args) <= 0 {
		return 0, errors.New("1 argument required, but only 0 present")
	}

	fn, err := args[0].AsFunction()
	if err != nil {
		return 0, err
	}

	var delay int32
	if len(args) > 1 && args[1].IsInt32() {
		delay = args[1].Int32()
	}
	if delay < 10 {
		delay = 10
	}

	var restArgs []v8go.Valuer
	if len(args) > 2 {
		restArgs = make([]v8go.Valuer, 0)
		for _, arg := range args[2:] {
			restArgs = append(restArgs, arg)
		}
	}

	item := &internal.Item{
		ID:       t.NextItemID,
		Done:     false,
		Cleared:  false,
		Delay:    delay,
		Interval: interval,
		FunctionCB: func() {
			_, _ = fn.Call(this, restArgs...)
		},
		ClearCB: func(id int32) {
			delete(t.Items, id)
		},
	}

	t.NextItemID++
	t.Items[item.ID] = item

	item.Start(ctx)

	return item.ID, nil
}

func newInt32Value(ctx *v8go.Context, i int32) *v8go.Value {
	iso := ctx.Isolate()
	v, _ := v8go.NewValue(iso, i)
	return v
}
