package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/web3sphere/backend/configs"
)

// TokenType distinguishes between access and refresh tokens.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims holds the custom JWT claims.
type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AccessID     string    `json:"access_id"`
	RefreshID    string    `json:"refresh_id"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Manager handles JWT operations.
type Manager struct {
	cfg *configs.JWTConfig
}

// NewManager creates a new JWT manager.
func NewManager(cfg *configs.JWTConfig) *Manager {
	return &Manager{cfg: cfg}
}

// GenerateTokenPair creates a new access/refresh token pair.
func (m *Manager) GenerateTokenPair(userID, email, role string) (*TokenPair, error) {
	accessID := uuid.New().String()
	refreshID := uuid.New().String()
	now := time.Now()

	// Access token
	accessClaims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: string(AccessToken),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessID,
			Issuer:    m.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.cfg.AccessExpiry)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(m.cfg.AccessSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token
	refreshClaims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: string(RefreshToken),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshID,
			Issuer:    m.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.cfg.RefreshExpiry)),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(m.cfg.RefreshSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessID:     accessID,
		RefreshID:    refreshID,
		ExpiresAt:    now.Add(m.cfg.AccessExpiry),
	}, nil
}

// ValidateAccessToken parses and validates an access token.
func (m *Manager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	return m.parseToken(tokenStr, m.cfg.AccessSecret, AccessToken)
}

// ValidateRefreshToken parses and validates a refresh token.
func (m *Manager) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return m.parseToken(tokenStr, m.cfg.RefreshSecret, RefreshToken)
}

func (m *Manager) parseToken(tokenStr, secret string, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.TokenType != string(expectedType) {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.TokenType)
	}

	return claims, nil
}

// GetAccessExpiry returns the access token expiry duration.
func (m *Manager) GetAccessExpiry() time.Duration {
	return m.cfg.AccessExpiry
}

// GetRefreshExpiry returns the refresh token expiry duration.
func (m *Manager) GetRefreshExpiry() time.Duration {
	return m.cfg.RefreshExpiry
}
