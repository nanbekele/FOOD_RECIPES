package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nanbekele/food-recipes-backend/gql"
	graphql "github.com/hasura/go-graphql-client"
)

type publicUserRow struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Bio        string `json:"bio"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified" graphql:"is_verified"`
}

// PublicGetUser returns a limited public view of a user.
func PublicGetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var q struct {
		User *publicUserRow `graphql:"users_by_pk(id: $id)"`
	}
	vars := map[string]interface{}{"id": graphql.Int(id)}
	if err := gql.Client.Query(c.Request.Context(), &q, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if q.User == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": q.User})
}

type publicRecipeRow struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	FeaturedImage string  `json:"featured_image" graphql:"featured_image"`
	AverageRating float64 `json:"average_rating" graphql:"average_rating"`
	FavoritesAgg  struct {
		Aggregate struct {
			Count int `json:"count"`
		} `json:"aggregate"`
	} `json:"favorites_aggregate" graphql:"favorites_aggregate"`
}

// PublicListUserRecipes returns recipes authored by a user.
func PublicListUserRecipes(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var q struct {
		Recipes []publicRecipeRow `graphql:"recipes(where: { user_id: { _eq: $uid } }, order_by: { created_at: desc })"`
	}
	vars := map[string]interface{}{"uid": graphql.Int(id)}
	if err := gql.Client.Query(c.Request.Context(), &q, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recipes": q.Recipes})
}

type followRow struct {
	ID int `json:"id"`
}

// FollowUser creates a follow relationship from the authenticated viewer to the target user.
func FollowUser(c *gin.Context) {
	uidAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	viewerID, _ := strconv.Atoi(uidAny.(string))

	idStr := c.Param("id")
	targetID, err := strconv.Atoi(idStr)
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if viewerID == targetID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot follow yourself"})
		return
	}

	var m struct {
		Insert struct {
			Returning []followRow `json:"returning"`
		} `graphql:"insert_follows(objects: { follower_id: $follower, followee_id: $followee }, on_conflict: {constraint: follows_follower_id_followee_id_key, update_columns: []})"`
	}
	vars := map[string]interface{}{
		"follower": graphql.Int(viewerID),
		"followee": graphql.Int(targetID),
	}
	if err := gql.Client.Mutate(c.Request.Context(), &m, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UnfollowUser deletes a follow relationship from the authenticated viewer to the target user.
func UnfollowUser(c *gin.Context) {
	uidAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	viewerID, _ := strconv.Atoi(uidAny.(string))

	idStr := c.Param("id")
	targetID, err := strconv.Atoi(idStr)
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var m struct {
		Delete struct {
			AffectedRows int `json:"affected_rows"`
		} `graphql:"delete_follows(where: { follower_id: { _eq: $follower }, followee_id: { _eq: $followee } })"`
	}
	vars := map[string]interface{}{
		"follower": graphql.Int(viewerID),
		"followee": graphql.Int(targetID),
	}
	if err := gql.Client.Mutate(c.Request.Context(), &m, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetFollowStatus returns whether the authenticated viewer follows the target user.
func GetFollowStatus(c *gin.Context) {
	uidAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	viewerID, _ := strconv.Atoi(uidAny.(string))

	idStr := c.Param("id")
	targetID, err := strconv.Atoi(idStr)
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var q struct {
		Follows []struct{
			ID int `json:"id"`
		} `graphql:"follows(where: { follower_id: { _eq: $follower }, followee_id: { _eq: $followee } }, limit: 1)"`
	}
	vars := map[string]interface{}{
		"follower": graphql.Int(viewerID),
		"followee": graphql.Int(targetID),
	}
	if err := gql.Client.Query(c.Request.Context(), &q, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"following": len(q.Follows) > 0})
}
