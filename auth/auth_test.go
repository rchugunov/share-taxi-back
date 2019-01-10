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

func TestHandleFacebookLogin(t *testing.T) {
	body := gin.H{
		"token":   "123",
		"user_id": "123",
	}

	// Grab our router
	router := setupRouter(fb.FacebookApiTestImpl{MockFbGetEmail: func(token string, userId string) (s string, e error) {
		return "myemail@gmail.com", nil
	}})

	strReq, _ := jsoniter.MarshalToString(body)

	// Perform a GET request with that handler.
	w := performRequestWithBody(router, "POST", "/api/v1/login/fb", strings.NewReader(strReq))
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response map[string]string
	responseStr := w.Body.String()
	err := json.Unmarshal([]byte(responseStr), &response)
	assert.Error(t, err, responseStr)
	// Grab the value & whether or not it exists
	//value, exists := response["hello"]
	//// Make some assertions on the correctness of the response.
	//assert.Nil(t, err)
	//assert.True(t, exists)
	//assert.Equal(t, body["hello"], value)
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

func setupRouter(fbApi fb.FacebookApi) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	api := router.Group("/api/v1")
	{
		api.POST("/login/fb", func(c *gin.Context) {
			userDao := gorm.UserDaoImpl{}
			userDao.Connect()
			defer userDao.Disconnect()
			HandleFacebookLogin(c, &userDao, fbApi)
		})
	}

	return router
}
