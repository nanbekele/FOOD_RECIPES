package handlers

import (
	"net/http"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/nanbekele/food-recipes-backend/gql"
	"github.com/nanbekele/food-recipes-backend/utils"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

// ChangePassword allows an authenticated user to update their password
// Body: { "old_password": string, "new_password": string }
func ChangePassword(c *gin.Context) {
	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(body.NewPassword) == "" || len(body.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be at least 6 characters"})
		return
	}
	uidAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uidStr, _ := uidAny.(string)
	idInt, _ := strconv.Atoi(uidStr)

	// Fetch current user password
	var q struct {
		User *struct {
			ID       int    `json:"id"`
			Password string `json:"password"`
		} `graphql:"users_by_pk(id: $id)"`
	}
	if err := gql.Client.Query(c.Request.Context(), &q, map[string]interface{}{"id": graphql.Int(idInt)}); err != nil || q.User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(q.User.Password), []byte(body.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	var m struct {
		Update struct {
			ID int `json:"id"`
		} `graphql:"update_users_by_pk(pk_columns: {id: $id}, _set: {password: $password})"`
	}
	vars := map[string]interface{}{
		"id":       graphql.Int(idInt),
		"password": graphql.String(string(hash)),
	}
	if err := gql.Client.Mutate(c.Request.Context(), &m, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed"})
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Normalize email
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	// Check if email already exists
	var existing struct {
		Users []struct {
			ID int `json:"id"`
		} `graphql:"users(where: {email: {_ilike: $email}})"`
	}
	if err := gql.Client.Query(c.Request.Context(), &existing, map[string]interface{}{"email": graphql.String(input.Email)}); err == nil {
		if len(existing.Users) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
			return
		}
	}
	if strings.TrimSpace(input.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Determine role: only allow "guest" or "chef" from signup; everything else becomes "guest".
	role := strings.ToLower(strings.TrimSpace(input.Role))
	switch role {
	case "chef":
		// keep as chef
	case "guest", "":
		role = "guest"
	default:
		role = "guest"
	}

	// GraphQL mutation for user insert
	var mutation struct {
		InsertUser struct {
			ID int `json:"id"`
		} `graphql:"insert_users_one(object: {name: $name, email: $email, password: $password, role: $role})"`
	}
	vars := map[string]interface{}{
		"name":     graphql.String(input.Name),
		"email":    graphql.String(input.Email),
		"password": graphql.String(string(hashedPassword)),
		"role":     graphql.String(role),
	}
	err = gql.Client.Mutate(c.Request.Context(), &mutation, vars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered", "id": mutation.InsertUser.ID})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Normalize email
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	// GraphQL query for user by email
	var queryFull struct {
		Users []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		} `graphql:"users(where: {email: {_ilike: $email}})"`
	}
	vars := map[string]interface{}{"email": graphql.String(input.Email)}
	err := gql.Client.Query(c.Request.Context(), &queryFull, vars)
	if err != nil {
		// Fallback: try minimal selection without role
		var queryMinimal struct {
			Users []struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Email    string `json:"email"`
				Password string `json:"password"`
			} `graphql:"users(where: {email: {_ilike: $email}})"`
		}
		if err2 := gql.Client.Query(c.Request.Context(), &queryMinimal, vars); err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service unavailable"})
			return
		}
		if len(queryMinimal.Users) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		user := queryMinimal.Users[0]
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		role := "guest"
		token, err := utils.GenerateJWT(strconv.Itoa(user.ID), role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"name":  user.Name,
			"email": user.Email,
			"id":    user.ID,
			"role":  role,
		})
		return
	}
	if len(queryFull.Users) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	user := queryFull.Users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT using user's role (default guest)
	role := user.Role
	if role == "" {
		role = "guest"
	}
	token, err := utils.GenerateJWT(strconv.Itoa(user.ID), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"name":  user.Name,
		"email": user.Email,
		"id":    user.ID,
		"role":  role,
	})
}
