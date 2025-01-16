package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to generate UUID v7
func newUUIDv7(t *testing.T) string {
	id, err := uuid.NewV7()
	require.NoError(t, err)
	return id.String()
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		tenantID    string
		branchID    string
		assetName   string
		description string
		wantErr     bool
	}{
		{
			name:        "valid asset",
			tenantID:    newUUIDv7(t),
			branchID:    newUUIDv7(t),
			assetName:   "Test Asset",
			description: "Test Description",
			wantErr:     false,
		},
		{
			name:        "invalid tenant id",
			tenantID:    "invalid-uuid",
			branchID:    newUUIDv7(t),
			assetName:   "Test Asset",
			description: "Test Description",
			wantErr:     true,
		},
		{
			name:        "invalid branch id",
			tenantID:    newUUIDv7(t),
			branchID:    "invalid-uuid",
			assetName:   "Test Asset",
			description: "Test Description",
			wantErr:     true,
		},
		{
			name:        "invalid name - too short",
			tenantID:    newUUIDv7(t),
			branchID:    newUUIDv7(t),
			assetName:   "Te",
			description: "Test Description",
			wantErr:     true,
		},
		{
			name:        "empty name",
			tenantID:    newUUIDv7(t),
			branchID:    newUUIDv7(t),
			assetName:   "",
			description: "Test Description",
			wantErr:     true,
		},
		{
			name:        "description too short",
			tenantID:    newUUIDv7(t),
			branchID:    newUUIDv7(t),
			assetName:   "Valid Name",
			description: "ab",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset, err := New(tt.tenantID, tt.branchID, tt.assetName, tt.description)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, asset.ID())
			assert.Equal(t, tt.tenantID, asset.TenantID())
			assert.Equal(t, tt.branchID, asset.branchName)
			assert.Equal(t, tt.assetName, asset.Name())
			assert.Equal(t, tt.description, asset.Description())
			assert.False(t, asset.IsDeleted())
			assert.Nil(t, asset.Parent())
			assert.NotZero(t, asset.UpdatedAt())
		})
	}
}

func TestNewFromRepository(t *testing.T) {
	id := newUUIDv7(t)
	tenantID := newUUIDv7(t)
	branchID := newUUIDv7(t)
	name := "Test Asset"
	description := "Test Description"
	parent := newUUIDv7(t)
	now := time.Now()
	deletedAt := time.Now()

	asset := NewFromRepository(id, tenantID, branchID, name, description, &parent, now, &deletedAt)

	assert.Equal(t, id, asset.ID())
	assert.Equal(t, tenantID, asset.TenantID())
	assert.Equal(t, branchID, asset.branchName)
	assert.Equal(t, name, asset.Name())
	assert.Equal(t, description, asset.Description())
	assert.Equal(t, &parent, asset.Parent())
	assert.Equal(t, now, asset.UpdatedAt())
	assert.Equal(t, &deletedAt, asset.DeletedAt())
}

func TestAsset_Update(t *testing.T) {
    asset, err := New(newUUIDv7(t), newUUIDv7(t), "Original Name", "Original Description")
    require.NoError(t, err)

    originalUpdatedAt := asset.UpdatedAt()
    time.Sleep(time.Millisecond) // Ensure time difference

    tests := []struct {
        name        string
        newName     string
        parent      string
        description string
        wantErr     bool
    }{
        {
            name:        "valid update",
            newName:     "New Name",
            parent:      newUUIDv7(t),
            description: "New Description",
            wantErr:     false,
        },
        {
            name:        "invalid name - too short",
            newName:     "A",
            parent:      newUUIDv7(t),
            description: "New Description",
            wantErr:     true,
        },
        {
            name:        "invalid parent uuid",
            newName:     "Valid Name",
            parent:      "invalid-uuid",
            description: "Valid Description",
            wantErr:     true,
        },
        {
            name:        "invalid description - too short",
            newName:     "Valid Name",
            parent:      newUUIDv7(t),
            description: "ab",
            wantErr:     true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create a new asset for each test case to avoid state sharing
            testAsset, err := New(newUUIDv7(t), newUUIDv7(t), "Original Name", "Original Description")
            require.NoError(t, err)

            err = testAsset.Update(tt.newName, tt.parent, tt.description)
            if tt.wantErr {
                assert.Error(t, err)
                if err != nil {
                    t.Logf("Got expected error: %v", err)
                }
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.newName, testAsset.Name())
            assert.Equal(t, tt.description, testAsset.Description())
            if tt.parent != "" {
                require.NotNil(t, testAsset.Parent())
                assert.Equal(t, tt.parent, *testAsset.Parent())
            }
            assert.True(t, testAsset.UpdatedAt().After(originalUpdatedAt))
        })
    }
}

func TestAsset_Delete(t *testing.T) {
	asset, err := New(newUUIDv7(t), newUUIDv7(t), "Test Asset", "Test Description")
	require.NoError(t, err)

	assert.False(t, asset.IsDeleted())
	assert.Nil(t, asset.DeletedAt())

	asset.Delete()

	assert.True(t, asset.IsDeleted())
	assert.NotNil(t, asset.DeletedAt())
}

func TestAsset_WithParent(t *testing.T) {
	asset, err := New(newUUIDv7(t), newUUIDv7(t), "Test Asset", "Test Description")
	require.NoError(t, err)

	assert.Nil(t, asset.Parent())

	parentID := newUUIDv7(t)
	asset.WithParent(parentID)

	assert.Equal(t, &parentID, asset.Parent())
}

func TestAsset_Getters(t *testing.T) {
	id := newUUIDv7(t)
	tenantID := newUUIDv7(t)
	branchID := newUUIDv7(t)
	name := "Test Asset"
	description := "Test Description"
	now := time.Now()

	asset := NewFromRepository(id, tenantID, branchID, name, description, nil, now, nil)

	assert.Equal(t, id, asset.ID())
	assert.Equal(t, tenantID, asset.TenantID())
	assert.Equal(t, name, asset.Name())
	assert.Equal(t, description, asset.Description())
	assert.Equal(t, now, asset.UpdatedAt())
	assert.Nil(t, asset.DeletedAt())
	assert.Nil(t, asset.Parent())
}
