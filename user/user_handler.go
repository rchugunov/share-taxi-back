package user

import (
	"com.github.rchugunov/share-taxi-back/entities"
	"com.github.rchugunov/share-taxi-back/gorm"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response class is used for returning result back to client.
type Response struct {
	entities.BaseResponse
	User entities.User `json:"user,omitempty"`
}

// GetUser handles request, checks user token and returns user object if exists.
func GetUser(c *gin.Context, userDao gorm.UserDao, tokenDao gorm.TokenDao) {
	var token string
	if token = c.GetHeader("token"); token == "" {
		c.JSON(http.StatusForbidden, Response{
			BaseResponse: entities.BaseResponse{Message: "Please send token in header"},
		})
		return
	}

	var userId *string
	if userId = tokenDao.GetUserIdIfValidToken(token); userId == nil {
		c.JSON(http.StatusForbidden, Response{
			BaseResponse: entities.BaseResponse{Message: "Could not find user. Try to relogin"},
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
