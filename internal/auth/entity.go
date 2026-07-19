package auth

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table.
type User struct {
	ID            string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string         `gorm:"type:varchar(255);not null"`
	Status        string         `gorm:"type:varchar(50);not null;default:'active'"`
	Role          string         `gorm:"type:varchar(50);not null;default:'user'"`
	EmailVerified bool           `gorm:"default:false"`
	PhoneVerified bool           `gorm:"default:false"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Associations
	UserInfo *UserInfo `gorm:"foreignKey:UserID"`
	Sessions []UserSession `gorm:"foreignKey:UserID"`
	Devices  []UserDevice  `gorm:"foreignKey:UserID"`
}

// UserInfo represents the user_info table.
type UserInfo struct {
	UserID        string    `gorm:"type:uuid;primaryKey"`
	FirstName     string    `gorm:"type:varchar(100)"`
	LastName      string    `gorm:"type:varchar(100)"`
	Avatar        string    `gorm:"type:text"`
	CountryID     *int      `gorm:"type:int"`
	Timezone      string    `gorm:"type:varchar(100)"`
	Language      string    `gorm:"type:varchar(20)"`
	Bio           string    `gorm:"type:text"`
	Website       string    `gorm:"type:text"`
	Github        string    `gorm:"type:text"`
	Linkedin      string    `gorm:"type:text"`
	Twitter       string    `gorm:"type:text"`
	WalletAddress string    `gorm:"type:varchar(100);uniqueIndex"`
	KYCStatus     string    `gorm:"type:varchar(50);default:'unverified'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// UserSession represents the user_sessions table.
type UserSession struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        string    `gorm:"type:uuid;not null;index"`
	RefreshToken  string    `gorm:"type:varchar(512);uniqueIndex;not null"`
	AccessTokenID string    `gorm:"type:varchar(255);not null"`
	IP            string    `gorm:"type:varchar(50)"`
	Browser       string    `gorm:"type:varchar(255)"`
	OS            string    `gorm:"type:varchar(100)"`
	Country       string    `gorm:"type:varchar(100)"`
	Device        string    `gorm:"type:varchar(255)"`
	ExpiresAt     time.Time `gorm:"not null"`
	LastActivity  time.Time `gorm:"autoCreateTime"`
	Revoked       bool      `gorm:"default:false"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// UserDevice represents the user_devices table.
type UserDevice struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID           string    `gorm:"type:uuid;not null"`
	DeviceIdentifier string    `gorm:"type:varchar(255);not null"`
	DeviceName       string    `gorm:"type:varchar(255)"`
	Trusted          bool      `gorm:"default:false"`
	LastIP           string    `gorm:"type:varchar(50)"`
	LastLogin        time.Time
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}

// TempData represents the temp_data table.
type TempData struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Identifier string    `gorm:"type:varchar(255);not null;index:idx_temp_data_identifier_type"`
	Type       string    `gorm:"type:varchar(50);not null;index:idx_temp_data_identifier_type"`
	Data       string    `gorm:"type:jsonb;not null"`
	ExpiresAt  time.Time `gorm:"not null;index"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
