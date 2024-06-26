package routes

import (
	"ginapp/handlers"
	"ginapp/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// func setCORSHeaders() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://www.amanimotors.in")
// 		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		c.Header("Access-Control-Allow-Credentials", "true")

// 		// Allow OPTIONS method for preflight requests
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusNoContent)
// 			return
// 		}

// 		c.Next()
// 	}
// }

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	// r.Use(setCORSHeaders())

	myapp := r.Group("/myapp")

	r.GET("/admin/login", handlers.AdminLogin) // gettin
	r.POST("/adminlogin", handlers.AdminLogin) // submitting the form
	myapp.GET("/get_filter_types", handlers.GetFilterTypes(db))
	myapp.GET("/get_stock_car_all", handlers.Get_Stock_Car_All(db))
	myapp.GET("/get_youtube_links", handlers.GetYoutubeLinks(db))
	myapp.GET("/get_banner_vehicles", handlers.Get_Banner_Vehicles(db))
	myapp.GET("/get_choices", handlers.GetChoices)
	myapp.GET("/get_all_vehicles_homepage", handlers.GetAllVehicles(db))
	myapp.GET("/get_specific_vehicle/:id", handlers.Get_Specific_Vehicle(db))
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.GET("/", handlers.Dashboard(db))
		admin.POST("/edit_brand/:id", handlers.BrandEdit(db))
		admin.POST("/delete_brand/:id", handlers.BrandDelete(db))
		admin.GET("/get_youtube_link_form", handlers.Show_Youtube_Page(db))
		admin.GET("/get_brand_page", handlers.Get_Brand_Page(db))
		admin.POST("/add_brand", handlers.Add_Brand_Page(db))
		admin.POST("/adding_youtube_form", handlers.Adding_Youtube_Link(db))
		admin.POST("/get_youtube_link_form_edit/:id", handlers.Youtube_page_edit(db))
		admin.POST("/get_youtube_link_form_delete/:id", handlers.Youtube_page_delete(db))
		admin.POST("/cars/add", handlers.AddCar(db))
		admin.GET("/cars/pdf_report", handlers.Get_Pdf_Report(db))
		admin.GET("/logout", handlers.Logout)
		admin.GET("/edit-car/:id", handlers.EditCarPage(db))
		admin.POST("/cars/edit/:id", handlers.EditCar(db))

		admin.POST("/cars/delete/:id", handlers.DeleteCar(db))
	}

	return r
}
