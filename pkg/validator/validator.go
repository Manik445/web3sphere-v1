package validator

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// walletAddressRegex matches Ethereum-style wallet addresses.
var walletAddressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

// Setup registers custom validators with Gin's validator.
func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("wallet_address", validateWalletAddress)
		v.RegisterValidation("strong_password", validateStrongPassword)
		v.RegisterValidation("username", validateUsername)
		v.RegisterValidation("no_spaces", validateNoSpaces)
	}
}

// validateWalletAddress checks for valid Ethereum wallet addresses.
func validateWalletAddress(fl validator.FieldLevel) bool {
	return walletAddressRegex.MatchString(fl.Field().String())
}

// validateStrongPassword requires min 8 chars, uppercase, lowercase, number, special char.
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateUsername allows alphanumeric, underscores, hyphens; 3-30 chars.
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return matched
}

// validateNoSpaces disallows whitespace.
func validateNoSpaces(fl validator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), " ")
}

// FormatValidationErrors converts validator.ValidationErrors to a map of field errors.
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = "Invalid email address"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + e.Param() + " characters"
			case "wallet_address":
				errors[field] = "Invalid wallet address format"
			case "strong_password":
				errors[field] = "Password must contain at least 8 characters, uppercase, lowercase, number, and special character"
			case "username":
				errors[field] = "Username must be 3-30 characters, alphanumeric with underscores or hyphens"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}
	return errors
}
