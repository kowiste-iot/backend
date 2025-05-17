package httputil

import (
	"context"

	"gorm.io/gorm"
)

// Define a key for storing the transaction in context
type txKey struct{}

// Function to store a transaction in context
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// Function to retrieve a transaction from context
func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
}

// Helper to get DB connection (transaction or regular DB)
func GetDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return db
}
