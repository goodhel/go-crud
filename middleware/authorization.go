package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/goodhel/go-crud/initializers"
	"github.com/goodhel/go-crud/models"
)

func UserSession(c *gin.Context) {
	fmt.Println("middleware")
	// Get the cookie from request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"status": false,
			"error":  "Unauthorized No Token",
		})

		return
	}

	// Decode and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": false,
				"error":  "Token Expired",
			})

			return
		}

		// Find the user with token
		var user models.User
		result := initializers.DB.Where("id = ?", claims["id"]).First(&user)

		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": false,
				"error":  "Failed to get user",
			})

			return
		}

		// Attach to request
		c.Set("user", user)

		// Continue to next middleware
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"error":  "Unauthorized Token Not Valid",
		})
	}
}
