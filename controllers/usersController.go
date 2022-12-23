package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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
		Role:     []models.Role{{ID: 1}},
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

type MyData struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Roles []uint `json:"roles"`
}

// func indexOfMyData(id uint, data []MyData) int {
// 	for k, v := range data {
// 		if id == v.ID {
// 			return k
// 		}
// 	}

// 	return -1
// }

// IndexOfB is a function to find index of an element in a slice of map
func indexOfB(id uint, data []map[string]interface{}, variable string) int {
	for k, v := range data {
		if id == v[variable] {
			return k
		}
	}

	return -1
}

func ListUsers(c *gin.Context) {
	// var users []models.User

	type auser struct {
		ID     uint   `json:"id"`
		Email  string `json:"email"`
		Name   string `json:"name"`
		RoleID uint   `json:"role_id"`
	}

	var users []auser

	result := initializers.DB.Model(&models.User{}).Select("users.id, users.email, users.name, user_roles.role_id").Joins("LEFT JOIN user_roles ON user_roles.user_id = users.id").Scan(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Failed to get users",
		})

		return
	}

	// var testa = []MyData{}

	// for _, value := range users {
	// 	i := indexOfMyData(value.ID, testa)

	// 	if i == -1 {
	// 		testa = append(testa, MyData{
	// 			ID:    value.ID,
	// 			Email: value.Email,
	// 			Roles: []uint{value.RoleID},
	// 		})
	// 	} else {
	// 		testa[i].Roles = append(testa[i].Roles, value.RoleID)
	// 	}
	// }

	var testb []map[string]interface{}

	for _, value := range users {
		i := indexOfB(value.ID, testb, "id")
		if i == -1 {
			testb = append(testb, map[string]interface{}{
				"id":    value.ID,
				"email": value.Email,
				"roles": []uint{value.RoleID},
			})
		} else {
			testb[i]["roles"] = append(testb[i]["roles"].([]uint), value.RoleID)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   testb,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Check User
	var user models.User
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "User not found",
		})

		return
	}

	fmt.Println("deteail user", user)

	// Get File from Body
	file, err := c.FormFile("file")

	var fileName *string

	if err == nil {
		s := strings.Split(file.Filename, ".")
		ext := s[len(s)-1]
		randName := fmt.Sprintf("%v", time.Now().UnixMilli()) + "_file." + ext

		if err := c.SaveUploadedFile(file, "./public/"+randName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"error":  err.Error(),
			})
			return
		}

		filepath := "http://localhost:8080" + "/public/" + randName
		fileName = &filepath

		// Delete Old File
		av := strings.Split(*user.Avatar, "/")
		oldFile := "./public/" + av[len(av)-1]
		os.Remove(oldFile)
	} else {
		fileName = nil
	}

	// Get Data from Body
	var body struct {
		Name string `form:"name";json:"name"`
	}

	c.Bind(&body)

	var input struct {
		Name   string  `json:"name"`
		Avatar *string `json:"avatar"`
	}

	input.Name = body.Name
	input.Avatar = fileName

	// Update User
	initializers.DB.Model(&user).Updates(models.User{
		Name:   input.Name,
		Avatar: input.Avatar,
	})

	// Return Response
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   user,
	})
}
