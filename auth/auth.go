package auth

import (
	fb "com.github.rchugunov/share-taxi-back/auth/facebook_api"
	"com.github.rchugunov/share-taxi-back/entities"
	"com.github.rchugunov/share-taxi-back/gorm"
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

type Response struct {
	entities.BaseResponse
	User        entities.User `json:"user,omitempty"`
	Token       string        `json:"token"`
	NewPassword string        `json:"new_password"`
}

func HandleFacebookLogin(c *gin.Context, userDao gorm.UserDao, tokenDao gorm.TokenDao, fbApi fb.FacebookApi) {
	var json FacebookAuthData

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusForbidden,
			Response{BaseResponse: entities.BaseResponse{Message: "Couldn't parse auth data", Exception: err.Error()}})
		return
		//logging.Timber(fmt.Sprintf("Couldn't bind request data %s", err.Error()))
	}

	if strings.TrimSpace(json.Token) == "" || strings.TrimSpace(json.UserID) == "" {
		c.JSON(http.StatusForbidden,
			Response{BaseResponse: entities.BaseResponse{Message: "token or user_id are wrong format"}})
		return
	}

	res, err := fbApi.FbGetProfile(json.Token, json.UserID)
	if ae, ok := err.(fb.FbGetError); ok {
		c.JSON(http.StatusForbidden,
			Response{BaseResponse: entities.BaseResponse{Message: ae.Message, Exception: ae.Cause.Error()}})
		return
	}

	var user *gorm.User
	var newPassword *string
	if user, newPassword, err = createNewUserIfNotExists(userDao, res); user == nil || err != nil {
		c.JSON(http.StatusForbidden,
			Response{BaseResponse: entities.BaseResponse{Message: "Couldn't validate user", Exception: err.Error()}})
		return
	}

	var token *string
	if token, err = generateNewToken(tokenDao, *user); err != nil {
		c.JSON(http.StatusForbidden,
			Response{
				User: entities.User{
					Id:                   user.Id,
					Email:                user.Email,
					FirstName:            user.FirstName,
					LastName:             user.LastName,
					PhotoUrl:             user.PhotoUrl,
					HexBytesPhotoPreview: user.GetPhotoPreviewHex(),
				},
				BaseResponse: entities.BaseResponse{Message: "Couldn't create token for user", Exception: err.Error()}},
		)
		return
	}

	if newPassword == nil {
		c.JSON(http.StatusOK,
			Response{
				User: entities.User{
					Id:                   user.Id,
					Email:                user.Email,
					FirstName:            user.FirstName,
					LastName:             user.LastName,
					PhotoUrl:             user.PhotoUrl,
					HexBytesPhotoPreview: user.GetPhotoPreviewHex(),
				},
				Token: *token,
			},
		)
		return
	}

	c.JSON(http.StatusOK,
		Response{
			User: entities.User{
				Id:                   user.Id,
				Email:                user.Email,
				FirstName:            user.FirstName,
				LastName:             user.LastName,
				PhotoUrl:             user.PhotoUrl,
				HexBytesPhotoPreview: user.GetPhotoPreviewHex(),
			},
			Token:       *token,
			NewPassword: *newPassword,
		})
}

func HandleLoginWithPassword(c *gin.Context, userDao gorm.UserDao, tokenDao gorm.TokenDao) {
	var json BasicAuthData
	var err error
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusForbidden,
			Response{BaseResponse: entities.BaseResponse{Message: "Couldn't parse auth data", Exception: err.Error()}},
		)
		return
	}
	//line
	if strings.TrimSpace(json.Login) == "" || strings.TrimSpace(json.Password) == "" {
		c.JSON(http.StatusForbidden,
			Response{BaseResponse: entities.BaseResponse{Message: "login or password are wrong format"}})
		return
	}

	var userDTO *gorm.User
	if userDTO, err = checkUserInDb(&userDao, json.Login, json.Password); err != nil || userDTO.Id == "" {
		c.JSON(http.StatusForbidden,
			Response{BaseResponse: entities.BaseResponse{Message: "Couldn't validate user", Exception: err.Error()}})
		return
	}

	var token *string
	if token, err = generateNewToken(tokenDao, *userDTO); err != nil {
		c.JSON(http.StatusForbidden,
			Response{
				User: entities.User{
					Id:                   userDTO.Id,
					Email:                userDTO.Email,
					FirstName:            userDTO.FirstName,
					LastName:             userDTO.LastName,
					PhotoUrl:             userDTO.PhotoUrl,
					HexBytesPhotoPreview: userDTO.GetPhotoPreviewHex(),
				},
				BaseResponse: entities.BaseResponse{Message: "Couldn't create token for user", Exception: err.Error()}})
		return
	}
	c.JSON(http.StatusOK,
		Response{
			User: entities.User{
				Id:                   userDTO.Id,
				Email:                userDTO.Email,
				FirstName:            userDTO.FirstName,
				LastName:             userDTO.LastName,
				PhotoUrl:             userDTO.PhotoUrl,
				HexBytesPhotoPreview: userDTO.GetPhotoPreviewHex(),
			},
			Token: *token,
		})
}

func generateNewToken(tokenDao gorm.TokenDao, user gorm.User) (token *string, err error) {
	token = tokenDao.CreateSession(user.Id)
	if token == nil {
		err = fmt.Errorf("couldnt create new token")
	}
	return
}

func createNewUserIfNotExists(userDao gorm.UserDao, fbUser entities.User) (user *gorm.User, newPassword *string, err error) {
	dbUser, getUserError := userDao.GetUserByEmail(fbUser.Email)
	if dbUser.Id != "" {
		user = dbUser
		return
	} else if getUserError != nil {
		err = fmt.Errorf(fmt.Sprintf("Could not retrieve user: %s", getUserError.Error()))
		return
	} else {
		genPassword, generateError := password.Generate(8, 3, 0, true, false)
		if generateError != nil {
			err = fmt.Errorf(fmt.Sprintf("Could not generate password: %s", generateError.Error()))
			return
		} else {
			hasher := sha1.New()
			hasher.Write([]byte(genPassword))
			sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

			u := gorm.User{
				Email:        fbUser.Email,
				PasswordHash: sha,
				FirstName:    fbUser.FirstName,
				LastName:     fbUser.LastName,
				PhotoUrl:     fbUser.PhotoUrl,
			}
			user = &u
			err = nil
			newPassword = &genPassword
			userDao.AddNewUser(user)
			return
		}
	}
}

func checkUserInDb(userDao *gorm.UserDao, email string, password string) (user *gorm.User, err error) {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	dbUser, getUserError := (*userDao).GetUserByEmailAndPassword(email, sha)
	if dbUser != nil {
		return dbUser, nil
	} else if getUserError != nil {
		err = fmt.Errorf(fmt.Sprintf("Could not retrieve user: %s", (*getUserError).Error()))
		return
	}
	panic("Shouldn't get here")
}
