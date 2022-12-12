package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/goodhel/go-crud/initializers"
	"github.com/goodhel/go-crud/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	// Get Email and Password grom Body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Failed to read body",
		})

		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Failed to hash password",
		})

		return
	}

	// Create User
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Failed to create user",
		})

		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"status": true,
		"data":   user,
	})
}

func Login(c *gin.Context) {
	// Get Email and Password grom Body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Failed to read body",
		})

		return
	}

	// Get User
	var user models.User

	result := initializers.DB.Where("email = ?", body.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Failed to get user",
		})

		return
	}

	// Compare Password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"error":  "Sorry, wrong password",
		})

		return
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "Failed to generate token",
		})

		return
	}

	// Set Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	// Response
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"token":     tokenString,
			"expiredAt": time.Now().Add(time.Hour * 24),
		},
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   user,
	})
}
