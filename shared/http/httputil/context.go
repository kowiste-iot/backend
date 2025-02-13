package httputil

import (
	"backend/internal/features/tenant/domain"
	"context"
	"errors"
)

// TENANT
type tenantKey struct{}

func SetTenant(ctx context.Context, tenant *domain.Tenant) context.Context {
	return context.WithValue(ctx, tenantKey{}, tenant)
}
func GetTenant(ctx context.Context) (*domain.Tenant, bool) {
	val := ctx.Value(tenantKey{})
	p, ok := val.(*domain.Tenant)
	return p, ok
}

// BRANCH
type branchKey struct{}

func SetBranch(ctx context.Context, branch string) context.Context {
	return context.WithValue(ctx, branchKey{}, branch)
}
func GetBranch(ctx context.Context) (string, bool) {
	val := ctx.Value(branchKey{})
	p, ok := val.(string)
	return p, ok
}

func GetBase(ctx context.Context) (tenant *domain.Tenant, branch string, err error) {
	branch, ok := GetBranch(ctx)
	if !ok {
		err = errors.New("branch not found")
		return
	}
	tenant, ok = GetTenant(ctx)
	if !ok {
		err = errors.New("tenant not found")
		return
	}
	return
}

//TOKEN

type tokenKey struct{}

func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}
func GetToken(ctx context.Context) (string, bool) {
	val := ctx.Value(tokenKey{})
	p, ok := val.(string)
	return p, ok
}

// USERID
type userIDKey struct{}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	val := ctx.Value(userIDKey{})
	p, ok := val.(string)
	return p, ok
}
