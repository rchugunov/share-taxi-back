package gorm

import (
	"time"
)

type Token struct {
	UserId    string
	Value     string    `gorm:"primary_key;not null;column:token"`
	ExpiresAt time.Time `gorm:"column:expire_ts"`
}

type TokenDao interface {
	Connect()
	Disconnect()
	GetUserIdIfValidToken(token string) (userId *string)
	CreateSession(userId string) (newToken *string)
}

type TokenDaoImpl struct {
	Connection
}

func (dao TokenDaoImpl) CreateSession(userId string) (newToken *string) {
	token := Token{UserId: userId, ExpiresAt: time.Now().AddDate(0, 3, 0)}
	dao.dbInst.Create(&token)
	newToken = &token.Value
	return
}

func (dao TokenDaoImpl) GetUserIdIfValidToken(tokenValue string) (userId *string) {
	newToken := Token{}
	dao.dbInst.Where("token = ?", tokenValue).First(&newToken)
	userId = &newToken.UserId
	return
}
