package server

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Counter database model used to track visit count for the given ID
type Counter struct {
	ID    uint `gorm:"primarykey" json:"-"`
	Visit uint `gorm:"default:0" json:"visit"`
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Counter{})
	if err != nil {
		panic(err)
	}
	return db
}
