package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"example/go-rest-api/model"
)

func InitDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Book{})

	return db
}
