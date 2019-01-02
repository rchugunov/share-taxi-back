package main

import (
	"com/github/rchugunov/share-taxi-back/events"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	api := router.Group("/api/v1")
	{
		api.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message":    "pong",
				"extra_data": map[string]string{"fsd": "fsdfsf"},
			})
		})

		api.GET("/events", func(c *gin.Context) {
			events.TestEvent(c)
		})
	}

	router.Run(":" + port)
}
