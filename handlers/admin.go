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

func GetAllCustomers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			images     []domain.CustomerImage
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
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "1")) // Default limit to 10 if not provided

		// Calculate offset
		offset = (page - 1) * limit

		// Fetch total count of entries
		if err := db.Model(&domain.CustomerImage{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count entries"})
			return
		}

		// Fetch images with pagination
		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&images).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
			return
		}

		// Generate pagination information

		c.JSON(http.StatusOK, gin.H{
			"images":     images,
			"totalCount": totalCount,
		})
	}
}

func DeleteCustomerImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var customerImage domain.CustomerImage

		// Find the image record in the database
		if err := db.First(&customerImage, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}

		// Optionally, remove the file from the filesystem
		if err := os.Remove(customerImage.Path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the image file"})
			return
		}

		// Delete the record from the database
		if err := db.Delete(&customerImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the image record"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/get_uploads_page")
	}
}
func EditCustomerImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var customerImage domain.CustomerImage

		if err := db.First(&customerImage, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image"})
			return
		}

		// Save the new file
		newImagePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
		if err := c.SaveUploadedFile(file, newImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the new image"})
			return
		}

		// Update the database
		customerImage.Path = newImagePath
		if err := db.Save(&customerImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the image"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/get_uploads_page")
	}
}

func Add_Customer_Form(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer domain.CustomerImage

		customerImage, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch": "image"})
			return
		}

		customerImagePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), customerImage.Filename))
		if err := c.SaveUploadedFile(customerImage, customerImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
			return
		}

		customer.Path = customerImagePath

		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the customer image"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/get_uploads_page")
	}
}

func UploadImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			images     []domain.CustomerImage
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
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))

		// Calculate offset
		offset = (page - 1) * limit

		// Fetch total count of entries
		if err := db.Model(&domain.CustomerImage{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count entries"})
			return
		}

		// Fetch links with pagination
		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&images).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
			return
		}

		// Generate pagination links
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)
		for i := range pages {
			pages[i] = i + 1
		}

		c.HTML(http.StatusOK, "customers.html", gin.H{
			"Images":     images,
			"TotalCount": totalCount,
			"Page":       page,
			"Limit":      limit,
			"TotalPages": totalPages,
			"Pages":      pages,
		})
	}
}
func EditCarPage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		var car domain.Car
		var brands []domain.Brand

		if err := db.Preload("Brand").Preload("Images").First(&car, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch each car"})
		}

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the cars"})
			return
		}

		if err := db.Model(&domain.Car{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total count"})
			return
		}
		carTypes := []string{
			domain.CarTypeSedan,
			domain.CarTypeHatchback,
			domain.CarTypeSuv,
			domain.CarTypeBike,
		}
		fuelTypes := []string{
			domain.FuelTypePetrol,
			domain.FuelTypeDiesel,
			domain.FuelTypeCNG,
			domain.FuelTypeElectric,
		}

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"Car":       car,
			"Brands":    brands,
			"CarTypes":  carTypes,
			"FuelTypes": fuelTypes,
		})
	}

}

func BrandDelete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var brand domain.Brand
		id := c.Param("id")

		// Check if the brand exists
		if err := db.First(&brand, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}

		// here am
		if err := db.Model(&domain.Car{}).Where("brand_id = ?", brand.ID).Update("brand_id", nil).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update associated cars"})
			return
		}

		// Now delete the brand itself
		if err := db.Delete(&brand).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete brand"})
			return
		}

		// Success message or redirect
		c.Redirect(http.StatusSeeOther, "/admin/get_brand_page") // Redirect to brands list page
	}
}

func BrandEdit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var brand domain.Brand
		id := c.Param("id")
		fmt.Println("here is the id", id)

		// Validate if brand with given ID exists
		if err := db.First(&brand, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}

		// Bind form data to update brand name
		if err := c.ShouldBind(&brand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		// Check if the new name already exists
		var existingBrand domain.Brand
		if err := db.Where("name = ?", brand.Name).First(&existingBrand).Error; err == nil && existingBrand.ID != brand.ID {
			// A brand with this name already exists
			c.JSON(http.StatusConflict, gin.H{"error": "Brand name already exists"})
			return
		}

		// Update brand name
		brand.Name = c.PostForm("brand_name")

		// Save updated brand to the database
		if err := db.Save(&brand).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update brand"})
			return
		}

		// Redirect or respond with success message
		c.Redirect(http.StatusSeeOther, "/admin/get_brand_page") // Redirect to brands list page
	}
}

