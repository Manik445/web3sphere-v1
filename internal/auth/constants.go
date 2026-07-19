package auth

const (
	// User Roles
	RoleUser       = "user"
	RoleFreelancer = "freelancer"
	RoleCompany    = "company"
	RoleModerator  = "moderator"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "super_admin"
	RoleSupport    = "support"

	// User Status
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusSuspended = "suspended"
	StatusBanned    = "banned"
	
	// KYC Status
	KYCStatusUnverified = "unverified"
	KYCStatusPending    = "pending"
	KYCStatusVerified   = "verified"
	KYCStatusRejected   = "rejected"
)
