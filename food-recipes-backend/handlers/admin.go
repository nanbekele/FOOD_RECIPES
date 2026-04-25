package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/nanbekele/food-recipes-backend/gql"
	"golang.org/x/crypto/bcrypt"
)

func PublicListChefs(c *gin.Context) {
	verifiedParam := strings.ToLower(strings.TrimSpace(c.Query("verified")))
	type userRow struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Role       string `json:"role"`
		IsVerified bool   `json:"is_verified" graphql:"is_verified"`
	}

	var query struct {
		Users []userRow `graphql:"users(where: { role: { _eq: \"chef\" } })"`
	}
	if err := gql.Client.Query(c.Request.Context(), &query, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows := make([]userRow, 0, len(query.Users))
	for _, u := range query.Users {
		if verifiedParam == "true" && !u.IsVerified {
			continue
		}
		if verifiedParam == "false" && u.IsVerified {
			continue
		}
		rows = append(rows, u)
	}
	c.JSON(http.StatusOK, gin.H{"chefs": rows})
}

// AdminListChefs returns chefs, optionally filtered by is_verified via query param: verified=true|false
func AdminListChefs(c *gin.Context) {
	verifiedParam := strings.ToLower(strings.TrimSpace(c.Query("verified")))
	type userRow struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Role       string `json:"role"`
		IsVerified bool   `json:"is_verified" graphql:"is_verified"`
	}

	var query struct {
		Users []userRow `graphql:"users(where: { role: { _eq: \"chef\" } })"`
	}
	if err := gql.Client.Query(c.Request.Context(), &query, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows := make([]userRow, 0, len(query.Users))
	for _, u := range query.Users {
		if verifiedParam == "true" && !u.IsVerified {
			continue
		}
		if verifiedParam == "false" && u.IsVerified {
			continue
		}
		rows = append(rows, u)
	}
	c.JSON(http.StatusOK, gin.H{"chefs": rows})
}

// AdminVerifyChef sets is_verified for a chef by id. Body: { "is_verified": true/false }
func AdminVerifyChef(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		IsVerified bool `json:"is_verified"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var mutation struct {
		UpdateUser struct {
			ID int `json:"id"`
		} `graphql:"update_users_by_pk(pk_columns: {id: $id}, _set: {is_verified: $is_verified})"`
	}
	vars := map[string]interface{}{
		"id":          graphql.Int(id),
		"is_verified": graphql.Boolean(body.IsVerified),
	}
	if err := gql.Client.Mutate(c.Request.Context(), &mutation, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "is_verified": body.IsVerified})
}

// AdminResetPassword allows an admin to set a new password for a user by id.
// Body: { "new_password": string }
func AdminResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(strings.TrimSpace(body.NewPassword)) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be at least 6 characters"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}
	var m struct {
		UpdateUser struct {
			ID int `json:"id"`
		} `graphql:"update_users_by_pk(pk_columns: {id: $id}, _set: {password: $password})"`
	}
	vars := map[string]interface{}{
		"id":       graphql.Int(id),
		"password": graphql.String(string(hash)),
	}
	if err := gql.Client.Mutate(c.Request.Context(), &m, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Password updated"})
}
