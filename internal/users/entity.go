package users

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
