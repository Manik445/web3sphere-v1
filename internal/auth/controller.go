package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/internal/middleware"
	"github.com/web3sphere/backend/pkg/response"
)

// Controller handles authentication HTTP requests.
type Controller struct {
	service Service
}

// NewController creates a new Auth controller.
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := ctrl.service.Signup(c.Request.Context(), &req); err != nil {
		c.Error(err)
		return
	}

	response.Created(c, "Registration successful. Please check your email for the verification code.", nil)
}

func (ctrl *Controller) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := ctrl.service.VerifyEmail(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Email verified successfully", res)
}

func (ctrl *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ip := c.GetString("client_ip")
	if ip == "" {
		ip = c.ClientIP()
	}
	userAgent := c.Request.UserAgent()

	res, err := ctrl.service.Login(c.Request.Context(), &req, ip, userAgent)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

func (ctrl *Controller) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ip := c.GetString("client_ip")
	if ip == "" {
		ip = c.ClientIP()
	}
	userAgent := c.Request.UserAgent()

	res, err := ctrl.service.RefreshToken(c.Request.Context(), &req, ip, userAgent)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed", res)
}

func (ctrl *Controller) Logout(c *gin.Context) {
	accessTokenID := c.GetString("token_id")
	
	// Try to get refresh token from request body if provided
	var req RefreshRequest
	refreshToken := ""
	if err := c.ShouldBindJSON(&req); err == nil {
		refreshToken = req.RefreshToken
	}

	if err := ctrl.service.Logout(c.Request.Context(), accessTokenID, refreshToken); err != nil {
		c.Error(err)
		return
	}

	response.OK(c, "Logged out successfully", nil)
}

func (ctrl *Controller) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := ctrl.service.ForgotPassword(c.Request.Context(), &req); err != nil {
		c.Error(err)
		return
	}

	response.OK(c, "If your email is registered, you will receive a password reset code.", nil)
}

func (ctrl *Controller) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := ctrl.service.ResetPassword(c.Request.Context(), &req); err != nil {
		c.Error(err)
		return
	}

	response.OK(c, "Password reset successfully", nil)
}

func (ctrl *Controller) ResendOTP(c *gin.Context) {
	var req ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := ctrl.service.ResendOTP(c.Request.Context(), &req); err != nil {
		c.Error(err)
		return
	}

	response.OK(c, "A new verification code has been sent", nil)
}

func (ctrl *Controller) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	
	// Just an example, in a real scenario you would fetch the latest user info from the DB
	response.OK(c, "User profile", gin.H{"id": userID})
}
