package swagger

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Here's an example of how to document your API handlers with Swagger annotations
// These comments will be automatically parsed and included in your Swagger docs

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "User found"
// @Success 400 {object} ErrorResponse "Bad request"
// @Success 404 {object} ErrorResponse "User not found"
// @Success 500 {object} ErrorResponse "Internal server error"
// @Router /v1/users/{id} [get]
func GetUserByID(c *gin.Context) {
	// Implementation would go here
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":        id,
		"username":  "example_user",
		"full_name": "Example User",
		"email":     "user@example.com",
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRequest true "User information"
// @Success 201 {object} UserResponse "User created"
// @Success 400 {object} ErrorResponse "Bad request"
// @Success 500 {object} ErrorResponse "Internal server error"
// @Router /v1/users [post]
func CreateUser(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusCreated, gin.H{
		"id":        "new-user-123",
		"username":  "new_user",
		"full_name": "New User",
		"email":     "newuser@example.com",
	})
}

// Note: These are just example structs for documentation purposes
// In a real application, these would be defined in your models package

// UserRequest represents a request to create a user
type UserRequest struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserResponse represents a user in the system
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
