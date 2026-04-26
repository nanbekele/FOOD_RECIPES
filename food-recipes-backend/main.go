package main

import (
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nanbekele/food-recipes-backend/config"
	"github.com/nanbekele/food-recipes-backend/db"
	"github.com/nanbekele/food-recipes-backend/gql"
	"github.com/nanbekele/food-recipes-backend/handlers"
	"github.com/nanbekele/food-recipes-backend/middleware"
)

func main() {
	// Load .env file
	config.LoadEnv()

	// Read Hasura config from environment
	hasuraEndpoint := config.GetEnv("HASURA_GRAPHQL_ENDPOINT", "http://localhost:8082/v1/graphql")
	hasuraAdminSecret := config.GetEnv("HASURA_GRAPHQL_ADMIN_SECRET", "")

	// Run SQL migrations against Postgres
	dsn := config.GetEnv("POSTGRES_DSN", "postgres://postgres:mysecretpassword@localhost:5433/foodrecipes?sslmode=disable")
	if err := db.RunMigrations(dsn); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	// Initialize Hasura client with endpoint and admin secret
	gql.InitHasuraClient(hasuraEndpoint, hasuraAdminSecret)

	// Set up routes
	r := gin.Default()
	// Restrict CORS
	corsEnv := config.GetEnv("CORS_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000")
	origins := []string{}
	for _, o := range strings.Split(corsEnv, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/upload", handlers.UploadImage)
	r.GET("/categories", handlers.ListCategories)
	r.GET("/chefs", handlers.PublicListChefs)
	r.GET("/news", handlers.PublicListNews)
	r.GET("/users/:id", handlers.PublicGetUser)
	r.GET("/users/:id/recipes", handlers.PublicListUserRecipes)
	r.POST("/news", middleware.Authenticated(), handlers.CreateNews)
	r.POST("/password/change", middleware.Authenticated(), handlers.ChangePassword)
	r.GET("/users/:id/following", middleware.Authenticated(), handlers.GetFollowStatus)
	r.POST("/users/:id/follow", middleware.Authenticated(), handlers.FollowUser)
	r.DELETE("/users/:id/follow", middleware.Authenticated(), handlers.UnfollowUser)
	r.GET("/admin/chefs", middleware.RequireRole("admin"), handlers.AdminListChefs)
	r.PATCH("/admin/chefs/:id/verify", middleware.RequireRole("admin"), handlers.AdminVerifyChef)
	r.PATCH("/admin/users/:id/password", middleware.RequireRole("admin"), handlers.AdminResetPassword)

	// Run server
	log.Println("Server is running on http://localhost:8081")
	_ = r.Run(":8081")
}
