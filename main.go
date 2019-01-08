package main

import (
	"com/github/rchugunov/share-taxi-back/auth"
	"com/github/rchugunov/share-taxi-back/auth/facebook_api"
	"com/github/rchugunov/share-taxi-back/events"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := SetupRouter()

	router.Run(":" + port)
}

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	api := router.Group("/api/v1")
	{
		api.POST("/login/fb",
			func(c *gin.Context) {
				auth.HandleFacebookLogin(c, facebook_api.FacebookApiImpl{})
			})

		api.GET("/events", func(c *gin.Context) {
			events.TestEvent(c)
		})
	}

	return router
}
