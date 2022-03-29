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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zauberhaus/42/logger"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/stretchr/testify/assert"
)

type header struct {
	Key   string
	Value string
}

func TestJSONLogMiddleware(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	assert.NotNil(t, logs)

	router := gin.New()
	router.Use(logger.JSONLogMiddleware())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "{}")
	})

	w := PerformRequest(router, "GET", "/")

	assert.Equal(t, http.StatusOK, w.Code)

	all := logs.All()

	assert.Len(t, all, 1)

	entry := all[0]

	assert.NotNil(t, entry)
	assert.Equal(t, zapcore.InfoLevel, entry.Level)

	context := entry.ContextMap()

	assert.Contains(t, context, "client_ip")
	assert.Contains(t, context, "duration")
	assert.Contains(t, context, "method")
	assert.Contains(t, context, "path")
	assert.Contains(t, context, "status")

	assert.Equal(t, "GET", context["method"])
	assert.Equal(t, "/", context["path"])
	assert.Equal(t, int64(http.StatusOK), context["status"])
}

func TestJSONLogMiddlewareError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	logs := logger.Observe(logger.DebugLevel).(*observer.ObservedLogs)
	assert.NotNil(t, logs)

	router := gin.New()
	router.Use(logger.JSONLogMiddleware())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusInsufficientStorage, "{}")
	})

	w := PerformRequest(router, "GET", "/")

	assert.Equal(t, http.StatusInsufficientStorage, w.Code)

	all := logs.All()

	assert.Len(t, all, 1)

	entry := all[0]

	assert.NotNil(t, entry)
	assert.Equal(t, zapcore.ErrorLevel, entry.Level)

	context := entry.ContextMap()

	assert.Contains(t, context, "client_ip")
	assert.Contains(t, context, "duration")
	assert.Contains(t, context, "method")
	assert.Contains(t, context, "path")
	assert.Contains(t, context, "status")

	assert.Equal(t, "GET", context["method"])
	assert.Equal(t, "/", context["path"])
	assert.Equal(t, int64(http.StatusInsufficientStorage), context["status"])
}

func PerformRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
