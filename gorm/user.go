package gorm

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserDao interface {
	Connect()
	Disconnect()
	GetUserByEmail(email string) (user *User, err error)
	GetUserByEmailAndPassword(email string, passwordHash string) (user *User, err *error)
	AddNewUser(user *User)
	DeleteUser(user *User)
	DeleteUserByEmail(email string)
}

type UserDaoImpl struct {
	Connection
}

func (dao UserDaoImpl) AddNewUser(user *User) {
	dao.dbInst.Create(&user)
}

func (dao UserDaoImpl) DeleteUser(user *User) {
	dao.dbInst.Delete(&user)
}

func (dao UserDaoImpl) DeleteUserByEmail(email string) {
	dao.dbInst.Where("email = ?", email).Delete(User{})
}

func (dao UserDaoImpl) GetUserByEmail(email string) (user *User, err error) {
	user = &User{}
	dao.dbInst.Where("email = ?", email).First(user)
	return
}

func (dao UserDaoImpl) GetUserByEmailAndPassword(email string, passwordHash string) (user *User, err *error) {
	user = &User{}
	dao.dbInst.Where("email = ? AND password_hash = ?", email, passwordHash).First(user)
	return
}

type User struct {
	Id           string `gorm:"primary_key;not null"`
	Email        string `gorm:"unique;not null"` // set member number to unique and not null
	PasswordHash string `gorm:"not null"`
	FirstName    string
	LastName     string `gorm:"column:second_name"`
}

func (user User) MapToGin() gin.H {
	return gin.H{
		"email": user.Email,
		"id":    user.Id,
	}
}
