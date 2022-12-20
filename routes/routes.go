package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/goodhel/go-crud/controllers"
	"github.com/goodhel/go-crud/middleware"
)

func Routes() *gin.Engine {
	router := gin.Default()

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// Grouping routes under path /api
	api := router.Group("/api")
	{
		api.GET("/posts", controllers.ListPosts)
		api.GET("/posts/:id", controllers.DetailPosts)
		api.POST("/posts", controllers.PostsCreate)
		api.PUT("/posts/:id", controllers.UpdatePost)
		api.DELETE("/posts", controllers.DeletePost)
		api.GET("/validate", middleware.UserSession, controllers.Validate)
		api.GET("/users", middleware.UserSession, controllers.ListUsers)
		api.GET("/file/:id", controllers.DonwloadFile)
		api.POST("/file/upload", middleware.UserSession, controllers.UploadFile)
	}

	return router
}
