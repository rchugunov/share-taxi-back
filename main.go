package main

import (
	"com/github/rchugunov/share-taxi-back/auth"
	"com/github/rchugunov/share-taxi-back/auth/facebook_api"
	"com/github/rchugunov/share-taxi-back/gorm"
	"com/github/rchugunov/share-taxi-back/search"
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
				userDao.Connect()
				defer userDao.Disconnect()

				tokenDao := gorm.TokenDaoImpl{}
				tokenDao.Connect()
				defer tokenDao.Disconnect()

				auth.HandleFacebookLogin(c, &userDao, &tokenDao, facebook_api.FacebookApiImpl{})
			})

		api.POST("/login/basic",
			func(c *gin.Context) {
				userDao := gorm.UserDaoImpl{}
				userDao.Connect()
				defer userDao.Disconnect()

				tokenDao := gorm.TokenDaoImpl{}
				tokenDao.Connect()
				defer tokenDao.Disconnect()

				auth.HandleLoginWithPassword(c, &userDao, &tokenDao)
			})

		api.GET("/user/:id",
			func(c *gin.Context) {
				userDao := gorm.UserDaoImpl{}
				userDao.Connect()
				defer userDao.Disconnect()

				tokenDao := gorm.TokenDaoImpl{}
				tokenDao.Connect()
				defer tokenDao.Disconnect()

				c.JSON(200, gin.H{
					"result_ok": "super",
				})
			})

		api.POST("/search", func(c *gin.Context) {
			searchesDao := gorm.SearchesDaoImpl{}
			searchesDao.Connect()
			defer searchesDao.Disconnect()

			tokenDao := gorm.TokenDaoImpl{}
			tokenDao.Connect()
			defer tokenDao.Disconnect()

			search.NewSearch(c, &tokenDao, searchesDao)
		})
	}

	return router
}
