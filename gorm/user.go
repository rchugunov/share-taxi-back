package gorm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

type UserDao interface {
	Connect()
	Disconnect()
	GetUserByEmail(email string) (user *User, err error)
	GetUserByEmailAndPassword(email string, passwordHash string) (user *User, err *error)
	AddNewUser(user *User)
}

type UserDaoImpl struct {
	dbInst *gorm.DB
	schema string
}

func (dao *UserDaoImpl) Connect() {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_HOST"),
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_PORT"),
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_USER"),
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_DBNAME"),
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_PASSWORD"))

	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to db: %s", err.Error()))
	}
	db.SingularTable(true)
	dao.dbInst = db
	dao.schema = "schema_share_taxi_back"
}

func (dao UserDaoImpl) GetUserByEmail(email string) (user *User, err error) {
	dao.dbInst = dao.dbInst.Begin()
	dao.dbInst.Exec("set search_path to " + dao.schema)
	defer dao.dbInst.Rollback()

	user = &User{}
	dao.dbInst.Where("email = ?", email).First(user)
	return
}

func (dao UserDaoImpl) GetUserByEmailAndPassword(email string, passwordHash string) (user *User, err *error) {
	dao.dbInst = dao.dbInst.Begin()
	dao.dbInst.Exec("set search_path to " + dao.schema)
	defer dao.dbInst.Rollback()

	user = &User{}
	dao.dbInst.Where("email = ? AND password_hash = ?", email, passwordHash).First(user)
	return
}

func (dao UserDaoImpl) AddNewUser(user *User) {
	dao.dbInst = dao.dbInst.Begin()
	dao.dbInst.Exec("set search_path to " + dao.schema)
	defer dao.dbInst.Rollback()

	dao.dbInst.Create(&user)
}

func (dao UserDaoImpl) Disconnect() {
	dao.dbInst.Close()
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
