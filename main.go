package main

import (
	"fmt"
	"ginapp/config"
	"ginapp/database"
	routes "ginapp/router"
	"log"

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

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	routes.AdminRoutes(router.Group(""), db)

	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("localhost error  %v", err)
	}
}
