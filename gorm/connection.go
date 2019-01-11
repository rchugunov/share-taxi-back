package gorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
)

type Connection struct {
	dbInst *gorm.DB
}

func (Connection) GetNewDBInst() (db *gorm.DB) {
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
	db.LogMode(true)
	schema := "schema_share_taxi_back"
	db.Exec("set search_path to " + schema)
	return
}

func (conn *Connection) Connect() {
	conn.dbInst = conn.GetNewDBInst()
}

func (conn Connection) Disconnect() {
	conn.dbInst.Close()
}
