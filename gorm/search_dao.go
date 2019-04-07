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
	FinishedAt  time.Time     `gorm:"not null"`
	WaitingTime uint16        `gorm:"not null;column:wait_seconds"`
}

type SearchesDao interface {
	Connect()
	Disconnect()
	SearchUsersNearby(currentUserId string, request entities.SearchRequest) (*entities.SearchResult, error)
	Clear()
}

type SearchesDaoImpl struct {
	Connection
	userDao UserDao
}

func (dao *SearchesDaoImpl) SearchUsersNearby(currentUserId string, request entities.SearchRequest) (*entities.SearchResult, error) {
	var searches []Search
	dao.Where(
		"user_id != ? and finished_at > ? and created_at < ?",
		currentUserId,
		request.CreatedAt,
		request.CreatedAt,
	).Find(&searches)
	dbErrors := dao.GetErrors()
	if len(dbErrors) > 0 {
		return nil, errors2.Wrap(dbErrors[0], "SearchUsersNearby failed")
	}

	if _, err := dao.createNewSearch(currentUserId, request); err != nil {
		return nil, err
	}

	if len(searches) == 0 {
		return nil, nil
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

func (dao *SearchesDaoImpl) createNewSearch(currentUserId string, request entities.SearchRequest) (*Search, error) {
	timeStarted := request.CreatedAt
	searchToUpdate := &Search{
		UserId:      currentUserId,
		OrigPoint:   spatial.Point{Lat: request.OrigPoint.Lat, Lng: request.OrigPoint.Long},
		DestPoint:   spatial.Point{Lat: request.DestPoint.Lat, Lng: request.DestPoint.Long},
		CreatedAt:   request.CreatedAt,
		WaitingTime: request.WaitingSeconds,
		FinishedAt:  timeStarted.Add(time.Duration(request.WaitingSeconds) * time.Second),
	}
	search := &Search{UserId: currentUserId}
	dao.Where(search).Assign(searchToUpdate).FirstOrCreate(search)
	dbErrors := dao.GetErrors()
	if len(dbErrors) > 0 {
		return nil, errors2.Wrap(dbErrors[0], "SearchUsersNearby failed")
	}

	if len(search.Id) == 0 {
		return nil, fmt.Errorf("couldn't create or update search")
	}

	return search, nil
}

func (dao *SearchesDaoImpl) Clear() {
	dao.Delete(Search{})
}
