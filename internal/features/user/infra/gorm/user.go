// repository/user_repository.go
package repository

import (
	"backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"
	baseCmd "backend/shared/base/command"
	"context"
	"time"

	gormhelper "backend/shared/gorm"
	"backend/shared/http/httputil"
	"backend/shared/pagination"
	"errors"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, input *domain.User) error {
	dbUser := userMapper(input)

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create user
	if err := tx.Create(dbUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add user roles if any
	if len(input.Roles()) > 0 {
		if err := r.saveUserRoles(tx, input.ID(), input.TenantID(), input.Roles()); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *userRepository) Update(ctx context.Context, input *domain.User) error {
	dbUser := userMapper(input)

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update user
	if err := tx.Updates(dbUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update roles if any
	if len(input.Roles()) > 0 {
		// Clear existing roles
		if err := tx.Where("user_id = ?", input.ID()).Delete(&UserRoleDB{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Add updated roles
		if err := r.saveUserRoles(tx, input.ID(), input.TenantID(), input.Roles()); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *userRepository) FindByID(ctx context.Context, input *command.UserIDInput) (*domain.User, error) {
	var dbUser UserDB

	err := r.db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()+" AND id = ?", input.UserID).
		First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	tenant, _ := httputil.GetTenant(ctx)
	// Get user roles
	roles, err := r.getUserRoles(ctx, dbUser.ID)
	if err != nil {
		return nil, err
	}

	return domain.NewFromRepositoryWithRoles(
		dbUser.ID,
		tenant.Domain(),
		dbUser.AuthUserID,
		dbUser.Email,
		dbUser.FirstName,
		dbUser.LastName,
		dbUser.UpdatedAt,
		&dbUser.DeletedAt.Time,
		roles,
	), nil
}

func (r *userRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.User, error) {
	var dbUsers []UserDB

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}

	var total int64
	err := r.db.Model(&UserDB{}).Where(gormhelper.DeleteFilter()).Count(&total).Error

	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbUsers).Error
	if err != nil {
		return nil, err
	}
	tenant, _ := httputil.GetTenant(ctx)
	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		// Get user roles
		roles, err := r.getUserRoles(ctx, dbUser.ID)
		if err != nil {
			return nil, err
		}

		users[i] = domain.NewFromRepositoryWithRoles(
			dbUser.ID,
			tenant.Domain(),
			dbUser.AuthUserID,
			dbUser.Email,
			dbUser.FirstName,
			dbUser.LastName,
			dbUser.UpdatedAt,
			&dbUser.DeletedAt.Time,
			roles,
		)
	}
	return users, nil
}

func (r *userRepository) Remove(ctx context.Context, input *command.UserIDInput) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete user roles
	if err := tx.Where("user_id = ?", input.UserID).Delete(&UserRoleDB{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete user
	result := tx.Where(gormhelper.DeleteFilter()+" AND id = ?", input.UserID).
		Delete(&UserDB{})

	if result.RowsAffected == 0 {
		tx.Rollback()
		return domain.ErrUserNotFound
	}

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser UserDB
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

	// Get user roles
	roles, err := r.getUserRoles(ctx, dbUser.ID)
	if err != nil {
		return nil, err
	}
	tenant, _ := httputil.GetTenant(ctx)
	return domain.NewFromRepositoryWithRoles(
		dbUser.ID,
		tenant.Domain(),
		dbUser.AuthUserID,
		dbUser.Email,
		dbUser.FirstName,
		dbUser.LastName,
		dbUser.UpdatedAt,
		&dbUser.DeletedAt.Time,
		roles,
	), nil
}

// Helper method to save user roles
func (r *userRepository) saveUserRoles(tx *gorm.DB, userID string, tenantID string, roles []string) error {
	for _, roleName := range roles {
		// Get role ID or create if not exists
		var role RoleDB
		err := tx.Where("tenant_id = ? AND name = ?", tenantID, roleName).FirstOrCreate(&role, RoleDB{
			ID:          roleName + "_" + tenantID, // Unique role ID per tenant
			Name:        roleName,
			Description: roleName + " role",
			UpdatedAt:   time.Now(),
		}).Error
		if err != nil {
			return err
		}

		// Create user_role association
		err = tx.Create(&UserRoleDB{
			UserID:    userID,
			RoleID:    role.ID,
			UpdatedAt: time.Now(),
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper method to get user roles
func (r *userRepository) getUserRoles(ctx context.Context, userID string) ([]string, error) {
	var userRoles []UserRoleDB
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []string{}, nil
	}

	// Extract role IDs
	roleIDs := make([]string, len(userRoles))
	for i, userRole := range userRoles {
		roleIDs[i] = userRole.RoleID
	}

	// Get role names
	var roles []RoleDB
	err = r.db.WithContext(ctx).
		Where("id IN ?", roleIDs).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// Extract role names
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return roleNames, nil
}
