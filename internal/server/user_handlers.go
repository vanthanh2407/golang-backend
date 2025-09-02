package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserRequest represents the request body for user operations
type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateUserRequest represents the request body for updating user
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// UpdatePasswordRequest represents the request body for updating password
type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

// CreateUserHandler handles user creation
func (s *Server) CreateUserHandler(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Check if user already exists
	existingUser, _ := s.db.GetUserByEmail(c.Request.Context(), req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	existingUser, _ = s.db.GetUserByUsername(c.Request.Context(), req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this username already exists",
		})
		return
	}

	// Create user
	user, err := s.db.CreateUser(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUserHandler handles getting a user by ID
func (s *Server) GetUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := s.db.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetAllUsersHandler handles getting all users
func (s *Server) GetAllUsersHandler(c *gin.Context) {
	users, err := s.db.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// UpdateUserHandler handles updating user information
func (s *Server) UpdateUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Check if user exists
	existingUser, err := s.db.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Check if new email is already taken by another user
	if req.Email != existingUser.Email {
		userWithEmail, _ := s.db.GetUserByEmail(c.Request.Context(), req.Email)
		if userWithEmail != nil && userWithEmail.ID != id {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already taken by another user",
			})
			return
		}
	}

	// Check if new username is already taken by another user
	if req.Username != existingUser.Username {
		userWithUsername, _ := s.db.GetUserByUsername(c.Request.Context(), req.Username)
		if userWithUsername != nil && userWithUsername.ID != id {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Username already taken by another user",
			})
			return
		}
	}

	// Update user
	user, err := s.db.UpdateUser(c.Request.Context(), id, req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

// UpdatePasswordHandler handles updating user password
func (s *Server) UpdatePasswordHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Check if user exists
	_, err = s.db.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Update password
	err = s.db.UpdateUserPassword(c.Request.Context(), id, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update password: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// DeleteUserHandler handles deleting a user
func (s *Server) DeleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = s.db.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
