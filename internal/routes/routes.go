package routes

import (
	"golang-auth-api-mysql/config"
	"golang-auth-api-mysql/internal/handler"
	"golang-auth-api-mysql/internal/middleware"
	"golang-auth-api-mysql/internal/repository"
	"golang-auth-api-mysql/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router gin.IRouter, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	auditService := service.NewAuditService(auditRepo)

	// Initialize handler
	h := handler.NewHandler(authService, userService, productService, auditService)

	// Apply global middleware
	router.Use(middleware.RateLimit())

	//index route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to API lur!"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
			auth.POST("/logout", h.Logout)
			auth.GET("/verify-email", h.VerifyEmail)
			auth.POST("/resend-verification", h.ResendVerification)
			auth.POST("/forgot-password", h.ForgotPassword)
			auth.POST("/reset-password", h.ResetPassword)
		}

		// Public product routes
		v1.GET("/products", h.GetAllProducts)
		v1.GET("/products/:id", h.GetProduct)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg))
		protected.Use(middleware.AuditMiddleware(auditService))
		{
			user := protected.Group("/users")
			{
				user.GET("/me", h.GetProfile)
				user.PUT("/me", h.UpdateProfile)
				user.PUT("/me/password", h.ChangePassword)
			}

			admin := protected.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("/users", h.GetAllUsers)
				admin.DELETE("/users/:id", h.DeleteUser)

				products := admin.Group("/products")
				{
					products.POST("", h.CreateProduct)
					products.PUT("/:id", h.UpdateProduct)
					products.DELETE("/:id", h.DeleteProduct)
					products.POST("/:id/image", handler.UploadProductImage(productService, cfg))
				}

				admin.GET("/audit-logs", h.GetAuditLogs)
			}
		}
	}

	// Swagger
	// router.GET("/swagger/*any", func(c *gin.Context) {
	// 	c.File("./docs/swagger.yaml")
	// })

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
