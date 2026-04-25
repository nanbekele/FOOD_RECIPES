package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nanbekele/food-recipes-backend/config"
)

func parseJWT(c *gin.Context) (jwt.MapClaims, error) {
	secret := []byte(config.GetEnv("JWT_SECRET", "myhasurasecret"))
	h := c.GetHeader("Authorization")
	if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return nil, fmt.Errorf("missing token")
	}
	tok := strings.TrimSpace(h[len("Bearer "):])
	token, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

// Authenticated ensures there is a valid JWT and exposes user_id and role in context
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := parseJWT(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if hclaim, ok := claims["https://hasura.io/jwt/claims"].(map[string]interface{}); ok {
			if uid, ok2 := hclaim["x-hasura-user-id"].(string); ok2 {
				c.Set("user_id", uid)
			}
			if r, ok2 := hclaim["x-hasura-role"].(string); ok2 {
				c.Set("role", r)
			}
		}
		c.Next()
	}
}

// RequireRole ensures the JWT role matches the required role
func RequireRole(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := parseJWT(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		role := ""
		if hclaim, ok := claims["https://hasura.io/jwt/claims"].(map[string]interface{}); ok {
			if r, ok2 := hclaim["x-hasura-role"].(string); ok2 {
				role = r
			}
			if uid, ok2 := hclaim["x-hasura-user-id"].(string); ok2 {
				c.Set("user_id", uid)
			}
		}
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role missing"})
			return
		}
		if role != required {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
			return
		}
		c.Next()
	}
}
