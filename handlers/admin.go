package handlers

import (
	"fmt"
	"ginapp/domain"
	"math"
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

func Get_Stock_Car_All(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cars []domain.Car

		var count int64
		fmt.Println("aksdnjkancdkjasndc")
		brand := c.Query("brand")
		carType := c.Query("car_type")
		fuelType := c.Query("fuel_type")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")

		query := db.Model(&domain.Car{})

		if brand != "" {
			fmt.Println("here is the query", brand)
			query = query.Where("brand = ?", brand)
		}
		if carType != "" {
			query = query.Where("car_type = ?", carType)
		}
		if fuelType != "" {
			query = query.Where("fuel_type = ?", fuelType)
		}

		if minPrice != "" {
			minPriceFloat, err := strconv.ParseFloat(minPrice, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_price format"})
				return
			}
			query = query.Where("price >= ?", minPriceFloat)
		}

		if maxPrice != "" {
			maxPriceFloat, err := strconv.ParseFloat(maxPrice, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_price format"})
				return
			}
			query = query.Where("price <= ?", maxPriceFloat)
		}
		if err := query.Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count cars"})
			return
		}

		if err := query.Preload("Images").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cars"})
			return
		}

		// Define a new structure to hold the filtered data
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

		// Populate the new structure with the filtered data
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
				Transmission: car.Transmission,
				Price:        car.Price,
				Image:        image,
			}
			result = append(result, carWithImage)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "vehicles": result, "total_count": count})
	}
}

// Youtube_link handles POST requests to add multiple YouTube links
func Adding_Youtube_Link(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse multipart form data
		err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max size

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Extract links from the form data
		links := c.Request.PostForm["links[]"]
		fmt.Println("Received links:", links)

		var youtubeLinks []domain.YoutubeLink

		// Iterate over each link
		for _, link := range links {
			// Create YoutubeLink object and append to slice
			youtubeLink := domain.YoutubeLink{
				VideoLink: link, // Assuming link is already a URL string
			}
			youtubeLinks = append(youtubeLinks, youtubeLink)
		}

		// Insert into the database
		if err := db.Create(&youtubeLinks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save links to database"})
			return
		}

		// Respond with success message
		c.Redirect(http.StatusSeeOther, "/admin/get_youtube_link_form")
	}
}
func GetFilterTypes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var filterTypes struct {
			Brands    []string `json:"brands"`
			CarTypes  []string `json:"car_types"`
			FuelTypes []string `json:"fuel_types"`
		}

		// Fetch distinct brands
		var brands []string
		if err := db.Model(&domain.Car{}).Distinct("brand").Pluck("brand", &brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch brands"})
			return
		}
		filterTypes.Brands = brands

		// Use predefined car types and fuel types
		filterTypes.CarTypes = []string{
			domain.CarTypeSedan,
			domain.CarTypeHatchback,
			domain.CarTypeSuv,
			domain.CarTypeBike,
		}
		filterTypes.FuelTypes = []string{
			domain.FuelTypePetrol,
			domain.FuelTypeCNG,
			domain.FuelTypeDiesel,
			domain.FuelTypeElectric,
		}

		c.JSON(http.StatusOK, filterTypes)
	}
}

func Youtube_page_delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		video_id := c.Param("id")
		fmt.Println("here is the id", video_id)

		var links domain.YoutubeLink

		result := db.First(&links, video_id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to do the link id incorrect"})
			return
		}

		if err := db.Where("id=?", video_id).Delete(&domain.YoutubeLink{}).Error; err != nil {
			c.String(http.StatusInternalServerError, "failed to delete the car")
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/get_youtube_link_form")

	}
}

func Youtube_page_edit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		videoid := c.Param("id")
		fmt.Println("here is the id", videoid)

		newVideolink := c.PostForm("editVideoLink")

		var link domain.YoutubeLink

		result := db.First(&link, videoid)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to do the link"})
			return

		}
		link.VideoLink = newVideolink
		db.Save(&link)
		c.Redirect(http.StatusFound, "/admin/get_youtube_link_form")

	}
}

func Show_Youtube_Page(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			links      []domain.YoutubeLink
			totalCount int64
			page       int
			limit      int
			offset     int
		)

		// Parse query parameters for pagination
		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5")) // Default limit to 2 if not provided

		// Calculate offset
		offset = (page - 1) * limit

		// Fetch total count of entries
		// Fetch total count of entries
		if err := db.Model(&domain.YoutubeLink{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count entries"})
			return
		}
		// Fetch links with pagination
		if err := db.Limit(limit).Offset(offset).Find(&links).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
			return
		}

		// Generate pagination links
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)
		for i := range pages {
			pages[i] = i + 1
		}

		c.HTML(http.StatusOK, "show.html", gin.H{
			"links":      links,
			"TotalCount": totalCount,
			"Page":       page,
			"Limit":      limit,
			"TotalPages": totalPages,
			"Pages":      pages,
		})
	}
}

