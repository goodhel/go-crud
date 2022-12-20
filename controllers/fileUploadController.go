package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goodhel/go-crud/initializers"
	"github.com/goodhel/go-crud/models"
	"github.com/xuri/excelize/v2"
)

func UploadFile(c *gin.Context) {
	// single file
	user, token := c.Get("user") // get user from middleware
	if !token {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"error":  "Not Authenticated, no token provided",
		})
		return
	}

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "File not found",
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

	filepath := mydir + "/uploads/" + randName

	// Insert into database
	fileupload := models.FileUpload{
		Name:         randName,
		OriginalName: fileName,
		Mime:         mime,
		Path:         filepath,
		Extension:    ext,
		UserID:       user.(models.User).ID,
	}

	result := initializers.DB.Create(&fileupload)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Error save file to database",
		})

		return
	}

	f, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	var data []models.DWorkOrder

	for index, row := range rows {
		// Exclude Header (Row 1)
		if index != 0 {
			qty, e := strconv.Atoi(row[8])
			if e != nil {
				qty = 0
			}
			totalOrder, e := strconv.Atoi(row[9])
			if e != nil {
				totalOrder = 0
			}
			totalBox, e := strconv.Atoi(row[10])
			if e != nil {
				totalBox = 0
			}

			// Convert Date String From Excel to Match type Database
			prodParse, error := time.Parse("02/01/2006", row[2])
			if error != nil {
				prodParse = time.Now()
			}
			delivParse, error := time.Parse("02/01/2006", row[3])
			if error != nil {
				delivParse = time.Now()
			}

			prodDate, err := time.Parse("2006-01-02", prodParse.Format("2006-01-02"))
			if err != nil {
				prodDate = time.Now()
			}

			delivDate, err := time.Parse("2006-01-02", delivParse.Format("2006-01-02"))
			if err != nil {
				delivDate = time.Now()
			}

			data = append(data, models.DWorkOrder{
				NoWorkOrder:  row[1],
				Customer:     row[5],
				PartNumber:   row[6],
				PartName:     row[7],
				Qty:          qty,
				TotalOrder:   totalOrder,
				TotalBox:     totalBox,
				ProdDate:     &prodDate,
				DelivDate:    &delivDate,
				FileUploadID: fileupload.ID,
			})
		}
	}

	initializers.DB.CreateInBatches(&data, len(data))

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   data,
	})
}

func DonwloadFile(c *gin.Context) {
	var fileupload models.FileUpload
	id := c.Param("id")

	result := initializers.DB.First(&fileupload, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Data UploadFile not found",
		})
	}

	content, err := os.ReadFile(fileupload.Path)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "File not found",
		})
	}

	// Alternate 1
	c.Header("Content-Disposition", "attachment; filename="+fileupload.Name)
	c.Data(http.StatusOK, fileupload.Mime, content)

	// Alternate 2
	// c.FileAttachment(fileupload.Path, fileupload.Name)
}
