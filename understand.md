To add a new route in this modular Clean Architecture, you follow a strict top-to-bottom flow. Since we've already implemented the Auth module, we can use it as the perfect example of how this is done!

Here is the step-by-step process of how routes are added and how the Auth module was implemented:

1. Define the Data (DTOs)
First, define what data the route will accept and return. We use Gin's binding tags for automatic validation.

Where: internal/auth/dto.go
Example:
go
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
2. Implement the Business Logic (Service)
Next, define the interface method for your feature and implement the core logic. The service handles database calls, password hashing, JWT generation, etc., keeping the controller clean.

Where: internal/auth/service.go
Example:
go
type Service interface {
	Login(ctx context.Context, req *LoginRequest, ip, userAgent string) (*AuthResponse, error)
    // ...
}
func (s *service) Login(ctx context.Context, req *LoginRequest, ip, userAgent string) (*AuthResponse, error) {
    // 1. Fetch user from DB
    // 2. Compare passwords
    // 3. Generate JWT tokens
    // 4. Return AuthResponse
}
3. Create the HTTP Handler (Controller)
The controller's only job is to parse the incoming HTTP request, pass it to the service, and return the HTTP response. It does not contain business logic.

Where: internal/auth/controller.go
Example:
go
func (ctrl *Controller) Login(c *gin.Context) {
	var req LoginRequest
	// 1. Bind and Validate JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	// 2. Call Service
	res, err := ctrl.service.Login(c.Request.Context(), &req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.Error(err) // Centralized error handling will catch this
		return
	}
	// 3. Return Standardized Response
	response.Success(c, http.StatusOK, "Login successful", res)
}
4. Register the Route
Group your routes and attach any necessary middlewares (like Rate Limiting or JWT Auth).

Where: internal/auth/routes.go
Example:
go
func RegisterRoutes(router *gin.RouterGroup, ctrl *Controller, jwtManager *jwt.Manager, redis *cache.RedisClient) {
	authRoutes := router.Group("/auth")
	{
		// Apply Strict Rate Limiter to login
		strictLimit := middleware.StrictRateLimiter(redis)
		authRoutes.POST("/login", strictLimit, ctrl.Login)
        
		// Protected routes require the JWT Auth middleware
		protected := authRoutes.Group("")
		protected.Use(middleware.Auth(jwtManager, redis))
		{
			protected.POST("/logout", ctrl.Logout)
			protected.GET("/me", ctrl.Me)
		}
	}
}
5. Wire it up in Main
Finally, if this is a brand new module (like if you were building an escrow module), you wire the layers together in main.go and attach its routes to the global /api/v1 router.

Where: cmd/server/main.go
Example:
go
// Base API Group
	api := router.Group("/api/v1")
	// Auth Module Initialization
	authRepo := auth.NewRepository(db)
	authSvc := auth.NewService(authRepo, cfg, redisClient, jwtManager, mailSvc, log)
	authCtrl := auth.NewController(authSvc)
	
    // Register the routes to /api/v1/auth
	auth.RegisterRoutes(api, authCtrl, jwtManager, redisClient)
The auth routes (/api/v1/auth/signup, /api/v1/auth/login, etc.) are already fully implemented and ready to be used by your frontend!

If you want to create a brand new module (e.g., projects or escrow), you would just duplicate this exact folder structure (entity, dto, repo, service, controller, routes) inside a new folder in internal/.



I created a new file inside each module called module.go which exposes a single Setup() function. This centralizes the dependency injection and route registration inside the module itself.

For example, in internal/auth/module.go, I added:

go
package auth
// Setup wires up the Auth module dependencies and registers its routes.
func Setup(
	router *gin.RouterGroup,
	db *gorm.DB,
	cfg *configs.Config,
	redisClient *cache.RedisClient,
	jwtManager *jwt.Manager,
	mailSvc mailer.Mailer,
	log *logger.Logger,
) {
	repo := NewRepository(db)
	svc := NewService(repo, cfg, redisClient, jwtManager, mailSvc, log)
	ctrl := NewController(svc)
	
	RegisterRoutes(router, ctrl, jwtManager, redisClient)
}
And similarly for the users module in internal/users/module.go.

The Result
Now, your cmd/server/main.go remains incredibly clean and scalable. Adding a new module is just a single line of code!

go
// 12. Initialize Modules
	auth.Setup(api, db, cfg, redisClient, jwtManager, mailSvc, log)
	users.Setup(api, db, jwtManager, redisClient, log)
    // projects.Setup(...)
    // escrow.Setup(...)
This keeps all the wiring completely contained within the specific domain folder, which is exactly how large enterprise Go codebases keep their entry points tidy.