func GetYoutubeLinks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var links []domain.YoutubeLink

		// Fetch all YouTube links from the database
		if err := db.Find(&links).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch YouTube links"})
			return
		}

		// Respond with JSON containing YouTube links
		c.JSON(http.StatusOK, links)
	}
}

// Youtube_link handles POST requests to add multiple YouTube links
func Youtube_link(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse multipart form data
		err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max size

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Extract links from the form data
		links := c.Request.PostForm["links[]"]
		fmt.Println("Received links:", links)

		var youtubeLinks []domain.YoutubeLink

		// Iterate over each link
		for _, link := range links {
			// Create YoutubeLink object and append to slice
			youtubeLink := domain.YoutubeLink{
				VideoLink: link, // Assuming link is already a URL string
			}
			youtubeLinks = append(youtubeLinks, youtubeLink)
		}

		// Insert into the database
		if err := db.Create(&youtubeLinks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save links to database"})
			return
		}

		// Respond with success message
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}
func Logout(c *gin.Context) {
	// Clear the authentication cookie
	c.SetCookie("authenticated", "", -1, "/", "", false, true)
	// Redirect to the login page
	c.Redirect(http.StatusFound, "/admin/login")
}

func AdminLogin(c *gin.Context) {
	if cookie, err := c.Cookie("authenticated"); err == nil && cookie == "true" {
		// User is authenticated, redirect to the dashboard
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		fmt.Println("here is the username and asdasd", username)
		password := c.PostForm("password")
		fmt.Println("here is the username", username)

		if username == "amani" && password == "amani123" {
			c.SetCookie("authenticated", "true", 36000, "/", "", false, true)
			c.Redirect(http.StatusFound, "/admin")
			return

		}
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)

}
func Dashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			cars       []domain.Car
			totalCount int64
			page       int
			limit      int
			offset     int
		)

		// Parse query parameters for pagination
		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5")) // Default limit to 2 if not provided

		// Calculate offset
		offset = (page - 1) * limit

		// Fetch cars with pagination
		if err := db.Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cars"})
			return
		}

		// Fetch total count of cars (for pagination)
		if err := db.Model(&domain.Car{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total count"})
			return
		}

		// Fetch associated images for each car (if Image is a related entity)
		for i := range cars {
			if err := db.Model(&cars[i]).Association("Images").Find(&cars[i].Images); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
				return
			}

			// Fetch the existing CarType and FuelType for the car
			var existingCar domain.Car
			if err := db.Where("id = ?", cars[i].ID).First(&existingCar).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing car details"})
				return
			}
			cars[i].CarType = existingCar.CarType
			cars[i].FuelType = existingCar.FuelType
		}

		// Generate pagination links
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)
		for i := range pages {
			pages[i] = i + 1
		}

		// Pass cars, pagination info, and other necessary data to the HTML template
		c.HTML(http.StatusOK, "admin.html", gin.H{
			"Cars":       cars,
			"TotalCount": totalCount,
			"Page":       page,
			"Limit":      limit,
			"TotalPages": totalPages,
			"Pages":      pages,
		})
	}
}

func GetChoices(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"car_types": []string{
			domain.CarTypeSedan,
			domain.CarTypeHatchback,
			domain.CarTypeSuv,
			domain.CarTypeBike,
		},
		"fuel_types": []string{
			domain.FuelTypePetrol,
			domain.FuelTypeDiesel,
			domain.FuelTypeCNG,
			domain.FuelTypeElectric,
		},
	})

}

func Get_Stock_Car_All_unit(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Set CORS headers
		// // Set CORS headers
		// allowedOrigins := " https://www.amanimotors.in"
		// c.Header("Access-Control-Allow-Origin", allowedOrigins)
		// c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		// c.Header("Access-Control-Allow-Credentials", "true")

		var cars []domain.Car
		var totalcount int64

		if err := db.Model(&domain.Car{}).Count(&totalcount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count "})
			return

		}

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

		c.JSON(http.StatusOK, gin.H{"status": "success", "vehicles": result, "totalcount": totalcount})
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
func Get_Banner_Vehicles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "https://www.amanimotors.in")

		// Fetch the latest 5 cars
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

		var cars []domain.Car

		if err := db.Order("id desc").Limit(6).Preload("Images").Find(&cars).Error; err != nil {
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
		car.CarType = c.PostForm("car_type")
		fmt.Println("here is the car type", car.CarType)
		car.FuelType = c.PostForm("fuel_type")
		car.Engine_size = c.PostForm("engine_size")       //new
		car.Insurance_date = c.PostForm("insurance_date") //new
		car.Location = c.PostForm("location")             //new
		fmt.Println("here is the fuel type", car.FuelType)
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
		car.CarType = c.PostForm("car_type")
		car.FuelType = c.PostForm("fuel_type")

		car.Engine_size = c.PostForm("engine_size")       //new
		car.Insurance_date = c.PostForm("insurance_date") //new
		car.Location = c.PostForm("location")
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
