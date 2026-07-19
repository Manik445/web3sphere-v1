package users

import (
	"context"

	"github.com/web3sphere/backend/pkg/errors"
	"github.com/web3sphere/backend/pkg/logger"
)

// Service defines the business logic for Users.
type Service interface {
	GetProfile(ctx context.Context, userID string) (*UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserProfileResponse, error)
}

type service struct {
	repo Repository
	log  *logger.Logger
}

// NewService creates a new Users service.
func NewService(repo Repository, log *logger.Logger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) GetProfile(ctx context.Context, userID string) (*UserProfileResponse, error) {
	user, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	if user == nil {
		return nil, apperrors.NotFound("User")
	}

	return s.mapToResponse(user), nil
}

func (s *service) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserProfileResponse, error) {
	user, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	if user == nil {
		return nil, apperrors.NotFound("User")
	}

	if user.UserInfo == nil {
		user.UserInfo = &UserInfo{UserID: user.ID}
	}

	// Update fields if provided
	if req.FirstName != "" { user.UserInfo.FirstName = req.FirstName }
	if req.LastName != "" { user.UserInfo.LastName = req.LastName }
	if req.Bio != "" { user.UserInfo.Bio = req.Bio }
	if req.Timezone != "" { user.UserInfo.Timezone = req.Timezone }
	if req.Language != "" { user.UserInfo.Language = req.Language }
	if req.Website != "" { user.UserInfo.Website = req.Website }
	if req.Github != "" { user.UserInfo.Github = req.Github }
	if req.Linkedin != "" { user.UserInfo.Linkedin = req.Linkedin }
	if req.Twitter != "" { user.UserInfo.Twitter = req.Twitter }
	if req.CountryID != nil { user.UserInfo.CountryID = req.CountryID }

	if err := s.repo.UpdateUserInfo(ctx, user.UserInfo); err != nil {
		return nil, apperrors.DatabaseError(err)
	}

	return s.mapToResponse(user), nil
}

func (s *service) mapToResponse(user *User) *UserProfileResponse {
	resp := &UserProfileResponse{
		ID:            user.ID,
		Email:         user.Email,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		JoinedAt:      user.CreatedAt,
	}

	if user.UserInfo != nil {
		resp.FirstName = user.UserInfo.FirstName
		resp.LastName = user.UserInfo.LastName
		resp.Avatar = user.UserInfo.Avatar
		resp.Bio = user.UserInfo.Bio
		resp.Website = user.UserInfo.Website
		resp.Github = user.UserInfo.Github
		resp.Linkedin = user.UserInfo.Linkedin
		resp.Twitter = user.UserInfo.Twitter
		resp.WalletAddress = user.UserInfo.WalletAddress
		resp.KYCStatus = user.UserInfo.KYCStatus
	}

	return resp
}
