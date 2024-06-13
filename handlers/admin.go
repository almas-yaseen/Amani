package handlers

import (
	"fmt"
	"ginapp/domain"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

// func GetChoices(c *gin.Context) {

// 	c.JSON(http.StatusOk,gin.H{
// 		"car_types":[]string {
// 			domain.c
// 		}
// 	})

// }

func Get_Stock_Car_All(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		var cars []domain.Car

		if err := db.Order("id desc").Preload("Images").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the cars"})
			return
		}

		// Create a new structure to hold car with a single image
		type CarWithImage struct {
			ID           uint   `json:"id"`
			Brand        string `json:"brand"`
			Model        string `json:"model"`
			Year         int    `json:"year"`
			Color        string `json:"color"`
			Variant      string `json:"variant"`
			Kms          int    `json:"kms"`
			Ownership    int    `json:"ownership"`
			Transmission string `json:"transmission"`
			Price        int    `json:"price"`
			Image        string `json:"image"`
		}

		var result []CarWithImage

		// Populate the new structure
		for _, car := range cars {
			var image string
			if len(car.Images) > 0 {
				image = car.Images[0].Path // Select the first image path as the representative image
			}
			carWithImage := CarWithImage{
				ID:        car.ID,
				Brand:     car.Brand,
				Model:     car.Model,
				Year:      car.Year,
				Color:     car.Color,
				Variant:   car.Variant,
				Kms:       car.Kms,
				Ownership: car.Ownership,

				Transmission: car.Transmission,

				Price: car.Price,
				Image: image,
			}
			result = append(result, carWithImage)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "vehicles": result})
	}
}

func Get_Pdf_Report(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query all cars with their images
		var cars []domain.Car
		if err := db.Preload("Images").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create a new PDF
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()

		// Set font
		pdf.SetFont("Arial", "B", 12)

		// Write header
		pdf.Cell(190, 10, "Cars Report")
		pdf.Ln(12)

		// Write data to PDF
		for _, car := range cars {
			pdf.SetFont("Arial", "B", 10)
			pdf.Cell(190, 10, fmt.Sprintf("Brand: %s, Model: %s", car.Brand, car.Model))
			pdf.Ln(8)

			pdf.SetFont("Arial", "", 10)
			pdf.CellFormat(190, 10, fmt.Sprintf("Year: %s", car.Year), "", 0, "L", false, 0, "")
			pdf.Ln(8)
			pdf.CellFormat(190, 10, fmt.Sprintf("Color: %s", car.Color), "", 0, "L", false, 0, "")
			pdf.Ln(8)
			pdf.CellFormat(190, 10, fmt.Sprintf("Variant: %s", car.Variant), "", 0, "L", false, 0, "")
			pdf.Ln(8)
			pdf.CellFormat(190, 10, fmt.Sprintf("Kms: %d", car.Kms), "", 0, "L", false, 0, "")
			pdf.Ln(8)
			pdf.CellFormat(190, 10, fmt.Sprintf("Ownership: %d", car.Ownership), "", 0, "L", false, 0, "")
			pdf.Ln(8)
			pdf.CellFormat(190, 10, fmt.Sprintf("Transmission: %s", car.Transmission), "", 0, "L", false, 0, "")
			pdf.Ln(8)
			pdf.CellFormat(190, 10, fmt.Sprintf("Reg No: %s", car.RegNo), "", 0, "L", false, 0, "")
			pdf.Ln(8)
			pdf.CellFormat(190, 10, fmt.Sprintf("Status: %s", car.Status), "", 0, "L", false, 0, "")
			pdf.Ln(8)
			pdf.CellFormat(190, 10, fmt.Sprintf("Price: %d", car.Price), "", 0, "L", false, 0, "")
			pdf.Ln(8)

			pdf.Ln(8)

			pdf.SetFont("Arial", "B", 10)

			pdf.Ln(8)
			pdf.SetFont("Arial", "", 10)

			pdf.Ln(10)
		}

		// Serve the PDF file
		c.Header("Content-Type", "application/pdf")
		pdf.Output(c.Writer)
	}
}

// Register the route

