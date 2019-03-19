package gorm

import (
	"com.github.rchugunov/share-taxi-back/entities"
	"fmt"
	"github.com/dewski/spatial"
	errors2 "github.com/pkg/errors"
	"time"
)

type Search struct {
	Id          string        `gorm:"primary_key;not null"`
	UserId      string        `gorm:"unique;not null"`
	OrigPoint   spatial.Point `gorm:"not null"`
	DestPoint   spatial.Point `gorm:"not null"`
	CreatedAt   time.Time     `gorm:"not null"`
	WaitingTime uint16        `gorm:"not null;column:wait_seconds"`
}

type SearchesDao interface {
	Connect()
	Disconnect() error
	SearchUsersNearby(currentUserId string, request entities.SearchRequest) (*entities.SearchResult, error)
	Clear()
}

type SearchesDaoImpl struct {
	Connection
	userDao UserDao
}

func (dao *SearchesDaoImpl) SearchUsersNearby(currentUserId string, request entities.SearchRequest) (*entities.SearchResult, error) {
	var searches []Search
	dao.Where("user_id != ? ", currentUserId).Find(&searches)
	dbErrors := dao.GetErrors()
	if len(dbErrors) > 0 {
		return nil, errors2.Wrap(dbErrors[0], "SearchUsersNearby failed")
	}

	if err := dao.createNewSearch(currentUserId, request); err != nil {
		return nil, err
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
	return &result, nil
}

func (dao *SearchesDaoImpl) createNewSearch(currentUserId string, request entities.SearchRequest) error {
	search := &Search{
		UserId:      currentUserId,
		OrigPoint:   spatial.Point{Lat: request.OrigPoint.Lat, Lng: request.OrigPoint.Long},
		DestPoint:   spatial.Point{Lat: request.DestPoint.Lat, Lng: request.DestPoint.Long},
		CreatedAt:   time.Now(),
		WaitingTime: request.WaitingTime,
	}
	dao.Where(Search{UserId: currentUserId}).Assign(search).FirstOrCreate(search)
	if len(search.Id) == 0 {
		return fmt.Errorf("couldn't create or update search")
	}
	dbErrors := dao.GetErrors()
	if len(dbErrors) > 0 {
		return errors2.Wrap(dbErrors[0], "SearchUsersNearby failed")
	}
	return nil
}

func (dao *SearchesDaoImpl) Clear() {
	dao.Delete(Search{})
}
