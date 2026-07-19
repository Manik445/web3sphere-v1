package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/pkg/response"
)

// Role constants match the user roles defined in the database.
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleModerator  = "moderator"
	RoleCompany    = "company"
	RoleFreelancer = "freelancer"
	RoleUser       = "user"
	RoleSupport    = "support"
)

// roleHierarchy defines role precedence (higher = more permissions).
var roleHierarchy = map[string]int{
	RoleSuperAdmin: 100,
	RoleAdmin:      90,
	RoleModerator:  70,
	RoleSupport:    60,
	RoleCompany:    50,
	RoleFreelancer: 40,
	RoleUser:       10,
}

// RequireRole creates a middleware that checks if the user has any of the allowed roles.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		if userRole == "" {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}

// RequireMinRole checks if the user's role meets a minimum hierarchy level.
func RequireMinRole(minRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		if userRole == "" {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		userLevel, ok := roleHierarchy[userRole]
		if !ok {
			response.Forbidden(c, "Unknown role")
			c.Abort()
			return
		}

		minLevel, ok := roleHierarchy[minRole]
		if !ok {
			response.Forbidden(c, "Unknown minimum role")
			c.Abort()
			return
		}

		if userLevel < minLevel {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnly restricts access to admin and super admin roles.
func AdminOnly() gin.HandlerFunc {
	return RequireRole(RoleSuperAdmin, RoleAdmin)
}

// ModeratorOrAbove restricts access to moderators, admins, and super admins.
func ModeratorOrAbove() gin.HandlerFunc {
	return RequireMinRole(RoleModerator)
}
