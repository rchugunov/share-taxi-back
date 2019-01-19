package gorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
)

type Connection struct {
	gorm.DB
}

func (conn *Connection) Connect() {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_HOST"),
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_PORT"),
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_USER"),
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_DBNAME"),
		os.Getenv("SHARE_TAXI_HEROKU_POSTGRES_PASSWORD"))

	c, err := gorm.Open("postgres", connectionString)
	*conn = Connection{*c}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to db: %s", err.Error()))
	}
	conn.SingularTable(true)
	conn.LogMode(true)
	schema := "schema_share_taxi_back"
	conn.Exec("set search_path to " + schema)
}

func (conn *Connection) Disconnect() error {
	return conn.Close()
}
