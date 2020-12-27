package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ContextKeyDB is the key name for the database within the Gin context
const ContextKeyDB = "db"

// SetupRouter completes setup of the router, middleware, db middleware and routes and returns the default Engine instance
func SetupRouter() *gin.Engine {
	r := gin.Default()
	addMiddleware(r)
	addDatabaseMiddleware(r)
	addRoutes(r)
	return r
}

// GetDB retrieves the database from the request context
func GetDB(c *gin.Context) *gorm.DB {
	value, ok := c.Get(ContextKeyDB)
	if !ok {
		panic("database not found in context")
	}
	db, ok := value.(*gorm.DB)
	if !ok {
		panic("database was not the correct type")
	}
	return db
}

// adds basic middleware
func addMiddleware(r *gin.Engine) {
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
}

// adds the database to the context, it can be retrieved in routes by using GetDB
func addDatabaseMiddleware(r *gin.Engine) {
	db := initDB()
	// Add database to our context
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyDB, db)
	})
}

// adds routes to the server
func addRoutes(r *gin.Engine) {
	// Test ping route
	// r.GET("/ping", func(c *gin.Context) {
	// 	var counter Counter
	// 	if err := GetDB(c).FirstOrCreate(&counter).Error; err != nil {
	// 		panic(err)
	// 	}
	// 	counter.Visit++
	// 	if err := GetDB(c).Save(&counter).Error; err != nil {
	// 		panic(err)
	// 	}
	// 	c.JSON(200, &counter)
	// })
	// Test ping route
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "ping")
	})
	// List all sessions
	r.GET("/sessions", func(c *gin.Context) {
		var records []Session
		GetDB(c).Find(&records)
		c.JSON(200, records)
	})
	// Create an arbitrary game session
	r.POST("/sessions/create", func(c *gin.Context) {
		var session Session
		session.ID = uuid.NewV4()
		if err := GetDB(c).Create(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, &session)
	})
	// TODO: Look into how to do wildcards in routes with gin
	// Leave feedback for a given session ID
	r.POST("/sessions/feedback", func(c *gin.Context) {
		var input CreateSessionFeedbackInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		session := GetDB(c).Find(&Session{}, input.SessionID)
		// TODO: Figure out why this is broken - I've already made sure input is correct through independent testing - the issue is here
		GetDB(c).Model(&session).Association("Feedback").Append(&SessionFeedback{Rating: input.Rating, Comment: input.Comment})
	})
	// List all users
	r.GET("/users", func(c *gin.Context) {
		var records []User
		GetDB(c).Find(&records)
		c.JSON(200, records)
	})
	// Create an arbitrary user
	r.POST("/users/create", func(c *gin.Context) {
		var user User
		user.ID = uuid.NewV4()
		if err := GetDB(c).Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, &user)
	})
	// Reset the local database (TESTING PURPOSES ONLY)
	r.POST("/data/flush", func(c *gin.Context) {
		if err := GetDB(c).Migrator().DropTable(
			&CustomModel{},
			&Counter{},
			&User{},
			&Session{},
			&SessionFeedback{},
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := GetDB(c).AutoMigrate(
			&CustomModel{},
			&Counter{},
			&User{},
			&Session{},
			&SessionFeedback{},
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}
