package main

import (
	"github.com/goodhel/go-crud/initializers"
	"github.com/goodhel/go-crud/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Post{}, &models.User{})
}
