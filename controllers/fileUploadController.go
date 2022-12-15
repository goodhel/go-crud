package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goodhel/go-crud/initializers"
	"github.com/goodhel/go-crud/models"
)

func UploadFile(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "File not found",
		})
		return
	}

	fileName := file.Filename
	s := strings.Split(fileName, ".")
	ext := s[len(s)-1]
	mime := file.Header.Values("Content-Type")[0]

	randName := fmt.Sprintf("%v", time.Now().UnixMilli()) + "_file." + ext

	if err := c.SaveUploadedFile(file, "./uploads/"+randName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	// using the function
	mydir, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
	}

	// Insert into database
	fileupload := models.FileUpload{
		Name:         randName,
		OriginalName: fileName,
		Mime:         mime,
		Path:         mydir + "/uploads/" + randName,
		Extension:    ext,
	}

	result := initializers.DB.Create(&fileupload)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Error save file to database",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": fileupload,
	})
}
