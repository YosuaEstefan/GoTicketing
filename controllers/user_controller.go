// controllers/auth_controller.go
package controllers

import (
	"net/http"
	"ticket/models"
	"ticket/service"

	"github.com/gin-gonic/gin"

)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type authController struct {
	userService service.UserService
}

func NewAuthController(userService service.UserService) AuthController {
	return &authController{
		userService: userService,
	}
}

func (ctrl *authController) Register(c *gin.Context) {
	var user models.User

	// Validate input
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	err := ctrl.userService.Register(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return user data (without password)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (ctrl *authController) Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate token
	token, err := ctrl.userService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