func Dashboard(db *gorm.DB) gin.HandlerFunc {
	fmt.Println("here is the dashboard")
	return func(c *gin.Context) {
		var cars []domain.Car

		if err := db.Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the database"})
			return
		}

		// Fetch associated images for each car
		for i, car := range cars {
			fmt.Println("here is the i cand car", i, car)
			var images []domain.Image
			if err := db.Where("car_id = ?", car.ID).Find(&images).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
				return
			}
			cars[i].Images = images

		}

		// Now, fetch all images for all cars
		var allImages []domain.Image
		if err := db.Find(&allImages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all images"})
			return
		}

		// Pass both cars and images to the HTML template
		c.HTML(http.StatusOK, "admin.html", gin.H{"Cars": cars, "Images": allImages})
	}
}
func Get_Banner_Vehicles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		var cars []domain.Car

		if err := db.Order("id desc").Limit(5).Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tha database"})
			return

		}

		// Create a structure to hold the response data
		type CarDetail struct {
			BannerImage string `json:"bannerImage"`
			Model       string `json:"model"`
			Variant     string `json:"variant"`
			Price       int    `json:"price"`
			Color       string `json:"color"`
		}

		var carDetails []CarDetail

		for _, car := range cars {
			carDetail := CarDetail{
				BannerImage: car.Bannerimage,
				Model:       car.Model,
				Variant:     car.Variant,
				Price:       car.Price,
				Color:       car.Color,
			}
			carDetails = append(carDetails, carDetail)
			fmt.Println("Car details:", carDetail)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "carDetails": carDetails})
	}
}
func GetAllVehicles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		var cars []domain.Car

		if err := db.Order("id desc").Limit(6).Preload("Images").Find(&cars).Error; err != nil { //
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the cars"})
			return
		}

		// Create a new structure to hold car with a single image
		type CarWithImage struct {
			ID           uint   `json:"id"`
			Brand        string `json:"brand"`
			Model        string `json:"model"`
			Year         int    `json:"year"`
			Color        string `json:"color"`
			Variant      string `json:"variant"`
			Kms          int    `json:"kms"`
			Ownership    int    `json:"ownership"`
			Bannerimage  string `json:"bannerimage"`
			Transmission string `json:"transmission"`
			RegNo        string `json:"regno"`
			Status       string `json:"status"`
			Price        int    `json:"price"`
			Image        string `json:"image"`
		}

		var result []CarWithImage

		// Populate the new structure
		for _, car := range cars {
			var image string
			if len(car.Images) > 0 {
				image = car.Images[0].Path // Select the first image path as the representative image
			}
			carWithImage := CarWithImage{
				ID:           car.ID,
				Brand:        car.Brand,
				Model:        car.Model,
				Year:         car.Year,
				Color:        car.Color,
				Variant:      car.Variant,
				Kms:          car.Kms,
				Ownership:    car.Ownership,
				Bannerimage:  car.Bannerimage,
				Transmission: car.Transmission,
				RegNo:        car.RegNo,
				Status:       car.Status,
				Price:        car.Price,
				Image:        image,
			}
			result = append(result, carWithImage)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "vehicles": result})
	}
}

