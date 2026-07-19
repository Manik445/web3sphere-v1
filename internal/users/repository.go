package users

import (
	"context"

	"github.com/web3sphere/backend/internal/common"
	"gorm.io/gorm"
)

// Repository defines database operations for Users.
type Repository interface {
	common.BaseRepository[User]
	
	// Specific user queries
	GetProfile(ctx context.Context, userID string) (*User, error)
	UpdateUserInfo(ctx context.Context, userInfo *UserInfo) error
}

type repository struct {
	*common.GormRepository[User]
	db *gorm.DB
}

// NewRepository creates a new Users repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		GormRepository: common.NewGormRepository[User](db),
		db:             db,
	}
}

func (r *repository) GetProfile(ctx context.Context, userID string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Preload("UserInfo").Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUserInfo(ctx context.Context, userInfo *UserInfo) error {
	return r.db.WithContext(ctx).Save(userInfo).Error
}
