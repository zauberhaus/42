/*
Copyright © 2021 Dirk Lembke <dirk@lembke.nz>

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

package background_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/42/background"
	"github.com/zauberhaus/42/logger"
)

const (
	tick = "tick %v"
)

func TestBackgroundProcessFinish(t *testing.T) {
	val := 0

	logger := logger.NewZapLogger()

	process := background.Process{}
	process.Init(t.Name(), nil, nil, logger)

	p := func(ctx context.Context) (bool, error) {
		return exec(ctx, 10*time.Millisecond, func(p ...interface{}) (bool, error) {
			val := p[0].(*int)

			logger.Infof(tick, time.Now())
			*val = *val + 1

			return true, nil
		}, &val)
	}

	c := process.Run(p)
	<-c

	err := <-process.Done()
	assert.NoError(t, err)

	assert.Equal(t, 1, val)
}

func TestBackgroundProcessFinishAndStop(t *testing.T) {
	val := 0

	logger := logger.NewZapLogger()

	process := background.Process{}
	process.Init(t.Name(), nil, nil, logger)

	p := func(ctx context.Context) (bool, error) {
		return exec(ctx, 10*time.Millisecond, func(p ...interface{}) (bool, error) {
			val := p[0].(*int)

			logger.Infof(tick, time.Now())
			*val = *val + 1

			return true, nil
		}, &val)
	}

	<-process.Run(p)

	time.Sleep(50 * time.Millisecond)

	process.Stop(context.Background())

	assert.Equal(t, 1, val)
}

func TestBackgroundProcessCancel(t *testing.T) {
	val := 0
	started := make(chan bool)

	logger := logger.NewZapLogger()

	process := background.Process{}
	process.Init(t.Name(), nil, nil, logger)

	p := func(ctx context.Context) (bool, error) {
		close(started)
		return exec(ctx, 10*time.Millisecond, func(p ...interface{}) (bool, error) {
			val := p[0].(*int)

			logger.Infof(tick, time.Now())
			*val = *val + 1

			return false, nil
		}, &val)
	}

	process.Run(p)
	<-started

	time.Sleep(220 * time.Millisecond)

	process.Stop(context.Background())
}

func TestBackgroundProcessInitShutdown(t *testing.T) {
	val := 0

	logger := logger.NewZapLogger()

	i := func(ctx context.Context) error {
		val = 1

		return nil
	}

	s := func(ctx context.Context) error {
		val += 10

		return nil
	}

	process := background.Process{}
	process.Init(t.Name(), i, s, logger)

	p := func(ctx context.Context) (bool, error) {
		return exec(ctx, 10*time.Millisecond, func(p ...interface{}) (bool, error) {
			val := p[0].(*int)

			logger.Infof(tick, time.Now())
			*val = *val + 1

			return true, nil
		}, &val)
	}

	c := process.Run(p)
	<-c

	assert.Equal(t, 1, val)

	err := <-process.Done()
	assert.NoError(t, err)

	assert.Equal(t, 12, val)
}

func TestBackgroundProcessInitFailed(t *testing.T) {
	val := 0

	logger := logger.NewZapLogger()

	i := func(ctx context.Context) error {
		val = 99

		return fmt.Errorf("Init failed")
	}

	s := func(ctx context.Context) error {
		val += 10

		return nil
	}

	process := background.Process{}
	process.Init(t.Name(), i, s, logger)

	p := func(ctx context.Context) (bool, error) {
		return exec(ctx, 10*time.Millisecond, func(p ...interface{}) (bool, error) {
			val := p[0].(*int)

			logger.Infof(tick, time.Now())
			*val = *val + 1

			return true, nil
		}, &val)
	}

	c := process.Run(p)
	<-c

	assert.Equal(t, 99, val)

	err := <-process.Done()
	assert.NoError(t, err)

	assert.Equal(t, 99, val)
}

func TestBackgroundProcessShutdownFailed(t *testing.T) {
	val := 0

	logger := logger.NewZapLogger()

	i := func(ctx context.Context) error {
		val = 99
		return nil
	}

	s := func(ctx context.Context) error {
		val += 10

		return fmt.Errorf("Shutdown failed")
	}

	process := background.Process{}
	process.Init(t.Name(), i, s, logger)

	p := func(ctx context.Context) (bool, error) {
		return exec(ctx, 10*time.Millisecond, func(p ...interface{}) (bool, error) {
			val := p[0].(*int)

			logger.Infof(tick, time.Now())
			*val = *val + 1

			return true, nil
		}, &val)
	}

	c := process.Run(p)
	<-c

	assert.Equal(t, 99, val)

	err := <-process.Done()
	assert.NoError(t, err)

	assert.Equal(t, 110, val)
}

func exec(ctx context.Context, timeout time.Duration, f func(...interface{}) (bool, error), param ...interface{}) (bool, error) {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	for {
		timer.Reset(timeout)

		select {
		case <-ctx.Done():
			return false, nil
		case <-timer.C:
			return f(param...)
		}
	}
}
