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
		// TODO: Figure out how to get this to return related records (or figure out if there is a problem with adding the related records)
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
		var sessionFeedback SessionFeedback
		GetDB(c).Where(&SessionFeedback{SessionID: input.SessionID, UserID: input.UserID}).Find(&sessionFeedback)
		// Stop execution early (saving processing time) if the user has already provided feedback for this Session
		if sessionFeedback.ID != uuid.Nil {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "This user has already provided feedback for the given session"})
			return
		}
		var session Session
		// Defines session with the first Session record found by the given input.SessionID
		GetDB(c).First(&session, input.SessionID)
		var user User
		// Defines user with the first User record found by the given input.UserID
		GetDB(c).First(&user, input.UserID)
		sessionFeedback.ID = uuid.NewV4()
		sessionFeedback.Rating = input.Rating
		sessionFeedback.Comment = input.Comment
		session.SessionFeedback = []SessionFeedback{sessionFeedback}
		// Update the session with the feedback (inserts the feedback record into the DB)
		if err := GetDB(c).Updates(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Update the user with the feedback (doesn't perform insert this time since it already exists - just updates the User record)
		user.SessionFeedback = []SessionFeedback{sessionFeedback}
		if err := GetDB(c).Updates(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true, "message": "Thank you for your feedback!"})
	})
	// Get session feedback
	r.GET("/sessions/feedback", func(c *gin.Context) {
		var records []SessionFeedback
		GetDB(c).Find(&records)
		c.JSON(200, records)
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
}