func AddCar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var car domain.Car
		fmt.Println("Starting to process the AddCar request")

		car.Brand = c.PostForm("brand")
		car.Model = c.PostForm("model")
		year, err := strconv.Atoi(c.PostForm("year"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid year"})
			return
		}
		car.Year = year
		car.Color = c.PostForm("color")
		car.Variant = c.PostForm("variant")
		car.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		car.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		car.Transmission = c.PostForm("transmission")
		car.RegNo = c.PostForm("regno")
		car.Status = c.PostForm("status")
		car.Price, _ = strconv.Atoi(c.PostForm("price"))

		form, err := c.MultipartForm() // allows files to be uploaded along with other form fields
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get the form"})
			return
		}

		bannerImage, err := c.FormFile("bannerimage")

		// Create the full path for the banner image
		bannerImagePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", car.ID, bannerImage.Filename))
		fmt.Println("here is the banner image path come on let asscd", bannerImagePath)
		if err := c.SaveUploadedFile(bannerImage, bannerImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
			return
		}
		// Replace backslashes with forward slashes
		bannerImagePath = "/" + strings.ReplaceAll(bannerImagePath, "\\", "/")
		car.Bannerimage = bannerImagePath

		files := form.File["images[]"]
		var images []domain.Image

		for _, file := range files {
			filename := filepath.Base(fmt.Sprintf("%d_%d_%s", car.ID, time.Now().UnixNano(), file.Filename)) // Using current time to ensure unique filename
			uploadPath := filepath.Join("uploads", filename)
			if err := c.SaveUploadedFile(file, uploadPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
				return
			}
			imagePath := "/" + strings.ReplaceAll(uploadPath, "\\", "/")
			images = append(images, domain.Image{Path: imagePath})
			fmt.Println("here is the second lastone", images)
		}
		car.Images = images

		fmt.Println("here is the updated one", car.Images)
		if err := db.Create(&car).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the car"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin")
	}
}
func EditCar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var car domain.Car

		// Fetch the existing car
		if err := db.Preload("Images").First(&car, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
			return
		}

		// Update car details from the form
		car.Brand = c.PostForm("brand")
		car.Model = c.PostForm("model")
		year, err := strconv.Atoi(c.PostForm("year"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
			return
		}
		car.Year = year
		car.Color = c.PostForm("color")
		car.Variant = c.PostForm("variant")
		car.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		car.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		car.Transmission = c.PostForm("transmission")
		car.RegNo = c.PostForm("regno")
		car.Status = c.PostForm("status")
		car.Price, _ = strconv.Atoi(c.PostForm("price"))

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the form"})
			return
		}

		bannerImage, err := c.FormFile("bannerimage")
		if err == nil {
			// Delete the old banner image if it exists
			if car.Bannerimage != "" {
				if err := deleteFile(strings.TrimPrefix(car.Bannerimage, "/")); err != nil {
					fmt.Println("Failed to delete the old banner image:", err)
				}
			}

			// Upload new banner image
			bannerImagePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", car.ID, bannerImage.Filename))
			if err := c.SaveUploadedFile(bannerImage, bannerImagePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the banner image"})
				return
			}
			bannerImagePath = "/" + strings.ReplaceAll(bannerImagePath, "\\", "/")
			car.Bannerimage = bannerImagePath
		}

		// Handle the images update
		files := form.File["images[]"]
		var images []domain.Image

		for _, file := range files {
			filename := filepath.Base(fmt.Sprintf("%d_%d_%s", car.ID, time.Now().UnixNano(), file.Filename))
			uploadPath := filepath.Join("uploads", filename)
			if err := c.SaveUploadedFile(file, uploadPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the image"})
				return
			}
			imagePath := "/" + strings.ReplaceAll(uploadPath, "\\", "/")
			images = append(images, domain.Image{Path: imagePath})
		}

		// Update the car's images if new images are uploaded
		if len(images) > 0 {
			// Delete existing images from the file system and the database
			for _, img := range car.Images {
				if err := deleteFile(strings.TrimPrefix(img.Path, "/")); err != nil {
					fmt.Println("Failed to delete old image:", err)
				}
			}

			if err := db.Where("car_id = ?", car.ID).Delete(&domain.Image{}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete existing images"})
				return
			}
			// Save new images
			car.Images = images
		}

		// Save the updated car details
		if err := db.Save(&car).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the car"})
			return
		}

		// Redirect to the admin page
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}

func DeleteCar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		carID, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid car ID")
			return
		}

		var car domain.Car
		if err := db.Preload("Images").First(&car, carID).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch the car details")
			return
		}

		// Delete the banner image
		if err := deleteFile(strings.TrimPrefix(car.Bannerimage, "/")); err != nil {
			fmt.Println("Failed to delete the banner image file:", err)
		}

		// Delete associated images
		for _, image := range car.Images {
			if err := deleteFile(strings.TrimPrefix(image.Path, "/")); err != nil {
				fmt.Println("Failed to delete image file:", err)
			}
		}

		// Delete the images from the database
		if err := db.Where("car_id = ?", carID).Delete(&domain.Image{}).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to delete car images")
			return
		}

		// Delete the car from the database
		if err := db.Where("id = ?", carID).Delete(&domain.Car{}).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to delete car")
			return
		}

		// Redirect to admin page
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}

func deleteFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path is empty")
	}

	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return err
	}

	return nil
}
func Get_Specific_Vehicle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the vehicle ID from the URL parameters
		id := c.Param("id")

		// Convert the ID string to an integer
		vehicleID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
			return
		}

		// Fetch the specific vehicle from the database
		var vehicle domain.Car
		if err := db.Preload("Images").First(&vehicle, vehicleID).Error; err != nil {
			fmt.Println("here is the &vechilce")
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		// Return the vehicle details as JSON
		c.JSON(http.StatusOK, gin.H{"vehicle": vehicle})
	}
}
