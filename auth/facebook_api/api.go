package facebook_api

import (
	"com.github.rchugunov/share-taxi-back/entities"
	"fmt"
	"github.com/huandu/facebook"
	"strings"
)

func init() {
	facebook.Debug = facebook.DEBUG_ALL
}

type FacebookApi interface {
	FbGetProfile(token string, userId string) (entities.User, error)
}

type FacebookApiImpl struct {
}

func (api FacebookApiImpl) FbGetProfile(token string, userId string) (fbUser entities.User, err error) {
	res, errFacebookAuth := facebook.Get(userId+"/", facebook.Params{
		"access_token": token,
		"fields":       "email,name",
	})

	if errFacebookAuth != nil {
		err = FbGetError{Message: "Facebook login failed", Cause: errFacebookAuth}
		return
	}
	fbUser.Email = res.Get("email").(string)
	for index, str := range strings.Fields(res.Get("name").(string)) {
		switch index {
		case 0:
			fbUser.FirstName = str
		case 1:
			fbUser.LastName = str
		}
	}

	res, errGetPhoto := facebook.Get(userId+"/picture", facebook.Params{
		"access_token": token,
		"type":         "large",
		"redirect":     false,
	})

	if errGetPhoto != nil {
		err = FbGetError{Message: "Could not fetch facebook logo", Cause: errGetPhoto}
		return
	}

	fbUser.PhotoUrl = res.Get("data").(map[string]interface{})["url"].(string)
	return
}

type FbGetError struct {
	Message string
	Cause   error
}

func (error FbGetError) Error() string {
	return fmt.Sprintf("%s %s", error.Message, error.Cause.Error())
}
