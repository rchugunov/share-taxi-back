package gorm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenDaoImpl_CreateSession(t *testing.T) {

	userDao := UserDaoImpl{}
	userDao.Connect()
	defer userDao.Disconnect()

	tokenDao := TokenDaoImpl{}
	tokenDao.Connect()
	defer tokenDao.Disconnect()

	user := User{Email: "sjdfksdf@mksdfg.com", PasswordHash: "ksdfksldf"}
	userDao.AddNewUser(&user)

	if user.Id == "" {
		existingUser, _ := userDao.GetUserByEmail("sjdfksdf@mksdfg.com")
		user = *existingUser
	}

	token := tokenDao.CreateSession(user.Id)
	assert.NotEmpty(t, token)

	userId := tokenDao.GetUserIdIfValidToken(*token)
	assert.NotEmpty(t, userId)

	userDao.DeleteUser(&user)

	userId = tokenDao.GetUserIdIfValidToken(*token)
	assert.Empty(t, userId)

}