func Get_Brand_Page(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			brands     []domain.Brand
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
		if err := db.Model(&domain.Brand{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count entries"})
			return
		}
		// Fetch links with pagination
		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
			return
		}

		// Generate pagination links
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)
		for i := range pages {
			pages[i] = i + 1
		}

		c.HTML(http.StatusOK, "brand.html", gin.H{
			"brands":     brands,
			"TotalCount": totalCount,
			"Page":       page,
			"Limit":      limit,
			"TotalPages": totalPages,
			"Pages":      pages,
		})
	}
}

func Get_Stock_Car_All(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			cars   []domain.Car
			count  int64
			page   int
			limit  int
			offset int
		)

		// Validate and set pagination parameters
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			limit = 2
		}

		offset = (page - 1) * limit

		brandIDStr := c.Query("brand_id") //  the query parameter is "brand_id" instead of "brand"
		fmt.Println("here is the brandid", brandIDStr)
		carType := c.Query("car_type")
		fuelType := c.Query("fuel_type")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")

		query := db.Model(&domain.Car{})

		if brandIDStr != "" {
			brandID, err := strconv.ParseUint(brandIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand_id format"})
				return
			}
			query = query.Where("brand_id = ?", brandID)
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

		if err := query.Order("created_at desc").Preload("Brand").Preload("Images").Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cars"})
			return
		}

		// Define a new structure to hold the filtered data
		type CarWithImage struct {
			ID           uint   `json:"id"`
			Brand        string `json:"brand"`
			Model        string `json:"model"`
			Status       string `json:"status"`
			Year         int    `json:"year"`
			Color        string `json:"color"`
			CarType      string `json:"car_type"`
			FuelType     string `json:"fuel_type"`
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
				Brand:        car.Brand.Name,
				Model:        car.Model,
				Year:         car.Year,
				Color:        car.Color,
				CarType:      car.CarType,
				FuelType:     car.FuelType,
				Variant:      car.Variant,
				Kms:          car.Kms,
				Ownership:    car.Ownership,
				Transmission: car.Transmission,
				Status:       car.Status,
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
			Brands    []domain.Brand `json:"brands"`     // getting the brand
			CarTypes  []string       `json:"car_types"`  // CarType
			FuelTypes []string       `json:"fuel_types"` //FuelType
		}

		// Fetch distinct brands
		var brands []domain.Brand
		if err := db.Find(&brands).Error; err != nil {
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
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit to 2 if not provided

		// Calculate offset
		offset = (page - 1) * limit

		// Fetch total count of entries
		// Fetch total count of entries
		if err := db.Model(&domain.YoutubeLink{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count entries"})
			return
		}
		// Fetch links with pagination
		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&links).Error; err != nil {
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
		var (
			links      []domain.YoutubeLink
			page       int
			limit      int
			offset     int
			totalCount int64
		)

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			limit = 10
		}

		offset = (page - 1) * limit

		// Fetch all YouTube links from the database
		if err := db.Model(&domain.YoutubeLink{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count"})
			return
		}

		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&links).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch YouTube links"})
			return
		}

		// Respond with JSON containing YouTube links
		c.JSON(http.StatusOK, gin.H{
			"links":      links,
			"totalCount": totalCount,
		})
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
			brands     []domain.Brand
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
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5")) // Default limit to 5 if not provided

		// Calculate offset
		offset = (page - 1) * limit

		//  THis one for brand  fetching
		if err := db.Preload("Brand").Order("created_at desc").Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cars"})
			return
		}

		// Fetch total count of cars (for pagination)
		if err := db.Model(&domain.Car{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total count"})
			return
		}

		//--> this one for   fetch the brand seperately
		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch brands"})
			return
		}

		// Fetch associated images for each car (if Image is a related entity)
		for i := range cars {
			if err := db.Model(&cars[i]).Association("Images").Find(&cars[i].Images); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
				return
			}
		}

		// Define car types and fuel types
		carTypes := []string{
			domain.CarTypeSedan,
			domain.CarTypeHatchback,
			domain.CarTypeSuv,
			domain.CarTypeBike,
			domain.CarTypeSport,
		}
		fuelTypes := []string{
			domain.FuelTypePetrol,
			domain.FuelTypeDiesel,
			domain.FuelTypeCNG,
			domain.FuelTypeElectric,
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
			"Brands":     brands,
			"Pages":      pages,
			"CarTypes":   carTypes,
			"FuelTypes":  fuelTypes,
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
				Brand:     car.Brand.Name,
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
		if err := db.Preload("Brand").Preload("Images").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create a new PDF
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()

		// Set header font
		pdf.SetFont("Arial", "B", 16)

		// Write header
		pdf.Cell(190, 10, "Cars Report")
		pdf.Ln(12)

		// Loop through cars and write data to PDF
		for _, car := range cars {
			pdf.SetFont("Arial", "B", 14)
			pdf.Cell(190, 10, fmt.Sprintf("Brand: %s, Model: %s", car.Brand.Name, car.Model))
			pdf.Ln(8)

			pdf.SetFont("Arial", "", 12)

			// Create a table layout
			pdf.CellFormat(50, 10, "Year:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, fmt.Sprintf("%d", car.Year), "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(50, 10, "Color:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, car.Color, "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(50, 10, "Variant:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, car.Variant, "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(50, 10, "Kms:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, fmt.Sprintf("%d", car.Kms), "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(50, 10, "Ownership:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, fmt.Sprintf("%d", car.Ownership), "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(50, 10, "Transmission:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, car.Transmission, "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(50, 10, "Reg No:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, car.RegNo, "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(50, 10, "Status:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, car.Status, "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(50, 10, "Price:", "1", 0, "L", false, 0, "")
			pdf.CellFormat(140, 10, fmt.Sprintf("%d", car.Price), "1", 0, "L", false, 0, "")
			pdf.Ln(-1)

			// Add spacing between cars
			pdf.Ln(12)
		}

		// Serve the PDF file
		c.Header("Content-Type", "application/pdf")
		err := pdf.Output(c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		}
	}
}

// Register the route
func Get_Banner_Vehicles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cars []domain.Car

		if err := db.Order("created_at desc").Limit(5).Preload("Brand").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tha database"})
			return

		}

		// Create a structure to hold the response data
		type CarDetail struct {
			BannerImage string `json:"bannerImage"`
			Brand       string `json:"brand"`
			Id          int    `json:"id"`
			Year        int    `json:"year"`
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
				Brand:       car.Brand.Name,
				Year:        int(car.Year),
				Id:          int(car.ID),
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

		if err := db.Order("created_at desc").Limit(6).Preload("Brand").Preload("Images").Find(&cars).Error; err != nil {
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
			CarType      string `json:"car_type"` //new one
			FuelType     string `json:"fuel_type"`
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
				Brand:        car.Brand.Name,
				Model:        car.Model,
				Year:         car.Year,
				Color:        car.Color,
				Variant:      car.Variant,
				Kms:          car.Kms,
				Ownership:    car.Ownership,
				FuelType:     car.FuelType,
				CarType:      car.CarType,
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
		var brand domain.Brand

		// Capture brand name from form
		brandIDStr := c.PostForm("brand")

		brandID, err := strconv.ParseUint(brandIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
			return
		}
		if err := db.First(&brand, brandID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Brand not found"})
			return
		}

		// Assign the brand ID to the car
		car.BrandID = uint(brandID) // Assuming car.BrandID is of type uint

		// Populate other car fields from form data
		car.Model = c.PostForm("model")
		year, err := strconv.Atoi(c.PostForm("year"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
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
		fmt.Println("here is the cartype ajdnkljasdnjklasndkjasndcklansdcklnaskclnasjkcdnaskljcdnkajscnklasdncj", car.CarType)
		car.FuelType = c.PostForm("fuel_type")
		car.Engine_size = c.PostForm("engine_size")
		car.Insurance_date = c.PostForm("insurance_date")
		car.Location = c.PostForm("location")

		// Handle banner image upload
		bannerImage, err := c.FormFile("bannerimage")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the banner image"})
			return
		}
		bannerImagePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), bannerImage.Filename))
		if err := c.SaveUploadedFile(bannerImage, bannerImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the banner image"})
			return
		}
		car.Bannerimage = "/" + strings.ReplaceAll(bannerImagePath, "\\", "/")

		// Handle multiple images upload
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the form"})
			return
		}
		files := form.File["images[]"]
		var images []domain.Image

		for _, file := range files {
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			uploadPath := filepath.Join("uploads", filename)
			if err := c.SaveUploadedFile(file, uploadPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save an image"})
				return
			}
			imagePath := "/" + strings.ReplaceAll(uploadPath, "\\", "/")
			images = append(images, domain.Image{Path: imagePath})
		}
		car.Images = images

		// Save the car to the database
		if err := db.Create(&car).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add the car"})
			return
		}

		// Optionally preload the brand when querying the car
		if err := db.Preload("Brand").First(&car, car.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load car with brand"})
			return
		}

		// Redirect to admin page after successful car addition
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}

func Add_Brand_Page(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var brands domain.Brand
		brands.Name = c.PostForm("brand")

		if err := db.Create(&brands).Error; err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/get_brand_page")

		}
		c.Redirect(http.StatusSeeOther, "/admin/get_brand_page")
	}
}
func EditCar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var car domain.Car
		var brands []domain.Brand

		// Fetch the existing car with preloaded images and brand
		if err := db.Preload("Images").Preload("Brand").First(&car, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
			return
		}

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch brands"})
			return
		}

		// Update car details from the form
		car.Model = c.PostForm("model")
		year, err := strconv.Atoi(c.PostForm("year"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
			return
		}
		car.Year = year
		car.CarType = c.PostForm("car_type")
		car.FuelType = c.PostForm("fuel_type")
		car.Engine_size = c.PostForm("engine_size")
		car.Insurance_date = c.PostForm("insurance_date")
		car.Location = c.PostForm("location")
		car.Color = c.PostForm("color")
		car.Variant = c.PostForm("variant")
		car.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		car.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		car.Transmission = c.PostForm("transmission")
		car.RegNo = c.PostForm("regno")
		car.Status = c.PostForm("status")
		car.Price, _ = strconv.Atoi(c.PostForm("price"))

		// Handle banner image upload if provided
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

		// Handle the new images update
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images[]"]
			var newImages []domain.Image

			for _, file := range files {
				filename := filepath.Base(fmt.Sprintf("%d_%d_%s", car.ID, time.Now().UnixNano(), file.Filename))
				uploadPath := filepath.Join("uploads", filename)
				if err := c.SaveUploadedFile(file, uploadPath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the image"})
					return
				}
				imagePath := "/" + strings.ReplaceAll(uploadPath, "\\", "/")
				newImages = append(newImages, domain.Image{Path: imagePath})
			}

			if len(newImages) > 0 {
				car.Images = append(car.Images, newImages...)
			}
		}

		// Handle image deletion if requested
		deleteImageIDs := c.PostFormArray("delete_images")
		if len(deleteImageIDs) > 0 {
			var remainingImages []domain.Image
			for _, img := range car.Images {
				shouldDelete := false
				for _, id := range deleteImageIDs {
					if strconv.Itoa(int(img.ID)) == id {
						shouldDelete = true
						break
					}
				}
				if !shouldDelete {
					remainingImages = append(remainingImages, img)
				} else {
					// Delete image from filesystem if needed
					if err := deleteFile(strings.TrimPrefix(img.Path, "/")); err != nil {
						fmt.Println("Failed to delete image:", err)
					}
					// Also, delete the image from the database
					db.Delete(&img)
				}
			}
			car.Images = remainingImages
		}

		// Handle image replacement
		for _, img := range car.Images {
			file, err := c.FormFile(fmt.Sprintf("replace_image_%d", img.ID))
			if err == nil {
				// Delete the old image from filesystem
				if err := deleteFile(strings.TrimPrefix(img.Path, "/")); err != nil {
					fmt.Println("Failed to delete image:", err)
				}
				// Save the new image
				uploadPath := filepath.Join("uploads", fmt.Sprintf("%d_%d_%s", car.ID, time.Now().UnixNano(), file.Filename))
				if err := c.SaveUploadedFile(file, uploadPath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the image"})
					return
				}
				img.Path = "/" + strings.ReplaceAll(uploadPath, "\\", "/")
				db.Save(&img) // Save the updated path to the database
			}
		}

		brandID, err := strconv.ParseUint(c.PostForm("brand"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
			return
		}
		car.BrandID = uint(brandID)

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
		if err := db.Preload("Brand").Preload("Images").First(&vehicle, vehicleID).Error; err != nil {
			fmt.Println("here is the &vechilce")
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		// Return the vehicle details as JSON
		c.JSON(http.StatusOK, gin.H{"vehicle": vehicle})
	}
}
