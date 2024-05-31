package handlers

import (
	"fmt"
	"ginapp/domain"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Dashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cars []domain.Car

		if err := db.Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the database"})
			return

		}
		c.HTML(http.StatusOK, "admin.html", gin.H{"Cars": cars})

	}

}

func DeleteCar(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		id := c.Param("id")
		fmt.Println("here is the id", id)
		carID, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid id")
			return
		}
		if err := db.Where("id=?", carID).Delete(&domain.Car{}).Error; err != nil {
			c.String(http.StatusInternalServerError, "failed to delete the database")
			return

		}
		c.Redirect(http.StatusSeeOther, "/admin")

	}

}
func EditCar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		var car domain.Car

		if err := db.First(&car, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
			return

		}
		car.Brand = c.PostForm("brand")
		car.Model = c.PostForm("model")
		car.Year = c.PostForm("year")
		car.Color = c.PostForm("color")
		car.Variant = c.PostForm("variant")
		car.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		car.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		car.Transmission = c.PostForm("transmission")
		car.RegNo = c.PostForm("regno")
		car.Status = c.PostForm("status")

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get the form"})
			return
		}
		files := form.File["images[]"]
		fmt.Println("here is the files", files)

		imagePaths := map[int]string{

			0: car.ImagePath1,
			1: car.ImagePath2,
			2: car.ImagePath3,
			3: car.ImagePath4,
			4: car.ImagePath5,
			5: car.ImagePath6,
			6: car.ImagePath7,
		}

		for i, file := range files {
			if file != nil {
				filename := filepath.Base(file.Filename)
				uploadPath := filepath.Join("uploads", filename)
				if err := c.SaveUploadedFile(file, uploadPath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the image"})
					return
				}
				imagePaths[i] = "/" + strings.ReplaceAll(uploadPath, "\\", "/")

			}

		}
		car.ImagePath1 = imagePaths[0]
		car.ImagePath2 = imagePaths[1]
		car.ImagePath3 = imagePaths[2]
		car.ImagePath4 = imagePaths[3]
		car.ImagePath5 = imagePaths[4]
		car.ImagePath6 = imagePaths[5]
		car.ImagePath7 = imagePaths[6]

		if err := db.Save(&car).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload the car"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin")
	}

}

func AddCar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var car domain.Car

		// Retrieve form data
		car.Brand = c.PostForm("brand")
		car.Model = c.PostForm("model")
		car.Year = c.PostForm("year")
		car.Color = c.PostForm("color")
		car.Variant = c.PostForm("variant")
		car.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		car.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		car.Transmission = c.PostForm("transmission")
		car.RegNo = c.PostForm("regno")
		car.Status = c.PostForm("status")

		// Retrieve multiple files from the request
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get the form"})
			return
		}
		files := form.File["images[]"]

		// Ensure exactly 7 images are uploaded
		if len(files) != 7 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "exactly 7 images are required"})
			return
		}

		// Save each file and create image paths
		imagePaths := make([]string, 7)
		for i, file := range files {
			filename := filepath.Base(file.Filename)
			uploadPath := filepath.Join("uploads", filename)
			if err := c.SaveUploadedFile(file, uploadPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
				return
			}
			imagePaths[i] = "/" + strings.ReplaceAll(uploadPath, "\\", "/")
		}

		// Set the image paths in the Car struct
		car.ImagePath1 = imagePaths[0]
		car.ImagePath2 = imagePaths[1]
		car.ImagePath3 = imagePaths[2]
		car.ImagePath4 = imagePaths[3]
		car.ImagePath5 = imagePaths[4]
		car.ImagePath6 = imagePaths[5]
		car.ImagePath7 = imagePaths[6]

		// Save the car to the database
		if err := db.Create(&car).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the car to the database"})
			return
		}

		// Redirect to the listing page
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}
