package gorm

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserDaoImpl_AddNewUser(t *testing.T) {
	userDao := UserDaoImpl{}
	userDao.Connect()
	defer userDao.Disconnect()
	hasher := sha1.New()
	hasher.Write([]byte("my_password"))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	user := User{Email: "123@gmail.com", PasswordHash: sha, FirstName: "Roman", LastName: "Chugunov"}
	userDao.AddNewUser(&user)
	userId := user.Id
	userDao.DeleteUser(&user)
	existingUser, _ := userDao.GetUserByEmail("123@gmail.com")
	assert.NotEmpty(t, userId, "user was not added unfortunately")
	assert.Empty(t, existingUser, "couldn't delete user unfortunately")
}

func TestUserDaoImpl_GetUserByEmailAndPassword(t *testing.T) {
	userDao := UserDaoImpl{}
	userDao.Connect()
	defer userDao.Disconnect()

	user, err := userDao.GetUserByEmailAndPassword("222", "000")
	assert.Nil(t, user, "user should be nil")
	assert.NotNil(t, err, "err should not be nil")
}
