package dbconn

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DBSystem *gorm.DB

func OpenDB(dsn string) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database")
	}

	DBSystem = db
}
