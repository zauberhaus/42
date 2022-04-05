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

package background

import (
	"golang.org/x/net/context"
)

type Process struct {
	name   string
	ctx    context.Context
	cancel context.CancelFunc
	done   chan error
	init   func(ctx context.Context) error
	close  func(ctx context.Context) error
	logger Logger
}

func (p *Process) Init(name string, init func(ctx context.Context) error, close func(ctx context.Context) error, logger Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	p.ctx = ctx
	p.cancel = cancel
	p.name = name
	p.done = make(chan error, 1)
	p.init = init
	p.close = close
	p.logger = logger
}

func (p *Process) Run(process func(ctx context.Context) (bool, error)) chan error {

	rc := make(chan error, 1)

	go func() {
		if p.init != nil {
			p.logger.Infof("Init %s", p.name)
			if err := p.init(p.ctx); err != nil {
				p.logger.Errorf("Init %s failed: %v", p.name, err)
				rc <- err
				close(p.done)
				close(rc)
				return
			}
		}
		p.logger.Infof("%s started", p.name)
		close(rc)
		if process != nil {
			finished, err := process(p.ctx)
			if err != nil {
				p.logger.Errorf("Process %s failed: %v", p.name, err)
				p.done <- err
				close(p.done)
			}

			if finished {
				p.logger.Debugf("%s finished", p.name)
			} else {
				p.logger.Debugf("%s cancelled", p.name)
			}

			if p.close != nil {
				p.logger.Debugf("%s wait for shutdown", p.name)

				if err := p.close(p.ctx); err == nil {
					p.logger.Debugf("%s shutdown finished.", p.name)

				} else {
					p.logger.Errorf("%s shutdown failed: %v", p.name, err)
				}
			}

			close(p.done)
		}
	}()

	return rc
}

func (p *Process) Stop(ctx context.Context) error {
	p.logger.Infof("%s shutdown start", p.name)
	p.cancel()
	p.logger.Debugf("%s wait for shutdown done", p.name)

	var err error

	select {
	case err = <-p.done:
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err == nil {
		p.logger.Infof("%s done", p.name)
	} else {
		p.logger.Errorf("%s stop failed: %v", p.name, err)
	}
	return err
}

func (p *Process) Done() chan error {
	return p.done
}
