package handlers

import (
	"context"
	"go-webserver-performance-test/models/data"
	"go-webserver-performance-test/services"
	"go-webserver-performance-test/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (m *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := m.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
func (m *UserHandler) CheckRequiredRole(ctx *gin.Context) {
	accessToken := ctx.GetHeader("Authorization")
	if accessToken == "" && !strings.HasPrefix(accessToken, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := utils.DecodeToken(accessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	requiredRoles := []string{"Admin", "Super Admin"}
	for _, role := range requiredRoles {
		if !strings.Contains(claims.Roles, role) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	}

	ctx.Next()
}

func (m *UserHandler) LogIn(c *gin.Context) {
	// get username and password from the request
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind JSON request body to the struct
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user *data.User
	// check if username includes a '@' symbol
	if strings.Contains(loginData.Username, "@") {
		var err error
		user, err = m.userService.GetUserByEmail(loginData.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
	} else {
		var err error
		user, err = m.userService.GetUserByUsername(loginData.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
	}

	// check if the password is correct
	if !utils.ComparePassword(user.PasswordHash, loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Login successful, create JWT tokens
	access, refresh, err := utils.CreateTokens(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func (m *UserHandler) GetUser(c *gin.Context) {
	// The user has to be authenticated to get their own profile
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" && !strings.HasPrefix(accessToken, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := utils.DecodeToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := m.userService.GetUserByID(claims.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (m *UserHandler) CreateUser(c *gin.Context) {
	var user data.User
	var userCreateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&userCreateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// hash the password before storing it
	hashedPassword, err := utils.HashPassword(userCreateData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.ID = newUUID
	user.Username = userCreateData.Username
	user.Email = userCreateData.Email
	user.PasswordHash = hashedPassword
	id, err := m.userService.CreateUser(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// TODO: Add auth to only let the user update their own profile
func (m *UserHandler) UpdateUser(c *gin.Context) {
	// The user has to be authenticated to get their own profile
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" && !strings.HasPrefix(accessToken, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := utils.DecodeToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user data.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user.ID, err = uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := m.userService.UpdateUser(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// TOOD: Only let the user delete their own profile
func (m *UserHandler) DeleteUser(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" && !strings.HasPrefix(accessToken, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := utils.DecodeToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := m.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
