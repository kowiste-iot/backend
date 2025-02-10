package mock

import (
    "context"
    "sync"
)

// LogCall represents a single call to any logging method
type LogCall struct {
    Method string
    Msg    string
    Err    error
    Fields any
}

// MockLogger implements the Logger interface for testing
type MockLogger struct {
    mu    sync.RWMutex
    calls []LogCall
}

// NewMockLogger creates a new instance of MockLogger
func NewMockLogger() *MockLogger {
    return &MockLogger{
        calls: make([]LogCall, 0),
    }
}

// Info implements Logger.Info
func (m *MockLogger) Info(ctx context.Context, msg string, fields any) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.calls = append(m.calls, LogCall{
        Method: "Info",
        Msg:    msg,
        Fields: fields,
    })
}

// Error implements Logger.Error
func (m *MockLogger) Error(ctx context.Context, err error, msg string, fields any) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.calls = append(m.calls, LogCall{
        Method: "Error",
        Msg:    msg,
        Err:    err,
        Fields: fields,
    })
}

// Debug implements Logger.Debug
func (m *MockLogger) Debug(ctx context.Context, msg string, fields any) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.calls = append(m.calls, LogCall{
        Method: "Debug",
        Msg:    msg,
        Fields: fields,
    })
}

// Warn implements Logger.Warn
func (m *MockLogger) Warn(ctx context.Context, msg string, fields any) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.calls = append(m.calls, LogCall{
        Method: "Warn",
        Msg:    msg,
        Fields: fields,
    })
}

// GetCalls returns all recorded calls
func (m *MockLogger) GetCalls() []LogCall {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    calls := make([]LogCall, len(m.calls))
    copy(calls, m.calls)
    return calls
}

// GetCallsByMethod returns all calls for a specific method
func (m *MockLogger) GetCallsByMethod(method string) []LogCall {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    var methodCalls []LogCall
    for _, call := range m.calls {
        if call.Method == method {
            methodCalls = append(methodCalls, call)
        }
    }
    return methodCalls
}

// Reset clears all recorded calls
func (m *MockLogger) Reset() {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.calls = make([]LogCall, 0)
}