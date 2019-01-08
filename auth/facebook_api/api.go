package facebook_api

import (
	"fmt"
	"github.com/huandu/facebook"
)

type FacebookApi interface {
	FbGetEmail(token string, userId string) (string, error)
}

type FacebookApiImpl struct {
}

func (api FacebookApiImpl) FbGetEmail(token string, userId string) (string, error) {
	res, facebookAuthErr := facebook.Get(userId+"/", facebook.Params{
		"access_token": token,
		"fields":       "email,name",
	})

	if facebookAuthErr != nil {
		return "", FbGetError{Message: "Facebook login failed", Cause: facebookAuthErr}
		//c.JSON(http.StatusForbidden, gin.H{
		//	"Message":   "Facebook login failed",
		//	"exception": facebookAuthErr.Error(),
		//})
	} else {
		email := res.Get("email")
		return email.(string), nil
		//c.JSON(http.StatusOK, gin.H{"email": email})
	}
}

type FbGetError struct {
	Message string
	Cause   error
}

func (error FbGetError) Error() string {
	return fmt.Sprintf("%s %s", error.Message, error.Cause.Error())
}
