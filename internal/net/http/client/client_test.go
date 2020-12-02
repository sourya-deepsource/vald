//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package client

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	}

	transportComparator = []comparator.Option{
		comparator.AllowUnexported(transport{}),
		comparator.AllowUnexported(http.Transport{}),
		comparator.IgnoreFields(http.Transport{}, "idleLRU"),

		comparator.Comparer(func(x, y backoff.Option) bool {
			return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
		}),
		comparator.Comparer(func(x, y func(*http.Request) (*url.URL, error)) bool {
			return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
		}),
		comparator.Comparer(func(x, y func(ctx context.Context, network, addr string) (net.Conn, error)) bool {
			return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
		}),
		comparator.Comparer(func(x, y sync.Mutex) bool {
			return reflect.DeepEqual(x, y)
		}),
		comparator.Comparer(func(x, y atomic.Value) bool {
			return reflect.DeepEqual(x.Load(), y.Load())
		}),
		comparator.Comparer(func(x, y sync.Once) bool {
			return reflect.DeepEqual(x, y)
		}),
	}
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want *http.Client
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *http.Client, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *http.Client, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/
		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}