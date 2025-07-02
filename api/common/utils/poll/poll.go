/*
Copyright (c) 2022 Dell Inc, or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package poll

// Polling functionality inspired by https://github.com/kubernetes/apimachinery

import (
	"context"
	"errors"
	"time"
)

var ErrWaitTimeout = errors.New("timed out waiting for the condition")

type (
	WaitWithContextFunc      func(ctx context.Context) <-chan struct{}
	ConditionWithContextFunc func(context.Context) (done bool, err error)
)

func PollImmediateWithContext(ctx context.Context, interval, timeout time.Duration, condition ConditionWithContextFunc) error {
	return poll(ctx, true, poller(interval, timeout), condition)
}

func WaitForWithContext(ctx context.Context, wait WaitWithContextFunc, fn ConditionWithContextFunc) error {
	waitCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := wait(waitCtx)
	for {
		select {
		case _, open := <-c:
			ok, err := fn(ctx)
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
			if !open {
				return ErrWaitTimeout
			}
		case <-ctx.Done():
			return ErrWaitTimeout
		}
	}
}

func poll(ctx context.Context, immediate bool, wait WaitWithContextFunc, condition ConditionWithContextFunc) error {
	if immediate {
		done, err := condition(ctx)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}

	select {
	case <-ctx.Done():
		return ErrWaitTimeout
	default:
		return WaitForWithContext(ctx, wait, condition)
	}
}

func poller(interval, timeout time.Duration) WaitWithContextFunc {
	return func(ctx context.Context) <-chan struct{} {
		ch := make(chan struct{})

		go func() {
			defer close(ch)

			tick := time.NewTicker(interval)
			defer tick.Stop()

			var after <-chan time.Time
			if timeout != 0 {
				timer := time.NewTimer(timeout)
				after = timer.C
				defer timer.Stop()
			}

			for {
				select {
				case <-tick.C:
					select {
					case ch <- struct{}{}:
					default:
					}
				case <-after:
					return
				case <-ctx.Done():
					return
				}
			}
		}()

		return ch
	}
}
