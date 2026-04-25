package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanbekele/food-recipes-backend/gql"
	graphql "github.com/hasura/go-graphql-client"
)

type categoryRow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ListCategories returns categories using the Hasura admin client to avoid client permission issues.
func ListCategories(c *gin.Context) {
	var query struct {
		Categories []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `graphql:"categories(order_by: {name: asc})"`
	}
	if err := gql.Client.Query(c.Request.Context(), &query, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If empty, seed defaults then re-query
	if len(query.Categories) == 0 {
		type catObj struct {
			Name graphql.String `json:"name"`
			Slug graphql.String `json:"slug"`
		}
		objs := []catObj{
			{Name: "Breakfast", Slug: "breakfast"},
			{Name: "Lunch", Slug: "lunch"},
			{Name: "Dinner", Slug: "dinner"},
			{Name: "Dessert", Slug: "dessert"},
			{Name: "Vegan", Slug: "vegan"},
			{Name: "Vegetarian", Slug: "vegetarian"},
			{Name: "Gluten Free", Slug: "gluten-free"},
			{Name: "Keto", Slug: "keto"},
		}
		var m struct {
			InsertCategories struct {
				AffectedRows int `json:"affected_rows"`
			} `graphql:"insert_categories(objects: $objects)"`
		}
		vars := map[string]interface{}{"objects": objs}
		_ = gql.Client.Mutate(c.Request.Context(), &m, vars)
		// Re-query regardless of mutation result
		_ = gql.Client.Query(c.Request.Context(), &query, nil)
	}
	rows := make([]categoryRow, 0, len(query.Categories))
	for _, r := range query.Categories {
		rows = append(rows, categoryRow{ID: r.ID, Name: r.Name, Slug: r.Slug})
	}
	c.JSON(http.StatusOK, gin.H{"categories": rows})
}
