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
	OrigPoint:      entities.Location{Lat: 10, Long: 10},
	DestPoint:      entities.Location{Lat: 10, Long: 10},
	CreatedAt:      time.Now().UTC(),
	WaitingSeconds: 100,
}

var searchRequest2 = entities.SearchRequest{
	OrigPoint:      entities.Location{Lat: 12, Long: 11},
	DestPoint:      entities.Location{Lat: 13, Long: 9},
	CreatedAt:      time.Now().UTC().AddDate(0, 0, 1),
	WaitingSeconds: 120,
}

func TestSearchesDaoImpl_SearchUsersNearBy(t *testing.T) {

	var searchRequest1 = entities.SearchRequest{
		OrigPoint:      entities.Location{Lat: 10, Long: 10},
		DestPoint:      entities.Location{Lat: 10, Long: 10},
		CreatedAt:      time.Date(0, 0, 1, 0, 0, 30, 0, time.Local),
		WaitingSeconds: 120,
	}

	var searchRequest2 = entities.SearchRequest{
		OrigPoint:      entities.Location{Lat: 10, Long: 10},
		DestPoint:      entities.Location{Lat: 10, Long: 10},
		CreatedAt:      time.Date(0, 0, 1, 0, 2, 40, 0, time.Local),
		WaitingSeconds: 100,
	}

	testSearchUsersNearby(t, searchRequest1, searchRequest2, false, "Expected behavior: Second starts later then first finishes")

	var searchRequest3 = entities.SearchRequest{
		OrigPoint:      entities.Location{Lat: 10, Long: 10},
		DestPoint:      entities.Location{Lat: 10, Long: 10},
		CreatedAt:      time.Date(0, 0, 1, 0, 0, 30, 0, time.Local),
		WaitingSeconds: 120,
	}

	var searchRequest4 = entities.SearchRequest{
		OrigPoint:      entities.Location{Lat: 10, Long: 10},
		DestPoint:      entities.Location{Lat: 10, Long: 10},
		CreatedAt:      time.Date(0, 0, 1, 0, 2, 10, 0, time.Local),
		WaitingSeconds: 100,
	}

	testSearchUsersNearby(t, searchRequest3, searchRequest4, true, "Expected behavior: Second starts before then first finishes")

	var searchRequest5 = entities.SearchRequest{
		OrigPoint:      entities.Location{Lat: 10, Long: 10},
		DestPoint:      entities.Location{Lat: 10, Long: 10},
		CreatedAt:      time.Date(0, 0, 1, 0, 2, 30, 0, time.Local),
		WaitingSeconds: 120,
	}

	var searchRequest6 = entities.SearchRequest{
		OrigPoint:      entities.Location{Lat: 10, Long: 10},
		DestPoint:      entities.Location{Lat: 10, Long: 10},
		CreatedAt:      time.Date(0, 0, 1, 0, 2, 10, 0, time.Local),
		WaitingSeconds: 10,
	}

	testSearchUsersNearby(t, searchRequest5, searchRequest6, false, "Expected behavior: Second finishes before then first starts")

	var searchRequest7 = entities.SearchRequest{
		OrigPoint:      entities.Location{Lat: 10, Long: 10},
		DestPoint:      entities.Location{Lat: 10, Long: 10},
		CreatedAt:      time.Date(0, 0, 1, 0, 2, 30, 0, time.Local),
		WaitingSeconds: 120,
	}

	var searchRequest8 = entities.SearchRequest{
		OrigPoint:      entities.Location{Lat: 10, Long: 10},
		DestPoint:      entities.Location{Lat: 10, Long: 10},
		CreatedAt:      time.Date(0, 0, 1, 0, 2, 10, 0, time.Local),
		WaitingSeconds: 100,
	}

	testSearchUsersNearby(t, searchRequest7, searchRequest8, false, "Expected behavior: Second starts before then first starts, but finishes after first starts")
}

func testSearchUsersNearby(
	t *testing.T,
	searchRequest1 entities.SearchRequest,
	searchRequest2 entities.SearchRequest,
	searchShouldBeFound bool,
	message string,
) {
	searchesDao.Connect()
	defer searchesDao.Disconnect()
	userDao.Connect()
	defer userDao.Disconnect()

	userDao.Clear()
	testMyUser := User{Email: "testMyUser@gmail.co"}
	userDao.AddNewUser(&testMyUser)

	testOtherUser := User{Email: "testOtherUser@gmail.co"}
	userDao.AddNewUser(&testOtherUser)

	searchesDao.Clear()

	var search1 *Search
	var err error
	if search1, err = searchesDao.createNewSearch(testMyUser.Id, searchRequest1); err != nil {
		t.Error(err.Error())
	}

	assert.NotEmpty(t, search1.Id)

	var searchResult *entities.SearchResult
	if searchResult, err = searchesDao.SearchUsersNearby(testOtherUser.Id, searchRequest2); err != nil {
		t.Error(err.Error())
	}

	if searchShouldBeFound {
		assert.NotNil(t, searchResult, message)
		assert.Equal(t, testOtherUser.Id, searchResult.UserList[0].User.Id, message)
	} else {
		assert.Nil(t, searchResult, message)
	}
}

func TestSearchesDaoImpl_createNewSearch(t *testing.T) {
	searchesDao.Connect()
	defer searchesDao.Disconnect()
	userDao.Connect()
	defer userDao.Disconnect()

	userDao.Clear()
	testMyUser := User{Email: "testMyUser@gmail.co"}
	userDao.AddNewUser(&testMyUser)

	searchesDao.Clear()
	if _, err := searchesDao.createNewSearch(testMyUser.Id, searchRequest1); err != nil {
		t.Error(err.Error())
	}
}

func TestSearchesDaoImpl_createNewSearch_TheSameUserNewSearch(t *testing.T) {
	searchesDao.Connect()
	defer searchesDao.Disconnect()
	userDao.Connect()
	defer userDao.Disconnect()

	userDao.Clear()
	testMyUser := User{Email: "testMyUser@gmail.co"}
	userDao.AddNewUser(&testMyUser)

	searchesDao.Clear()
	var search1 *Search
	var err error
	if search1, err = searchesDao.createNewSearch(testMyUser.Id, searchRequest1); err != nil {
		t.Error(err.Error())
	}

	assert.NotEmpty(t, search1.Id)

	var search2 *Search
	if search2, err = searchesDao.createNewSearch(testMyUser.Id, searchRequest2); err != nil {
		t.Error(err.Error())
	}

	assert.NotEmpty(t, search2.Id)
	assert.Equal(t, search1.Id, search2.Id)
	assert.NotEqual(t, search1.WaitingTime, search2.WaitingTime)
}
