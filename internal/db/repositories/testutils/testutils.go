package testutils

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetTestDb(t *testing.T, models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatalf("failed to migrate model: %v", err)
		}
	}
	return db
}
