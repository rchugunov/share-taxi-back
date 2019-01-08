package auth

import (
	fb "com/github/rchugunov/share-taxi-back/auth/facebook_api"
	"com/github/rchugunov/share-taxi-back/gorm"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-password/password"
	"net/http"
	"strings"
)

type FacebookAuthData struct {
	Token  string `json:"token" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

type BasicAuthData struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func HandleFacebookLogin(c *gin.Context, userDao gorm.UserDao, fbApi fb.FacebookApi) {
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
				user, newPassword, err := createNewUserIfNotExists(&userDao, res)
				if user != nil {
					if newPassword == nil {
						c.JSON(http.StatusOK, gin.H{
							"user": user.MapToGin(),
						})
					} else {
						c.JSON(http.StatusOK, gin.H{
							"user":         user.MapToGin(),
							"new_password": newPassword,
						})
					}
				} else {
					c.JSON(http.StatusForbidden, gin.H{
						"message":   "Couldn't validate user",
						"exception": (*err).Error(),
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

func HandleLoginWithPassword(c *gin.Context, userDao gorm.UserDao) {
	var json BasicAuthData
	bindErr := c.BindJSON(&json)
	if bindErr == nil {
		if strings.TrimSpace(json.Login) != "" && strings.TrimSpace(json.Password) != "" {
			user, err := checkUserInDb(&userDao, json.Login, json.Password)

			if user != nil {
				c.JSON(http.StatusOK, gin.H{
					"user": user.MapToGin(),
				})
			} else {
				c.JSON(http.StatusForbidden, gin.H{
					"message":   "Couldn't validate user",
					"exception": (*err).Error(),
				})
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "login or password are wrong format",
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

func createNewUserIfNotExists(userDao *gorm.UserDao, email string) (user *gorm.User, newPassword *string, err *error) {
	dbUser, getUserError := (*userDao).GetUserByEmail(email)
	if dbUser != nil {
		return user, nil, nil
	} else if getUserError != nil {
		resError := fmt.Errorf(fmt.Sprintf("Could not retrieve user: %s", &getUserError))
		return nil, nil, &resError
	} else {
		genPassword, generateError := password.Generate(8, 3, 0, true, false)
		if generateError != nil {
			resError := fmt.Errorf(fmt.Sprintf("Could not generate password: %s", &generateError))
			return nil, nil, &resError
		} else {
			hasher := sha1.New()
			hasher.Write([]byte(genPassword))
			sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

			newUser := gorm.User{Email: email, PasswordHash: sha}

			return &newUser, &genPassword, nil
		}
	}
}

func checkUserInDb(userDao *gorm.UserDao, email string, password string) (user *gorm.User, err *error) {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	dbUser, getUserError := (*userDao).GetUserByEmailAndPassword(email, sha)
	if dbUser != nil {
		return dbUser, nil
	} else if getUserError != nil {
		resError := fmt.Errorf(fmt.Sprintf("Could not retrieve user: %s", &getUserError))
		return nil, &resError
	}
	panic("Shouldn't get here")
}
