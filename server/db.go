package server

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CustomModel is the base model struct for all models used in this project - it emulates gorm.Model except that it uses a UUID for the ID field
type CustomModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime:mili"  json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:mili"  json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

// Counter database model used to track visit count for the given ID
type Counter struct {
	ID    uint `gorm:"primarykey" json:"-"`
	Visit uint `gorm:"default:0" json:"visit"`
}

// User database model representing the data collected for a user
type User struct {
	CustomModel
	SessionFeedback []SessionFeedback `json:"sessionFeedback"`
}

// Session database model representing the data for an arbitrary game session
type Session struct {
	CustomModel
	SessionFeedback []SessionFeedback `json:"feedback"`
}

// CreateSessionFeedbackInput represents the fields expected when the session feedback endpoint is hit with a POST request
type CreateSessionFeedbackInput struct {
	// The ID of the session being reviewed
	SessionID uuid.UUID `json:"sessionId"`
	// The ID of the user who left the review
	UserID  uuid.UUID `json:"userId"`
	Rating  int       `json:"rating"`
	Comment string    `json:"comment"`
}

// SessionFeedback database model representing the data for an arbitrary feedback response from a user about an arbitrary game session
type SessionFeedback struct {
	CustomModel
	// A rating of 1 to 5 (1 being the "worst", 5 being the "best")
	Rating int `gorm:"not null" json:"rating"`
	// An optional comment where the player can describe their experience in a small comment
	Comment string `json:"comment"`
	// FK
	SessionID uuid.UUID
	// FK
	UserID uuid.UUID
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&Counter{},
		&User{},
		&Session{},
		&SessionFeedback{},
	)
	if err != nil {
		panic(err)
	}
	return db
}
