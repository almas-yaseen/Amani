package routes

import (
	"ginapp/handlers"
	"ginapp/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	r.GET("/admin/login", handlers.AdminLogin) // getting the form
	r.POST("/adminlogin", handlers.AdminLogin) // submitting the form
	r.GET("myapp/get_youtube_links", handlers.GetYoutubeLinks(db))
	r.GET("/myapp/get_banner_vehicles", handlers.Get_Banner_Vehicles(db))
	r.GET("/myapp/get_choices", handlers.GetChoices)
	r.GET("/myapp/get_all_vehicles_homepage", handlers.GetAllVehicles(db))
	r.GET("/myapp/get_stockcar_all", handlers.Get_Stock_Car_All(db))
	r.GET("/myapp/get_specific_vehicle/:id", handlers.Get_Specific_Vehicle(db))

	admin := r.Group("/admin")

	admin.Use(middleware.AuthMiddleware())
	{
		admin.GET("/", handlers.Dashboard(db))
		admin.POST("/get_youtube_link_form", handlers.Youtube_link(db))
		admin.POST("/cars/add", handlers.AddCar(db))
		admin.GET("/cars/pdf_report", handlers.Get_Pdf_Report(db))
		admin.GET("/logout", handlers.Logout)
		admin.POST("/cars/edit/:id", handlers.EditCar(db))
		admin.POST("/cars/delete/:id", handlers.DeleteCar(db))
	}

	return r
}
