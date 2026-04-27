package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nanbekele/food-recipes-backend/gql"
)

// POST /recipes/:id/rate — add or update a recipe rating
func RateRecipe(c *gin.Context) {
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

	var body struct {
		Rating int `json:"rating"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Rating < 1 || body.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating must be between 1 and 5"})
		return
	}

	query := `mutation Rate($uid: Int!, $rid: uuid!, $rating: Int!) {
		insert_recipe_ratings_one(
			object: { user_id: $uid, recipe_id: $rid, rating: $rating },
			on_conflict: {constraint: uq_user_recipe_rating, update_columns: [rating]}
		) {
			id
		}
	}`

	var m struct {
		InsertRecipeRatingsOne struct {
			ID int `graphql:"id"`
		} `graphql:"insert_recipe_ratings_one"`
	}
	vars := map[string]interface{}{
		"uid":    viewerID,
		"rid":    recipeID,
		"rating": body.Rating,
	}

	if err := gql.Client.Exec(c.Request.Context(), query, &m, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": m.InsertRecipeRatingsOne.ID})
}
