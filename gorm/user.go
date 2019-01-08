package gorm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
)

type UserDao interface {
	Connect()
	Disconnect()
	GetUserByEmail(email string) (user *User, err error)
	GetUserByEmailAndPassword(email string, passwordHash string) (user *User, err *error)
}

type UserDaoImpl struct {
	dbInst *gorm.DB
}

func (dao UserDaoImpl) Connect() {
	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_HOST"),
			os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_PORT"),
			os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_USER"),
			os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_DBNAME"),
			os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_PASSWORD")))

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to db: %s", err.Error()))
	}
	db.SingularTable(true)
	dao.dbInst = db
}

func (dao UserDaoImpl) GetUserByEmail(email string) (user *User, err error) {
	dao.dbInst.Where("email = ?", email).First(user)
	return user, nil
}

func (dao UserDaoImpl) GetUserByEmailAndPassword(email string, passwordHash string) (user *User, err *error) {
	dao.dbInst.Where("email = ? AND password_hash = ?", email, passwordHash).First(user)
	return user, nil
}

func (dao UserDaoImpl) Disconnect() {
	dao.dbInst.Close()
}

type User struct {
	Id           string `gorm:"primary_key;not null"`
	Email        string `gorm:"unique;not null"` // set member number to unique and not null
	PasswordHash string `gorm:"not null"`
	FirstName    string
	LastName     string
}

func (user User) MapToGin() gin.H {
	return gin.H{
		"email": user.Email,
		"id":    user.Id,
	}
}

//func AddUser(user User) (res) {
//
//}
