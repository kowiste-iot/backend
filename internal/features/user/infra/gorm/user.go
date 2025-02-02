// repository/user_repository.go
package repository

import (
	"backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"
	baseCmd "backend/shared/base/command"
	"context"

	gormhelper "backend/shared/gorm"
	"backend/shared/http/httputil"
	"backend/shared/pagination"
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         string `gorm:"primaryKey"`
	TenantID   string `gorm:"index"`
	BranchID   string `gorm:"index"`
	AuthUserID string `gorm:"uniqueIndex"`
	Email      string
	FirstName  string
	LastName   string
	Name       string
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.UserRepository {
	db.AutoMigrate(&User{})
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_tenant_email ON users(tenant_id, email)")

	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, input *domain.User) error {
	dbUser := User{
		ID:         input.ID(),
		TenantID:   input.TenantID(),
		BranchID:   input.Branch(),
		AuthUserID: input.AuthID(),
		Email:      input.Email(),
		FirstName:  input.FirstName(),
		LastName:   input.LastName(),
	}
	return r.db.WithContext(ctx).Create(&dbUser).Error
}
func (r *userRepository) Update(ctx context.Context, input *command.UpdateUserInput) error {
	dbUser := User{
		ID:        input.ID,
		TenantID:  input.TenantDomain,
		BranchID:  input.BranchName,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	return r.db.WithContext(ctx).Updates(&dbUser).Error
}

func (r *userRepository) FindByID(ctx context.Context, input *command.UserIDInput) (*domain.User, error) {
	var dbUser User

	err := r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", input.UserID).
		First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return domain.NewFromRepository(
		dbUser.ID,
		dbUser.TenantID,
		dbUser.AuthUserID,
		dbUser.Email,
		dbUser.FirstName,
		dbUser.LastName,
		dbUser.UpdatedAt,
		&dbUser.DeletedAt.Time,
	), nil
}

func (r *userRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.User, error) {
	var dbUsers []User

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}

	var total int64
	err := r.db.Model(&User{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).Count(&total).Error

	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbUsers).Error
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = domain.NewFromRepository(
			dbUser.ID,
			dbUser.TenantID,
			dbUser.AuthUserID,
			dbUser.Email,
			dbUser.FirstName,
			dbUser.LastName,
			dbUser.UpdatedAt,
			&dbUser.DeletedAt.Time,
		)
	}
	return users, nil
}

func (r *userRepository) Remove(ctx context.Context, input *command.UserIDInput) error {

	result := r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", input.UserID).
		Delete(&User{})

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return result.Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser User
	tenantID, ok := httputil.GetTenant(ctx)
	if !ok {
		return nil, errors.New("not tenant id")
	}

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return domain.NewFromRepository(
		dbUser.ID,
		dbUser.TenantID,
		dbUser.AuthUserID,
		dbUser.Email,
		dbUser.FirstName,
		dbUser.LastName,
		dbUser.UpdatedAt,
		&dbUser.DeletedAt.Time,
	), nil
}
