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

package internal

import (
	"time"
)

type FunctionCallback func()

type ClearCallback func(id int32)

type Item struct {
	ID       int32
	Done     bool
	Cleared  bool
	Interval bool
	Delay    int32

	ClearCB    ClearCallback
	FunctionCB FunctionCallback
}

func (t *Item) Clear() {
	if !t.Cleared {
		t.Cleared = true

		if t.ClearCB != nil {
			t.ClearCB(t.ID)
		}
	}

	t.Done = true
}

func (t *Item) Start() {
	go func() {
		defer t.Clear() // self clear

		ticker := time.NewTicker(time.Duration(t.Delay) * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			if t.Done {
				break
			}

			if t.FunctionCB != nil {
				t.FunctionCB()
			}

			if !t.Interval {
				t.Done = true
				break
			}
		}
	}()
}
