package facebook_api

type FacebookApiTestImpl struct {
	MockFbGetEmail func(token string, userId string) (string, error)
}

func (api FacebookApiTestImpl) FbGetEmail(token string, userId string) (string, error) {
	return api.MockFbGetEmail(token, userId)
}
