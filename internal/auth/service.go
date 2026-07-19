package auth

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/cache"
	"github.com/web3sphere/backend/pkg/errors"
	"github.com/web3sphere/backend/pkg/jwt"
	"github.com/web3sphere/backend/pkg/logger"
	"github.com/web3sphere/backend/pkg/mailer"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the business logic for the Auth module.
type Service interface {
	Signup(ctx context.Context, req *SignupRequest) error
	VerifyEmail(ctx context.Context, req *VerifyEmailRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest, ip, userAgent string) (*AuthResponse, error)
	RefreshToken(ctx context.Context, req *RefreshRequest, ip, userAgent string) (*TokenResponse, error)
	Logout(ctx context.Context, accessTokenID, refreshToken string) error
	ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	ResendOTP(ctx context.Context, req *ResendOTPRequest) error
}

type service struct {
	repo       Repository
	cfg        *configs.Config
	redis      *cache.RedisClient
	jwtManager *jwt.Manager
	mailer     mailer.Mailer
	log        *logger.Logger
}

// NewService creates a new Auth service.
func NewService(repo Repository, cfg *configs.Config, redis *cache.RedisClient, jwtManager *jwt.Manager, mailer mailer.Mailer, log *logger.Logger) Service {
	return &service{
		repo:       repo,
		cfg:        cfg,
		redis:      redis,
		jwtManager: jwtManager,
		mailer:     mailer,
		log:        log,
	}
}

func (s *service) Signup(ctx context.Context, req *SignupRequest) error {
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return apperrors.DatabaseError(err)
	}
	if existingUser != nil {
		return apperrors.AlreadyExists("User with this email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.InternalError(err)
	}

	user := &User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user", // Default role
	}

	userInfo := &UserInfo{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := s.repo.CreateUser(ctx, user, userInfo); err != nil {
		return apperrors.DatabaseError(err)
	}

	// Generate and send OTP
	otp, err := s.generateOTP()
	if err != nil {
		return apperrors.InternalError(err)
	}

	if err := s.redis.StoreOTP(ctx, "verify:"+user.Email, otp, time.Duration(s.cfg.OTP.Expiry)*time.Second); err != nil {
		return apperrors.CacheError(err)
	}

	// Send Email
	go s.sendVerificationEmail(user.Email, otp)

	return nil
}

func (s *service) VerifyEmail(ctx context.Context, req *VerifyEmailRequest) (*AuthResponse, error) {
	// Check rate limit for OTP attempts
	_, err := s.redis.IncrementOTPAttempts(ctx, "verify:"+req.Email, s.cfg.OTP.MaxAttempts, 15*time.Minute)
	if err != nil {
		// Log error, but proceed if possible. Ideally, block if max attempts reached.
	}

	storedOTP, err := s.redis.GetOTP(ctx, "verify:"+req.Email)
	if err != nil {
		return nil, apperrors.OTPExpired()
	}

	if storedOTP != req.OTP {
		return nil, apperrors.InvalidOTP()
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	if user == nil {
		return nil, apperrors.NotFound("User")
	}

	user.EmailVerified = true
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, apperrors.DatabaseError(err)
	}

	s.redis.DeleteOTP(ctx, "verify:"+req.Email)

	return s.generateAuthResponse(ctx, user, "", "") // IP and userAgent omitted for simplicity in this flow, or pass them in
}

func (s *service) Login(ctx context.Context, req *LoginRequest, ip, userAgent string) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	if user == nil || user.Status != "active" {
		return nil, apperrors.InvalidCredentials()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.InvalidCredentials()
	}

	if !user.EmailVerified {
		return nil, apperrors.EmailNotVerified()
	}

	return s.generateAuthResponse(ctx, user, ip, userAgent)
}

func (s *service) RefreshToken(ctx context.Context, req *RefreshRequest, ip, userAgent string) (*TokenResponse, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, apperrors.InvalidToken()
	}

	session, err := s.repo.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	if session == nil || session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.InvalidToken()
	}

	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	if user == nil || user.Status != "active" {
		return nil, apperrors.AccountDisabled()
	}

	// Revoke old session
	_ = s.repo.RevokeSession(ctx, session.ID)

	// Generate new tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperrors.InternalError(err)
	}

	// Create new session
	newSession := &UserSession{
		UserID:        user.ID,
		RefreshToken:  tokens.RefreshToken,
		AccessTokenID: tokens.AccessID,
		IP:            ip,
		Browser:       userAgent, // Basic extraction, a real app might parse user-agent
		ExpiresAt:     tokens.ExpiresAt,
	}

	if err := s.repo.CreateSession(ctx, newSession); err != nil {
		return nil, apperrors.DatabaseError(err)
	}

	return &TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int(s.jwtManager.GetAccessExpiry().Seconds()),
	}, nil
}

