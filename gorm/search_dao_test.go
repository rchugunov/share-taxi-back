package gorm

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"com.github.rchugunov/share-taxi-back/entities"
)

var userDao = &UserDaoImpl{}
var conn = Connection{}
var searchesDao = SearchesDaoImpl{
	Connection: conn,
	userDao:    userDao,
}

var searchRequest1 = entities.SearchRequest{
	OrigPoint:   entities.Location{Lat: 10, Long: 10},
	DestPoint:   entities.Location{Lat: 10, Long: 10},
	CreatedAt:   time.Now().UTC(),
	WaitingTime: 120,
}

var searchRequest2 = entities.SearchRequest{
	OrigPoint:   entities.Location{Lat: 12, Long: 11},
	DestPoint:   entities.Location{Lat: 13, Long: 9},
	CreatedAt:   time.Now().UTC().AddDate(0, 0, 1),
	WaitingTime: 120,
}

var searchResponse = entities.SearchResult{
	UserList: []entities.UserWithLocation{},
}

func TestSearchesDaoImpl_SearchUsersNearBy(t *testing.T) {
	userDao.Connect()
	defer userDao.Disconnect()

	testMyUser := &User{Email: "testMyUser@gmail.co"}
	testMyUser, _ = userDao.GetUserByEmail("testMyUser@gmail.co")
	if testMyUser == nil || testMyUser.Id == "" {
		userDao.AddNewUser(testMyUser)
	}

	testOtherUser := &User{Email: "testOtherUser@gmail.co"}
	testOtherUser, _ = userDao.GetUserByEmail("testOtherUser@gmail.co")
	if testOtherUser == nil || testOtherUser.Id == "" {
		userDao.AddNewUser(testOtherUser)
	}

	searchesDao.Connect()
	defer searchesDao.Disconnect()
	defer searchesDao.Clear()

	err := searchesDao.createNewSearch(testMyUser.Id, searchRequest1)
	assert.NoError(t, err)

	result, err := searchesDao.SearchUsersNearby(testOtherUser.Id, searchRequest2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.UserList), "Expected only one search")
	assert.Equal(t, testMyUser.Id, result.UserList[0].User.Id, "Previous search not found")
}

func TestSearchesDaoImpl_createNewSearch(t *testing.T) {
	searchesDao.Connect()
	defer searchesDao.Disconnect()
	userDao.Connect()
	defer userDao.Disconnect()

	testMyUser := &User{}
	testMyUser, _ = userDao.GetUserByEmail("testMyUser@gmail.co")
	if testMyUser == nil || testMyUser.Id == "" {
		testMyUser = &User{Email: "testMyUser@gmail.co"}
		userDao.AddNewUser(testMyUser)
	}

	if err := searchesDao.createNewSearch(testMyUser.Id, searchRequest1); err != nil {
		t.Error(err.Error())
	}
}
