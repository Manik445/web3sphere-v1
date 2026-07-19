package users

import "time"

// UpdateProfileRequest represents the payload for updating a user profile.
type UpdateProfileRequest struct {
	FirstName     string `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName      string `json:"last_name" binding:"omitempty,min=2,max=100"`
	Bio           string `json:"bio" binding:"omitempty,max=500"`
	Timezone      string `json:"timezone" binding:"omitempty"`
	Language      string `json:"language" binding:"omitempty,max=20"`
	Website       string `json:"website" binding:"omitempty,url"`
	Github        string `json:"github" binding:"omitempty,url"`
	Linkedin      string `json:"linkedin" binding:"omitempty,url"`
	Twitter       string `json:"twitter" binding:"omitempty,url"`
	CountryID     *int   `json:"country_id" binding:"omitempty,min=1"`
}

// UserProfileResponse represents a complete public user profile.
type UserProfileResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Avatar        string    `json:"avatar"`
	Bio           string    `json:"bio"`
	Website       string    `json:"website"`
	Github        string    `json:"github"`
	Linkedin      string    `json:"linkedin"`
	Twitter       string    `json:"twitter"`
	WalletAddress string    `json:"wallet_address"`
	KYCStatus     string    `json:"kyc_status"`
	JoinedAt      time.Time `json:"joined_at"`
}
