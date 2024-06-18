package routes

import (
	"ginapp/handlers"
	"ginapp/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

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
		admin.GET("/get_youtube_link_form", handlers.Show_Youtube_Page(db))
		admin.POST("/adding_youtube_form", handlers.Adding_Youtube_Link(db))
		admin.POST("/get_youtube_link_form_edit/:id", handlers.Youtube_page_edit(db))
		admin.POST("/get_youtube_link_form_delete/:id", handlers.Youtube_page_delete(db))
		admin.POST("/cars/add", handlers.AddCar(db))
		admin.GET("/cars/pdf_report", handlers.Get_Pdf_Report(db))
		admin.GET("/logout", handlers.Logout)
		admin.POST("/cars/edit/:id", handlers.EditCar(db))

		admin.POST("/cars/delete/:id", handlers.DeleteCar(db))
	}

	return r
}
