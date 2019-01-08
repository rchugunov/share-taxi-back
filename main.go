package main

import (
	"com/github/rchugunov/share-taxi-back/auth"
	"com/github/rchugunov/share-taxi-back/auth/facebook_api"
	"com/github/rchugunov/share-taxi-back/events"
	"com/github/rchugunov/share-taxi-back/gorm"
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
				userDao := gorm.UserDaoImpl{}
				auth.HandleFacebookLogin(c, userDao, facebook_api.FacebookApiImpl{})
				userDao.Disconnect()
			})

		api.POST("/login/basic",
			func(c *gin.Context) {
				userDao := gorm.UserDaoImpl{}
				auth.HandleLoginWithPassword(c, userDao)
				userDao.Disconnect()
			})

		api.GET("/events", func(c *gin.Context) {
			events.TestEvent(c)
		})
	}

	return router
}
