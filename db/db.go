package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	user     = "abidlutfi"
	password = "0000"
	dbname   = "products"
	port     = "5432"
	timeZone = "Asia/Jakarta"
	sslMode  = "disable"
)

var (
	db  *gorm.DB
	err error
)

func InitializeDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslMode, timeZone,
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("error connecting to db: (%s)", err.Error())
	}
}

func GetDB() *gorm.DB {
	return db
}
