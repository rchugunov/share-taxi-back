package user

import (
	"com/github/rchugunov/share-taxi-back/entities"
	"com/github/rchugunov/share-taxi-back/gorm"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	entities.BaseResponse
	User entities.User `json:"user,omitempty"`
}

func GetUser(c *gin.Context, userDao gorm.UserDao, tokenDao gorm.TokenDao) {
	var token string
	if token = c.GetHeader("token"); token == "" {
		c.JSON(http.StatusForbidden, Response{
			BaseResponse: entities.BaseResponse{Message: "please send user_token in header"},
		})
		return
	}

	var userId *string
	if userId = tokenDao.GetUserIdIfValidToken(token); userId == nil {
		c.JSON(http.StatusForbidden, Response{
			BaseResponse: entities.BaseResponse{Message: "could not find user. try to relogin"},
		})

		return
	}

	user := userDao.GetUserById(*userId)
	c.JSON(http.StatusOK, Response{
		User: entities.User{
			Id:                   user.Id,
			Email:                user.Email,
			FirstName:            user.FirstName,
			LastName:             user.LastName,
			PhotoUrl:             user.PhotoUrl,
			HexBytesPhotoPreview: user.GetPhotoPreviewHex(),
		},
	})

}
