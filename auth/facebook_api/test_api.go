package facebook_api

import "com/github/rchugunov/share-taxi-back/entities"

type FacebookApiTestImpl struct {
	MockFbGetEmail func(token string, userId string) (entities.User, error)
}

func (api FacebookApiTestImpl) FbGetProfile(token string, userId string) (entities.User, error) {
	return api.MockFbGetEmail(token, userId)
}
