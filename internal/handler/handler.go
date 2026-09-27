package handler

import (
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/service"
	"golang-auth-api-mysql/pkg/jwt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService    service.AuthService
	userService    service.UserService
	productService service.ProductService
	auditService   service.AuditService
}

func NewHandler(
	authService service.AuthService,
	userService service.UserService,
	productService service.ProductService,
	auditService service.AuditService,
) *Handler {
	return &Handler{
		authService:    authService,
		userService:    userService,
		productService: productService,
		auditService:   auditService,
	}
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{} "data: AuthResponse"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{} "data: AuthResponse"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshRequest true "Refresh Token Request"
// @Success 200 {object} map[string]interface{} "data: AuthResponse"
// @Failure 401 {object} map[string]string "error: Invalid refresh token"
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// Logout godoc
// @Summary Logout user
// @Description Revoke refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshRequest true "Refresh Token"
// @Success 200 {object} map[string]string "message: logged out successfully"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "data: User"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 404 {object} map[string]string "error: User not found"
// @Router /users/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	user, err := h.userService.GetProfile(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Update authenticated user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.UpdateUserRequest true "Update User Request"
// @Success 200 {object} map[string]interface{} "data: User"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Router /users/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(claims.UserID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change authenticated user's password
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} map[string]string "message: password changed successfully"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Router /users/me/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.ChangePassword(claims.UserID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// GetAllUsers godoc
// @Summary Get all users (Admin only)
// @Description Get paginated list of all users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "data, page, limit, total_items"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Router /users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.userService.GetAllUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        users,
		"page":        page,
		"limit":       limit,
		"total_items": total,
	})
}

// DeleteUser godoc
// @Summary Delete user (Admin only)
// @Description Delete user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "message: user deleted successfully"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// CreateProduct godoc
// @Summary Create product (Admin only)
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateProductRequest true "Create Product Request"
// @Success 201 {object} map[string]interface{} "data: Product"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Router /products [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get single product details
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "data: Product"
// @Failure 404 {object} map[string]string "error: Product not found"
// @Router /products/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// UpdateProduct godoc
// @Summary Update product (Admin only)
// @Description Update existing product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body domain.UpdateProductRequest true "Update Product Request"
// @Success 200 {object} map[string]interface{} "data: Product"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 404 {object} map[string]string "error: Product not found"
// @Router /products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// DeleteProduct godoc
// @Summary Delete product (Admin only)
// @Description Delete product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string "message: product deleted successfully"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Router /products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	if err := h.productService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get paginated and filtered list of products
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search keyword"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {object} map[string]interface{} "data: ProductListResponse"
// @Router /products [get]
func (h *Handler) GetAllProducts(c *gin.Context) {
	var query domain.ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.productService.GetAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetAuditLogs godoc
// @Summary Get audit logs (Admin only)
// @Description Get paginated and filtered audit logs
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param user_id query int false "Filter by user ID"
// @Param action query string false "Filter by action"
// @Param resource query string false "Filter by resource"
// @Success 200 {object} map[string]interface{} "data: AuditListResponse"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Router /audit-logs [get]
func (h *Handler) GetAuditLogs(c *gin.Context) {
	var query domain.AuditQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.auditService.GetAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify user email address using verification token
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]string "message: Email verified successfully"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 404 {object} map[string]string "error: Not Found"
// @Router /auth/verify-email [get]
func (h *Handler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	err := h.authService.VerifyEmail(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// ResendVerification godoc
// @Summary Resend verification email
// @Description Resend email verification link to user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.ResendVerificationRequest true "Resend Verification Request"
// @Success 200 {object} map[string]string "message: Verification email sent"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 404 {object} map[string]string "error: Not Found"
// @Router /auth/resend-verification [post]
func (h *Handler) ResendVerification(c *gin.Context) {
	var req domain.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ResendVerification(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset link to user's email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} map[string]string "message: Password reset email sent"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 404 {object} map[string]string "error: Not Found"
// @Router /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req domain.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ForgotPassword(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password using reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} map[string]string "message: Password reset successfully"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 404 {object} map[string]string "error: Not Found"
// @Router /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ResetPassword(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
