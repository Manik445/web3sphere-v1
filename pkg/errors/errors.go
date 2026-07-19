package apperrors

import (
	"fmt"
	"net/http"
)

// ErrorCode defines application-specific error codes.
type ErrorCode string

const (
	// Authentication errors
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeInvalidToken     ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired     ErrorCode = "TOKEN_EXPIRED"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeAccountDisabled  ErrorCode = "ACCOUNT_DISABLED"
	ErrCodeEmailNotVerified ErrorCode = "EMAIL_NOT_VERIFIED"

	// Validation errors
	ErrCodeValidation       ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"

	// Resource errors
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	ErrCodeConflict         ErrorCode = "CONFLICT"

	// Permission errors
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeInsufficientRole ErrorCode = "INSUFFICIENT_ROLE"

	// Rate limiting
	ErrCodeRateLimited      ErrorCode = "RATE_LIMITED"

	// Server errors
	ErrCodeInternal         ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceDown      ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeDatabaseError    ErrorCode = "DATABASE_ERROR"
	ErrCodeCacheError       ErrorCode = "CACHE_ERROR"

	// OTP errors
	ErrCodeInvalidOTP       ErrorCode = "INVALID_OTP"
	ErrCodeOTPExpired       ErrorCode = "OTP_EXPIRED"
	ErrCodeMaxOTPAttempts   ErrorCode = "MAX_OTP_ATTEMPTS"

	// Mail errors
	ErrCodeMailFailed       ErrorCode = "MAIL_FAILED"
)

// AppError represents a structured application error.
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
	Err        error     `json:"-"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError.
func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Wrap wraps an existing error into an AppError.
func Wrap(code ErrorCode, message string, httpStatus int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

// --- Convenience constructors ---

func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

func InvalidCredentials() *AppError {
	return New(ErrCodeInvalidCredentials, "Invalid email or password", http.StatusUnauthorized)
}

func InvalidToken() *AppError {
	return New(ErrCodeInvalidToken, "Invalid or malformed token", http.StatusUnauthorized)
}

func TokenExpired() *AppError {
	return New(ErrCodeTokenExpired, "Token has expired", http.StatusUnauthorized)
}

func EmailNotVerified() *AppError {
	return New(ErrCodeEmailNotVerified, "Email address not verified", http.StatusForbidden)
}

func AccountDisabled() *AppError {
	return New(ErrCodeAccountDisabled, "Account is disabled", http.StatusForbidden)
}

func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message, http.StatusForbidden)
}

func NotFound(resource string) *AppError {
	return New(ErrCodeNotFound, fmt.Sprintf("%s not found", resource), http.StatusNotFound)
}

func AlreadyExists(resource string) *AppError {
	return New(ErrCodeAlreadyExists, fmt.Sprintf("%s already exists", resource), http.StatusConflict)
}

func ValidationFailed(message string) *AppError {
	return New(ErrCodeValidation, message, http.StatusUnprocessableEntity)
}

func InternalError(err error) *AppError {
	return Wrap(ErrCodeInternal, "An internal error occurred", http.StatusInternalServerError, err)
}

func DatabaseError(err error) *AppError {
	return Wrap(ErrCodeDatabaseError, "Database operation failed", http.StatusInternalServerError, err)
}

func CacheError(err error) *AppError {
	return Wrap(ErrCodeCacheError, "Cache operation failed", http.StatusInternalServerError, err)
}

func RateLimited() *AppError {
	return New(ErrCodeRateLimited, "Too many requests, please try again later", http.StatusTooManyRequests)
}

func InvalidOTP() *AppError {
	return New(ErrCodeInvalidOTP, "Invalid OTP", http.StatusBadRequest)
}

func OTPExpired() *AppError {
	return New(ErrCodeOTPExpired, "OTP has expired", http.StatusBadRequest)
}

func MaxOTPAttempts() *AppError {
	return New(ErrCodeMaxOTPAttempts, "Maximum OTP attempts exceeded", http.StatusTooManyRequests)
}

func MailFailed(err error) *AppError {
	return Wrap(ErrCodeMailFailed, "Failed to send email", http.StatusInternalServerError, err)
}
