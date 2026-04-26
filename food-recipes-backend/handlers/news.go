package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nanbekele/food-recipes-backend/gql"
	graphql "github.com/hasura/go-graphql-client"
)

type newsRow struct {
	ID        string `json:"id"`
	UserID    int    `json:"user_id" graphql:"user_id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Tag       string `json:"tag"`
	CreatedAt string `json:"created_at" graphql:"created_at"`

	Author struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Role       string `json:"role"`
		IsVerified bool   `json:"is_verified" graphql:"is_verified"`
	} `json:"author" graphql:"user"`
}

// PublicListNews returns news articles ordered by newest first.
func PublicListNews(c *gin.Context) {
	var q struct {
		News []newsRow `graphql:"news(order_by: { created_at: desc })"`
	}
	if err := gql.Client.Query(c.Request.Context(), &q, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"news": q.News})
}

// CreateNews creates a news article for the authenticated user.
func CreateNews(c *gin.Context) {
	uidAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uidStr, _ := uidAny.(string)
	uid, _ := strconv.Atoi(uidStr)
	if uid <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var body struct {
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Content string `json:"content"`
		Image   string `json:"image"`
		Tag     string `json:"tag"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(body.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	var m struct {
		Insert struct {
			ID string `json:"id"`
		} `graphql:"insert_news_one(object: {user_id: $uid, title: $title, summary: $summary, content: $content, image: $image, tag: $tag})"`
	}
	vars := map[string]interface{}{
		"uid":     graphql.Int(uid),
		"title":   graphql.String(strings.TrimSpace(body.Title)),
		"summary": graphql.String(strings.TrimSpace(body.Summary)),
		"content": graphql.String(strings.TrimSpace(body.Content)),
		"image":   graphql.String(strings.TrimSpace(body.Image)),
		"tag":     graphql.String(strings.TrimSpace(body.Tag)),
	}
	if err := gql.Client.Mutate(c.Request.Context(), &m, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": m.Insert.ID})
}
