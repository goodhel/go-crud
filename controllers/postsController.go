package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goodhel/go-crud/initializers"
	"github.com/goodhel/go-crud/models"
)

func PostsCreate(c *gin.Context) {
	// Get data from request
	var body struct {
		Body   string `json:"body" binding:"required"`
		Titile string `json:"title" binding:"required"`
	}

	c.Bind(&body)

	// Create post
	post := models.Post{Title: body.Titile, Body: body.Body}

	result := initializers.DB.Create(&post) // pass pointer of data to Create

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": true,
		"data":   post,
	})
}

func ListPosts(c *gin.Context) {
	var posts []models.Post

	result := initializers.DB.Find(&posts)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"status": true,
		"data":   posts,
	})
}

func DetailPosts(c *gin.Context) {
	var post models.Post

	result := initializers.DB.First(&post, c.Param("id"))

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"status": true,
		"data":   post,
	})
}

func UpdatePost(c *gin.Context) {
	// Get id from the url
	id := c.Param("id")

	// Get data from request
	var body struct {
		Body   string `json:"body" binding:"required"`
		Titile string `json:"title" binding:"required"`
	}

	c.Bind(&body)

	// Check Post is exist or not
	var post models.Post
	result := initializers.DB.First(&post, id)

	if result.Error != nil {
		c.Status(404)
		return
	}

	// Update post
	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Titile,
		Body:  body.Body,
	})

	// Return response
	c.JSON(200, gin.H{
		"status": true,
		"data":   post,
	})
}

func DeletePost(c *gin.Context) {
	// Get id from the url with query param
	id := c.Query("id")

	// Check Post is exist or not
	var post models.Post
	result := initializers.DB.First(&post, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Post not found",
		})

		return
	}

	// Delete post
	initializers.DB.Delete(&post)

	// Return Response
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   post,
	})
}
