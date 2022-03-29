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

package logger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/42/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestDebug(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	logger.Debug("test")
	logger.Debug("test", 1)

	expected := []string{
		"test",
		"test1",
	}

	for idx, l := range logs.All() {
		assert.Equal(t, expected[idx], l.Message)
		assert.Equal(t, zapcore.DebugLevel, l.Level)
	}
}

func TestDebugf(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	logger.Debugf("test %v", 1)
	logger.Debugf("test %v", 1, 2)

	expected := []string{
		"test 1",
		"test 1%!(EXTRA int=2)",
	}

	for idx, l := range logs.All() {
		assert.Equal(t, expected[idx], l.Message)
		assert.Equal(t, zapcore.DebugLevel, l.Level)
	}
}

func TestInfo(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	logger.Info("test")
	logger.Info("test", 1)

	expected := []string{
		"test",
		"test1",
	}

	for idx, l := range logs.All() {
		assert.Equal(t, expected[idx], l.Message)
		assert.Equal(t, zapcore.InfoLevel, l.Level)
	}
}

func TestInfof(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	logger.Infof("test %v", 1)
	logger.Infof("test %v", 1, 2)

	expected := []string{
		"test 1",
		"test 1%!(EXTRA int=2)",
	}

	for idx, l := range logs.All() {
		assert.Equal(t, expected[idx], l.Message)
		assert.Equal(t, zapcore.InfoLevel, l.Level)
	}
}

func TestWarn(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	logger.Warn("test")
	logger.Warn("test", 1)

	expected := []string{
		"test",
		"test1",
	}

	for idx, l := range logs.All() {
		assert.Equal(t, expected[idx], l.Message)
		assert.Equal(t, zapcore.WarnLevel, l.Level)
	}
}

func TestWarnf(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	logger.Warnf("test %v", 1)
	logger.Warnf("test %v", 1, 2)

	expected := []string{
		"test 1",
		"test 1%!(EXTRA int=2)",
	}

	for idx, l := range logs.All() {
		assert.Equal(t, expected[idx], l.Message)
		assert.Equal(t, zapcore.WarnLevel, l.Level)
	}
}

func TestError(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	logger.Error("test")
	logger.Error("test", 1)

	expected := []string{
		"test",
		"test1",
	}

	for idx, l := range logs.All() {
		assert.Equal(t, expected[idx], l.Message)
		assert.Equal(t, zapcore.ErrorLevel, l.Level)
	}
}

func TestErrorf(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	logger.Errorf("test %v", 1)
	logger.Errorf("test %v", 1, 2)

	expected := []string{
		"test 1",
		"test 1%!(EXTRA int=2)",
	}

	for idx, l := range logs.All() {
		assert.Equal(t, expected[idx], l.Message)
		assert.Equal(t, zapcore.ErrorLevel, l.Level)
	}
}

func TestWith(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)

	l := logger.With(zap.String("client_ip", "127.0.0.1"))

	l.Info("test")

	assert.Len(t, logs.All(), 1)

	entry := logs.All()[0]
	context := entry.ContextMap()

	assert.Equal(t, "test", entry.Message)
	assert.Contains(t, context, "client_ip")
	assert.Equal(t, "127.0.0.1", context["client_ip"])
}

func TestSetLogger(t *testing.T) {
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	l := logger.With(zap.String("client_ip", "127.0.0.1"))
	logger.SetLogger(l)
	logger.Info("test")

	assert.Len(t, logs.All(), 1)

	entry := logs.All()[0]
	context := entry.ContextMap()

	assert.Equal(t, "test", entry.Message)
	assert.Contains(t, context, "client_ip")
	assert.Equal(t, "127.0.0.1", context["client_ip"])
}

func TestGetLogger(t *testing.T) {
	l := logger.GetLogger()
	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	l.Info("test")
	assert.Len(t, logs.All(), 0)
}

func TestZapGetLevelMap(t *testing.T) {
	l := logger.NewZapLogger(logger.InfoLevel, 0)
	m := l.GetLevelMap()
	assert.Len(t, m, 6)

}

func TestZapGetLevelNames(t *testing.T) {
	l := logger.NewZapLogger(logger.InfoLevel, 0)
	m := l.GetLevelNames()
	assert.Len(t, m, 7)

}
