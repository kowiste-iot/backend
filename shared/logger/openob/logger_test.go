package openob

import (
	"backend/shared/logger"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				ServiceName:   "test-service",
				Environment:   "test",
				Endpoint:      "localhost:8080",
				Headers:       "test-headers",
				TenantID:      "default",
				StreamName:    "test-stream",
				MinLevel:      logger.InfoLevel,
				ConsoleOutput: true,
				EnableTracing: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLogger(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.config.ServiceName, got.serviceName)
				assert.Equal(t, tt.config.Environment, got.environment)
				assert.Equal(t, tt.config.Endpoint, got.endpoint)
				assert.Equal(t, tt.config.MinLevel, got.minLevel)
			}
		})
	}
}

func TestShouldLog(t *testing.T) {
	log := &Logger{
		minLevel: logger.InfoLevel,
	}

	tests := []struct {
		name  string
		level logger.Level
		want  bool
	}{
		{
			name:  "debug level when min is info",
			level: logger.DebugLevel,
			want:  false,
		},
		{
			name:  "info level when min is info",
			level: logger.InfoLevel,
			want:  true,
		},
		{
			name:  "error level when min is info",
			level: logger.ErrorLevel,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := log.shouldLog(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToMap(t *testing.T) {
	type testStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	tests := []struct {
		name  string
		input interface{}
		want  map[string]interface{}
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name: "map input",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
			},
			want: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
			},
		},
		{
			name: "struct input",
			input: testStruct{
				Field1: "test",
				Field2: 123,
			},
			want: map[string]interface{}{
				"field1": "test",
				"field2": float64(123), // JSON marshaling converts integers to float64
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToMap(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatFields(t *testing.T) {
	tests := []struct {
		name        string
		fields      map[string]interface{}
		want        string
		containsStr string
	}{
		{
			name:   "nil fields",
			fields: nil,
			want:   "",
		},
		{
			name:   "empty fields",
			fields: map[string]interface{}{},
			want:   "",
		},
		{
			name: "single field",
			fields: map[string]interface{}{
				"key": "value",
			},
			want: "> key:value <",
		},
		{
			name: "multiple fields",
			fields: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
			},
			containsStr: "key1:value1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFields(tt.fields)
			if tt.containsStr != "" {
				assert.Contains(t, got, tt.containsStr)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestLoggerMethods(t *testing.T) {
	var receivedRequests []*http.Request
	var receivedBodies []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := new(strings.Builder)
		_, err := io.Copy(body, r.Body)
		require.NoError(t, err)

		receivedRequests = append(receivedRequests, r)
		receivedBodies = append(receivedBodies, body.String())
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := Config{
		ServiceName:   "test-service",
		Environment:   "test",
		Endpoint:      strings.TrimPrefix(server.URL, "http://"),
		Headers:       "test-headers",
		TenantID:      "default",
		StreamName:    "test-stream",
		MinLevel:      logger.DebugLevel,
		ConsoleOutput: true,
		EnableTracing: true,
	}

	logger, err := NewLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)

	ctx := context.Background()
	testFields := map[string]interface{}{
		"test_key": "test_value",
	}

	tests := []struct {
		name     string
		logFunc  func()
		level    string
		message  string
		hasError bool
	}{
		{
			name: "info log",
			logFunc: func() {
				logger.Info(ctx, "info message", testFields)
			},
			level:    "info",
			message:  "info message",
			hasError: false,
		},
		{
			name: "error log",
			logFunc: func() {
				logger.Error(ctx, fmt.Errorf("test error"), "error message", testFields)
			},
			level:    "error",
			message:  "error message",
			hasError: true,
		},
		{
			name: "debug log",
			logFunc: func() {
				logger.Debug(ctx, "debug message", testFields)
			},
			level:    "debug",
			message:  "debug message",
			hasError: false,
		},
		{
			name: "warn log",
			logFunc: func() {
				logger.Warn(ctx, "warn message", testFields)
			},
			level:    "warn",
			message:  "warn message",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialReqCount := len(receivedRequests)
			tt.logFunc()
			time.Sleep(100 * time.Millisecond) // Small delay to ensure log processing

			// Verify request was made
			assert.Equal(t, initialReqCount+1, len(receivedRequests), "Expected one new request")
			req := receivedRequests[len(receivedRequests)-1]

			// Verify request headers
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, fmt.Sprintf("Basic %s", cfg.Headers), req.Header.Get("Authorization"))

			// Verify request body
			body := receivedBodies[len(receivedBodies)-1]
			assert.Contains(t, body, tt.message)
			assert.Contains(t, body, tt.level)
			assert.Contains(t, body, "test_key")
			assert.Contains(t, body, "test_value")

			if tt.hasError {
				assert.Contains(t, body, "test error")
			}
		})
	}
}

func TestLoggerWithoutTracing(t *testing.T) {
	cfg := Config{
		ServiceName:   "test-service",
		Environment:   "test",
		Endpoint:      "localhost:8080",
		Headers:       "test-headers",
		TenantID:      "default",
		StreamName:    "test-stream",
		MinLevel:      logger.DebugLevel,
		ConsoleOutput: true,
		EnableTracing: false,
	}

	logger, err := NewLogger(cfg)
	require.NoError(t, err)

	ctx := context.Background()
	// This shouldn't make any HTTP requests since tracing is disabled
	logger.Info(ctx, "test message", nil)
}
