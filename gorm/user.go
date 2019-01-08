package gorm

import "github.com/gin-gonic/gin"

type USerDao interface {
}

type User struct {
	email        string
	id           string
	passwordHash string
}

func (user User) MapToGin() gin.H {
	return gin.H{
		"email": user.email,
		"id":    user.id,
	}
}

//func AddUser(user User) (res) {
//
//}
