/*
Copyright (c) 2025 Dell Inc, or its subsidiaries.

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

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPollImmediateWithContext_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Condition function that returns true immediately
	condition := func(context.Context) (bool, error) {
		return true, nil
	}

	err := ImmediateWithContext(ctx, 1*time.Second, 3*time.Second, condition)
	assert.NoError(t, err)
}

func TestPollImmediateWithContext_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Condition function that never returns true
	condition := func(context.Context) (bool, error) {
		return false, nil
	}

	err := ImmediateWithContext(ctx, 1*time.Second, 2*time.Second, condition)
	assert.Error(t, err)
	assert.Equal(t, ErrWaitTimeout, err)
}

func TestWaitForWithContext_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait function that uses a short interval
	wait := func(ctx context.Context) <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					ch <- struct{}{}
				case <-ctx.Done():
					close(ch)
					return
				}
			}
		}()
		return ch
	}

	// Condition function that returns true after a short delay
	condition := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	err := WaitForWithContext(ctx, wait, condition)
	assert.NoError(t, err)
}

func TestWaitForWithContext_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Wait function that uses a short interval
	wait := func(ctx context.Context) <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					ch <- struct{}{}
				case <-ctx.Done():
					close(ch)
					return
				}
			}
		}()
		return ch
	}

	// Condition function that never returns true
	condition := func(context.Context) (bool, error) {
		return false, nil
	}

	err := WaitForWithContext(ctx, wait, condition)
	assert.Error(t, err)
	assert.Equal(t, ErrWaitTimeout, err)
}

func TestPollImmediateWithContext_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Condition function that returns an error
	condition := func(context.Context) (bool, error) {
		return false, errors.New("condition error")
	}

	err := ImmediateWithContext(ctx, 1*time.Second, 3*time.Second, condition)
	assert.Error(t, err)
	assert.Equal(t, "condition error", err.Error())
}

func TestWaitForWithContext_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait function that uses a short interval
	wait := func(ctx context.Context) <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					ch <- struct{}{}
				case <-ctx.Done():
					close(ch)
					return
				}
			}
		}()
		return ch
	}

	// Condition function that returns an error
	condition := func(context.Context) (bool, error) {
		return false, errors.New("condition error")
	}

	err := WaitForWithContext(ctx, wait, condition)
	assert.Error(t, err)
	assert.Equal(t, "condition error", err.Error())
}
