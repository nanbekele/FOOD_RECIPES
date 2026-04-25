package utils

import (
	"time"

	"github.com/nanbekele/food-recipes-backend/config"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(config.GetEnv("JWT_SECRET", "myhasurasecret"))

func GenerateJWT(userID, role string) (string, error) {
	// Determine allowed roles based on user's role
	allowed := []string{"public"}
	switch role {
	case "admin":
		allowed = []string{"public", "guest", "chef", "admin"}
	case "chef":
		allowed = []string{"public", "guest", "chef"}
	case "guest", "":
		role = "guest"
		allowed = []string{"public", "guest"}
	default:
		// fallback to guest
		role = "guest"
		allowed = []string{"public", "guest"}
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-user-id":       userID,
			"x-hasura-allowed-roles": allowed,
			"x-hasura-default-role":  role,
			"x-hasura-role":          role,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
