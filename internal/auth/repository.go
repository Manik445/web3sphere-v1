package auth

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines the database operations for the Auth module.
type Repository interface {
	CreateUser(ctx context.Context, user *User, userInfo *UserInfo) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	
	CreateSession(ctx context.Context, session *UserSession) error
	GetSessionByRefreshToken(ctx context.Context, token string) (*UserSession, error)
	RevokeSession(ctx context.Context, id string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	
	TrackDevice(ctx context.Context, device *UserDevice) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new Auth repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User, userInfo *UserInfo) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		userInfo.UserID = user.ID
		if err := tx.Create(userInfo).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Preload("UserInfo").Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Preload("UserInfo").Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *repository) CreateSession(ctx context.Context, session *UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *repository) GetSessionByRefreshToken(ctx context.Context, token string) (*UserSession, error) {
	var session UserSession
	err := r.db.WithContext(ctx).Where("refresh_token = ?", token).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *repository) RevokeSession(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&UserSession{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *repository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&UserSession{}).Where("user_id = ?", userID).Update("revoked", true).Error
}

func (r *repository) TrackDevice(ctx context.Context, device *UserDevice) error {
	// Upsert device
	return r.db.WithContext(ctx).Where(UserDevice{UserID: device.UserID, DeviceIdentifier: device.DeviceIdentifier}).
		Assign(UserDevice{
			DeviceName: device.DeviceName,
			LastIP:     device.LastIP,
			LastLogin:  device.LastLogin,
		}).FirstOrCreate(device).Error
}
