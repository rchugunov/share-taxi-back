package auth

import (
	fb "com/github/rchugunov/share-taxi-back/auth/facebook_api"
	"com/github/rchugunov/share-taxi-back/gorm"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type FacebookAuthData struct {
	Token  string `json:"token" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

func HandleFacebookLogin(c *gin.Context, fbApi fb.FacebookApi) {
	var json FacebookAuthData
	bindErr := c.BindJSON(&json)
	if bindErr == nil {
		if strings.TrimSpace(json.Token) != "" && strings.TrimSpace(json.UserID) != "" {

			res, err := fbApi.FbGetEmail(json.Token, json.UserID)

			if ae, ok := err.(fb.FbGetError); ok {
				c.JSON(http.StatusForbidden, gin.H{
					"message":   ae.Message,
					"exception": ae.Cause.Error(),
				})
			} else {
				user, exists, newPassword := validateUserInDb(res)
				if exists {
					c.JSON(http.StatusOK, gin.H{
						"user": user.MapToGin(),
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"user":         user.MapToGin(),
						"new_password": newPassword,
					})
				}

			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "token or user_id are wrong format",
			})
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"message":   "Couldn't parse auth data",
			"exception": bindErr.Error(),
		})
		//logging.Timber(fmt.Sprintf("Couldn't bind request data %s", err.Error()))
	}
}

func validateUserInDb(email string) (user gorm.User, exists bool, newPassword string) {
	return gorm.User{}, false, ""
}
