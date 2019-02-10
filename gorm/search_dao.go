package gorm

import (
	"com.github.rchugunov/share-taxi-back/entities"
	"github.com/dewski/spatial"
	"time"
)

type Search struct {
	Id          string        `gorm:"primary_key;not null"`
	UserId      string        `gorm:"unique;not null"`
	OrigPoint   spatial.Point `gorm:"not null"`
	DestPoint   spatial.Point `gorm:"not null"`
	CreatedAt   time.Time     `gorm:"not null"`
	WaitingTime uint16        `gorm:"not null"`
}

type SearchesDao interface {
	Connect()
	Disconnect() error
	SearchUsersNearBy(currentUserId string, request entities.SearchRequest) *entities.SearchResult
}

type SearchesDaoImpl struct {
	Connection
	userDao UserDao
}

func (dao *SearchesDaoImpl) SearchUsersNearBy(currentUserId string, request entities.SearchRequest) *entities.SearchResult {
	var searches []Search
	db := dao.Where("user_id != ? ", currentUserId).Find(searches)
	errors := db.GetErrors()
	if len(errors) > 0 {
		return nil
	}

	result := entities.SearchResult{
		UserList: []entities.UserWithLocation{},
	}

	for _, search := range searches {

		user := dao.userDao.GetUserById(search.UserId)

		result.UserList = append(
			result.UserList,
			entities.UserWithLocation{
				Location: entities.Location{Lat: search.OrigPoint.Lat, Long: search.OrigPoint.Lng},
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
	return &result
}
