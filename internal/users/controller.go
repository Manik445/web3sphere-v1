package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/internal/middleware"
	"github.com/web3sphere/backend/pkg/response"
)

// Controller handles Users HTTP requests.
type Controller struct {
	service Service
}

// NewController creates a new Users controller.
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// GetProfile handles getting a user's profile.
func (ctrl *Controller) GetProfile(c *gin.Context) {
	// If getting own profile, ID is from context
	userID := c.Param("id")
	if userID == "me" || userID == "" {
		userID = middleware.GetUserID(c)
	}

	res, err := ctrl.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", res)
}

// UpdateProfile handles updating a user's profile.
func (ctrl *Controller) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := ctrl.service.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", res)
}
