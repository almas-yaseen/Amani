package main

import (
	"fmt"
	"ginapp/config"
	"ginapp/database"
	routes "ginapp/router"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading the config file: %v", err)
	}

	db, err := database.ConnectDatabase(cfg)
	fmt.Println("db is here", db)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	log.Println("Database connection successful!")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://www.amanimotors.in"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://www.amanimotors.in"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	routes.AdminRoutes(router.Group(""), db)

	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("localhost error  %v", err)
	}
}
