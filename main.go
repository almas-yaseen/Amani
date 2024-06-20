package main

import (
	"fmt"
	"ginapp/config"
	"ginapp/database"
	routes "ginapp/router"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://www.amanimotors.in")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusOK)
// 			return
// 		}

// 		c.Next()
// 	}
// }

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "https://www.amanimotors.in" || origin == "http://localhost:5173" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

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
	router.Use(CORSMiddleware())

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	routes.AdminRoutes(router.Group(""), db)

	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("localhost error  %v", err)
	}
}