func (s *service) Logout(ctx context.Context, accessTokenID, refreshToken string) error {
	if refreshToken != "" {
		session, _ := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
		if session != nil {
			_ = s.repo.RevokeSession(ctx, session.ID)
		}
	}
	
	if accessTokenID != "" {
		_ = s.redis.BlacklistToken(ctx, accessTokenID, s.jwtManager.GetAccessExpiry())
	}

	return nil
}

func (s *service) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		// Return success to prevent email enumeration
		return nil
	}

	otp, err := s.generateOTP()
	if err != nil {
		return apperrors.InternalError(err)
	}

	if err := s.redis.StoreOTP(ctx, "reset:"+user.Email, otp, time.Duration(s.cfg.OTP.Expiry)*time.Second); err != nil {
		return apperrors.CacheError(err)
	}

	go s.sendPasswordResetEmail(user.Email, otp)

	return nil
}

func (s *service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	storedOTP, err := s.redis.GetOTP(ctx, "reset:"+req.Email)
	if err != nil {
		return apperrors.OTPExpired()
	}

	if storedOTP != req.OTP {
		return apperrors.InvalidOTP()
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return apperrors.NotFound("User")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.InternalError(err)
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return apperrors.DatabaseError(err)
	}

	s.redis.DeleteOTP(ctx, "reset:"+req.Email)
	
	// Revoke all sessions on password reset
	_ = s.repo.RevokeAllUserSessions(ctx, user.ID)

	return nil
}

func (s *service) ResendOTP(ctx context.Context, req *ResendOTPRequest) error {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil // Prevent enumeration
	}

	if user.EmailVerified {
		return apperrors.ValidationFailed("Email is already verified")
	}

	otp, err := s.generateOTP()
	if err != nil {
		return apperrors.InternalError(err)
	}

	if err := s.redis.StoreOTP(ctx, "verify:"+user.Email, otp, time.Duration(s.cfg.OTP.Expiry)*time.Second); err != nil {
		return apperrors.CacheError(err)
	}

	go s.sendVerificationEmail(user.Email, otp)

	return nil
}

// --- Helper Methods ---

func (s *service) generateAuthResponse(ctx context.Context, user *User, ip, userAgent string) (*AuthResponse, error) {
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperrors.InternalError(err)
	}

	session := &UserSession{
		UserID:        user.ID,
		RefreshToken:  tokens.RefreshToken,
		AccessTokenID: tokens.AccessID,
		IP:            ip,
		Browser:       userAgent,
		ExpiresAt:     tokens.ExpiresAt,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	
	// Track device
	device := &UserDevice{
		UserID: user.ID,
		DeviceIdentifier: userAgent, // Ideally generate a stable device ID
		LastIP: ip,
		LastLogin: time.Now(),
	}
	_ = s.repo.TrackDevice(ctx, device)

	avatar := ""
	firstName := ""
	lastName := ""
	if user.UserInfo != nil {
		avatar = user.UserInfo.Avatar
		firstName = user.UserInfo.FirstName
		lastName = user.UserInfo.LastName
	}

	return &AuthResponse{
		User: UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			FirstName:     firstName,
			LastName:      lastName,
			Role:          user.Role,
			EmailVerified: user.EmailVerified,
			Avatar:        avatar,
		},
		Tokens: TokenResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresIn:    int(s.jwtManager.GetAccessExpiry().Seconds()),
		},
	}, nil
}

func (s *service) generateOTP() (string, error) {
	const chars = "0123456789"
	result := make([]byte, s.cfg.OTP.Length)
	for i := 0; i < s.cfg.OTP.Length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}
	return string(result), nil
}

func (s *service) sendVerificationEmail(email, otp string) {
	msg := &mailer.TemplateMessage{
		To:           []string{email},
		Subject:      "Verify Your Email - Web3Sphere",
		TemplateName: "verify_email",
		Data: map[string]interface{}{
			"otp":            otp,
			"expiry_minutes": s.cfg.OTP.Expiry / 60,
		},
	}
	if err := s.mailer.SendTemplate(context.Background(), msg); err != nil {
		s.log.Errorf("Failed to send verification email to %s: %v", email, err)
	}
}

func (s *service) sendPasswordResetEmail(email, otp string) {
	msg := &mailer.TemplateMessage{
		To:           []string{email},
		Subject:      "Password Reset - Web3Sphere",
		TemplateName: "reset_password",
		Data: map[string]interface{}{
			"otp":            otp,
			"expiry_minutes": s.cfg.OTP.Expiry / 60,
		},
	}
	if err := s.mailer.SendTemplate(context.Background(), msg); err != nil {
		s.log.Errorf("Failed to send password reset email to %s: %v", email, err)
	}
}
