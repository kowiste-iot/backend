package mock

import (
	"context"
	"sync"

	"backend/shared/authorization/domain/command"
)

// MockPermissionProvider is a mock implementation of PermissionProvider interface
type MockPermissionProvider struct {
	mu        sync.RWMutex
	calls     []PermissionCall
	responses map[string]mockResponse
}

type PermissionCall struct {
	Token    string
	Resource string
	Action   string
	TenantID string
	BranchID string
}

type mockResponse struct {
	result bool
	err    error
}

// NewMockPermissionProvider creates a new instance of MockPermissionProvider
func NewMockPermissionProvider() *MockPermissionProvider {
	return &MockPermissionProvider{
		responses: make(map[string]mockResponse),
		calls:     make([]PermissionCall, 0),
	}
}

// HasPermission implements the PermissionProvider interface
func (m *MockPermissionProvider) HasPermission(ctx context.Context, input *command.PermissionInput) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := PermissionCall{
		Token:    input.Token,
		Resource: input.Resource,
		Action:   input.Action,
		TenantID: input.TenantID,
		BranchID: input.BranchID,
	}

	m.calls = append(m.calls, call)

	key := m.buildKey(input)
	if response, exists := m.responses[key]; exists {
		return response.result, response.err
	}

	// Default response if no mock is set
	return false, nil
}

// SetResponse sets up the mock to return specific values for given inputs
func (m *MockPermissionProvider) SetResponse(input *command.PermissionInput, result bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := m.buildKey(input)
	m.responses[key] = mockResponse{
		result: result,
		err:    err,
	}
}

// GetCalls returns all recorded calls to HasPermission
func (m *MockPermissionProvider) GetCalls() []PermissionCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]PermissionCall, len(m.calls))
	copy(calls, m.calls)
	return calls
}

// Reset clears all recorded calls and configured responses
func (m *MockPermissionProvider) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = make([]PermissionCall, 0)
	m.responses = make(map[string]mockResponse)
}

// buildKey creates a unique key for storing mock responses
func (m *MockPermissionProvider) buildKey(input *command.PermissionInput) string {
	return input.Token + "|" + input.Resource + "|" + input.Action + "|" + input.TenantID + "|" + input.BranchID
}
