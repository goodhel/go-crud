package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goodhel/go-crud/controllers"
	"github.com/goodhel/go-crud/initializers"
	"github.com/goodhel/go-crud/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello this is api from Backend Go",
		})
	})

	r.GET("/posts", controllers.ListPosts)
	r.GET("/posts/:id", controllers.DetailPosts)
	r.POST("/posts", controllers.PostsCreate)
	r.PUT("/posts/:id", controllers.UpdatePost)
	r.DELETE("/posts", controllers.DeletePost)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.UserSession, controllers.Validate)

	r.Run() // listen and serve on 0.0.0.0:8081
}
