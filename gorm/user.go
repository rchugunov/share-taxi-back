package gorm

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserDao interface {
	Connect()
	Disconnect() error
	GetUserByEmail(email string) (user *User, err error)
	GetUserByEmailAndPassword(email string, passwordHash string) (user *User, err *error)
	AddNewUser(user *User)
	DeleteUser(user *User)
	DeleteUserByEmail(email string)
	GetUserById(id string) (user *User)
}

type UserDaoImpl struct {
	Connection
}

func (dao UserDaoImpl) GetUserById(id string) (user *User) {
	user = &User{}
	dao.Where("id = ?", id).First(user)
	return
}

func (dao UserDaoImpl) AddNewUser(user *User) {
	dao.Create(&user)
}

func (dao UserDaoImpl) DeleteUser(user *User) {
	dao.Delete(&user)
}

func (dao UserDaoImpl) DeleteUserByEmail(email string) {
	dao.Where("email = ?", email).Delete(User{})
}

func (dao UserDaoImpl) GetUserByEmail(email string) (user *User, err error) {
	user = &User{}
	dao.Where("email = ?", email).First(user)
	return
}

func (dao UserDaoImpl) GetUserByEmailAndPassword(email string, passwordHash string) (user *User, err *error) {
	user = &User{}
	dao.Where("email = ? AND password_hash = ?", email, passwordHash).First(user)
	return
}

// User represents DTO model from user table in database.
type User struct {
	Id                string `gorm:"primary_key;not null"`
	Email             string `gorm:"unique;not null"` // set member number to unique and not null
	PasswordHash      string `gorm:"not null"`
	FirstName         string
	LastName          string `gorm:"column:second_name"`
	PhotoPreviewBytes []byte
	PhotoUrl          string
}

func (user User) GetPhotoPreviewHex() string {
	return ""
}
