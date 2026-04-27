package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nanbekele/food-recipes-backend/gql"
)

// POST /recipes/:id/favorite — add a recipe to the authenticated user's favorites
func AddFavorite(c *gin.Context) {
	uidAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	viewerID, _ := strconv.Atoi(uidAny.(string))

	recipeID := c.Param("id")
	if recipeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipe id required"})
		return
	}

	query := `mutation AddFav($uid: Int!, $rid: uuid!) {
		insert_favorites_one(object: { user_id: $uid, recipe_id: $rid }, on_conflict: {constraint: uq_fav_user_recipe, update_columns: []}) {
			id
		}
	}`

	var m struct {
		InsertFavoritesOne struct {
			ID string `graphql:"id"`
		} `graphql:"insert_favorites_one"`
	}
	vars := map[string]interface{}{
		"uid": viewerID,
		"rid": recipeID,
	}
	if err := gql.Client.Exec(c.Request.Context(), query, &m, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": m.InsertFavoritesOne.ID})
}

// DELETE /recipes/:id/favorite — remove a recipe from the authenticated user's favorites
func RemoveFavorite(c *gin.Context) {
	uidAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	viewerID, _ := strconv.Atoi(uidAny.(string))

	recipeID := c.Param("id")
	if recipeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipe id required"})
		return
	}

	query := `mutation RemoveFav($uid: Int!, $rid: uuid!) {
		delete_favorites(where: { user_id: { _eq: $uid }, recipe_id: { _eq: $rid } }) {
			affected_rows
		}
	}`

	var m struct {
		DeleteFavorites struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"delete_favorites"`
	}
	vars := map[string]interface{}{
		"uid": viewerID,
		"rid": recipeID,
	}
	if err := gql.Client.Exec(c.Request.Context(), query, &m, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "affected_rows": m.DeleteFavorites.AffectedRows})
}
