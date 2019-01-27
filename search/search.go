package search

import (
	"com.github.rchugunov/share-taxi-back/entities"
	"com.github.rchugunov/share-taxi-back/gorm"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Request struct {
	OriginDistanceMeters      float32           `json:"distance_origin" binding:"required"`
	DestinationDistanceMeters float32           `json:"distance_destination" binding:"required"`
	Origin                    entities.Location `json:"position_start" binding:"required"`
	Destination               entities.Location `json:"position_finish" binding:"required"`
	WaitingTimeInSeconds      uint32            `json:"wait_in_sec" binding:"required"`
	CreatedAt                 uint32            `json:"created_at" binding:"required"`
}

type Response struct {
	entities.BaseResponse
	Data *[]entities.SearchResult `json:"data"`
}

func NewSearch(c *gin.Context, tokenDao gorm.TokenDao, searchesDao gorm.SearchesDao) {
	request := Request{}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			BaseResponse: entities.BaseResponse{Message: "Could not parse request", Exception: err.Error()}, Data: nil,
		})
		return
	}

	var token string
	if token = c.GetHeader("user_token"); token == "" {
		c.JSON(http.StatusForbidden, Response{
			BaseResponse: entities.BaseResponse{Message: "please send user_token in header"}, Data: nil,
		})
		return
	}

	var userId *string
	if userId = tokenDao.GetUserIdIfValidToken(token); userId == nil {
		c.JSON(http.StatusForbidden, Response{
			BaseResponse: entities.BaseResponse{Message: "could not find user. try to relogin"}, Data: nil,
		})

		return
	}

	var data *[]entities.SearchResult
	if data = findOtherSearches(*userId, request, searchesDao); data == nil {
		c.JSON(http.StatusOK, Response{
			BaseResponse: entities.BaseResponse{Message: "could not find any users nearby"}, Data: &[]entities.SearchResult{},
		})

		return
	}

	c.JSON(http.StatusOK, Response{Data: data})
}

func findOtherSearches(userId string, request Request, searchesDao gorm.SearchesDao) *[]entities.SearchResult {
	return &[]entities.SearchResult{}
}
