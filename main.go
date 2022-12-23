package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goodhel/go-crud/initializers"
	"github.com/goodhel/go-crud/routes"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	router := routes.Routes()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello this is api from Backend Go",
		})
	})

	router.Static("/public", "./public")

	router.Run() // listen and serve on 0.0.0.0:8081
}
