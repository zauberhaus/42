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

package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type ZapLogger struct {
	*zap.SugaredLogger
	level Level
}

func NewZapLogger(options ...interface{}) Logger {
	skip := 0
	level := InfoLevel

	for _, o := range options {
		switch v := o.(type) {
		case Level:
			level = v
		case int:
			skip = v
		}
	}

	return NewZapLoggerWithOptions(
		level,
		zap.AddCallerSkip(skip),
	)
}

func NewZapLoggerWithOptions(level Level, opts ...zap.Option) Logger {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(level)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(opts...)
	if err != nil {
		log.Fatal(err)
	}

	return &ZapLogger{
		logger.Sugar(),
		level,
	}
}

func NewObservableZapLogger(level Level) (Logger, *observer.ObservedLogs) {
	return NewObservableZapLoggerWithOptions(level,
		zap.AddCaller(),
	)
}

func NewObservableZapLoggerWithOptions(level Level, opts ...zap.Option) (Logger, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.Level(level))
	logger := zap.New(core, opts...)

	return &ZapLogger{
		logger.Sugar(),
		level,
	}, recorded
}

func (z *ZapLogger) With(args ...interface{}) Logger {
	return &ZapLogger{
		z.SugaredLogger.With(args...),
		z.level,
	}
}

func (z *ZapLogger) GetLevel() int8 {
	return int8(z.level)
}

func (z *ZapLogger) GetLevelMap() map[int][]string {
	return map[int][]string{
		int(zapcore.DebugLevel): {"debug"},
		int(zapcore.InfoLevel):  {"info"},
		int(zapcore.WarnLevel):  {"warning", "warn"},
		int(zapcore.ErrorLevel): {"error"},
		int(zapcore.PanicLevel): {"panic"},
		int(zapcore.FatalLevel): {"fatal"},
	}
}

func (z *ZapLogger) GetLevelNames() []string {
	m := z.GetLevelMap()

	list := []string{}
	for _, values := range m {
		for _, v := range values {
			list = append(list, v)
		}
	}

	return list
}
