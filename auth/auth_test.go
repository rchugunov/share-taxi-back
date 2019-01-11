package auth

import (
	fb "com/github/rchugunov/share-taxi-back/auth/facebook_api"
	"com/github/rchugunov/share-taxi-back/gorm"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const TEST_EMAIL string = "myemail@gmail.com"

func TestHandleFacebookLogin(t *testing.T) {
	body := gin.H{
		"token":   "123",
		"user_id": "123",
	}

	userDao := gorm.UserDaoImpl{}
	tokenDao := gorm.TokenDaoImpl{}
	// Grab our router
	router := setupRouter(fb.FacebookApiTestImpl{MockFbGetEmail: func(token string, userId string) (s string, e error) {
		return TEST_EMAIL, nil
	}}, &userDao, &tokenDao)
	userDao.Connect()
	tokenDao.Connect()
	defer userDao.Disconnect()
	defer tokenDao.Disconnect()

	strReq, _ := jsoniter.MarshalToString(body)

	// Perform a GET request with that handler.
	w := performRequestWithBody(router, "POST", "/api/v1/login/fb", strings.NewReader(strReq))
	// Delete user which we not gonna use anymore
	userDao.DeleteUserByEmail(TEST_EMAIL)
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response Response
	responseBytes := w.Body.Bytes()
	err := json.Unmarshal(responseBytes, &response)
	assert.NoError(t, err, string(responseBytes))
	assert.NotEmpty(t, response.User)
	assert.NotEmpty(t, *response.Token)
}

func TestHandleLoginWithPassword(t *testing.T) {

	userDao := gorm.UserDaoImpl{}
	tokenDao := gorm.TokenDaoImpl{}

	router := setupRouter(nil, &userDao, &tokenDao)
	userDao.Connect()
	tokenDao.Connect()
	defer userDao.Disconnect()
	defer tokenDao.Disconnect()

	// Making request
	request := BasicAuthData{Login: TEST_EMAIL, Password: "slfsdkfmksdflk"}
	strReq, _ := jsoniter.MarshalToString(request)

	user, err := userDao.GetUserByEmail(TEST_EMAIL)
	if user.Id == "" {
		userDao.AddNewUser(&gorm.User{Email: TEST_EMAIL, PasswordHash: "jDuA0aQQUKIy6h9P6eNP2Ez-mEc="})
	}
	w := performRequestWithBody(router, "POST", "/api/v1/login/basic", strings.NewReader(strReq))
	// Delete user which we not gonna use anymore
	userDao.DeleteUserByEmail(TEST_EMAIL)

	// Checking response
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response Response
	responseBytes := w.Body.Bytes()
	err = json.Unmarshal(responseBytes, &response)
	assert.NoError(t, err, string(responseBytes))
	assert.NotEmpty(t, response.User)
	assert.NotEmpty(t, *response.Token)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	return performRequestWithBody(r, method, path, nil)
}

func performRequestWithBody(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func setupRouter(fbApi fb.FacebookApi, userDao gorm.UserDao, tokenDao gorm.TokenDao) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	api := router.Group("/api/v1")
	{
		api.POST("/login/fb", func(c *gin.Context) {
			HandleFacebookLogin(c, userDao, tokenDao, fbApi)
		})

		api.POST("/login/basic", func(c *gin.Context) {
			HandleLoginWithPassword(c, userDao, tokenDao)
		})
	}

	return router
}
