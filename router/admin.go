package routes

import (
	"ginapp/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/admin", handlers.Dashboard(db))
	r.GET("/myapp/get_banner_images", handlers.Get_Banner_Images(db))
	r.GET("/myapp/get_all_vehicles", handlers.GetAllVehicles(db))

	admin := r.Group("/admin")
	{
		admin.POST("/cars/add", handlers.AddCar(db))
		admin.GET("/cars/pdf_report", handlers.Get_Pdf_Report)
		admin.POST("/cars/edit/:id", handlers.EditCar(db))
		admin.POST("/cars/delete/:id", handlers.DeleteCar(db))
	}

	return r
}
