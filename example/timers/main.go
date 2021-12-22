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

package main

import (
	"fmt"
	"time"

	"go.kuoruan.net/v8go-polyfills/timers"
	"rogchap.com/v8go"
)

func main() {
	iso := v8go.NewIsolate()
	global := v8go.NewObjectTemplate(iso)

	if err := timers.InjectTo(iso, global); err != nil {
		panic(err)
	}

	ctx := v8go.NewContext(iso, global)

	val, err := ctx.RunScript(
		"new Promise((resolve) => setTimeout(function(name) {resolve(`Hello, ${name}!`)}, 1000, 'Tom'))",
		"resolve.js",
	)
	if err != nil {
		panic(err)
	}

	proms, err := val.AsPromise()
	if err != nil {
		panic(err)
	}

	done := make(chan bool, 1)

	go func() {
		for proms.State() == v8go.Pending {
			continue
		}

		done <- proms.State() == v8go.Fulfilled
	}()

	select {
	case succ := <-done:
		if !succ {
			panic("except success but not")
		}

		fmt.Println(proms.Result().String())
	case <-time.After(time.Second * 2):
		panic("timeout")
	}
}
