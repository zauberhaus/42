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

import "go.uber.org/zap"

type (
	// Logger is an interface that can be passed to ClientOptions.Logger.
	Logger interface {
		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Panic(args ...interface{})
		Fatal(args ...interface{})

		Debugf(template string, args ...interface{})
		Infof(template string, args ...interface{})
		Warnf(template string, args ...interface{})
		Errorf(template string, args ...interface{})
		Panicf(template string, args ...interface{})
		Fatalf(template string, args ...interface{})

		With(args ...interface{}) Logger
	}
)

var logger = NewZapLogger(
	zap.AddCallerSkip(2),
)

func GetLogger() Logger {
	return logger
}

func SetLogger(l Logger) {
	logger = l
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Warn(args ...interface{}) {
	logger.Warn(args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Panic(args ...interface{}) {
	logger.Panic(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args)
}

func Warnf(template string, args ...interface{}) {
	logger.Warn(template, args)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args)
}

func Fatalf(template string, args ...interface{}) {
	logger.Panicf(template, args)
}

func With(args ...interface{}) Logger {
	return logger.With(args)
}
