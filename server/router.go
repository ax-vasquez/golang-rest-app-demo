package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
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
	// Set logrus to use JSON formatting (e.g., "structured formatting") - more easily-consumable by services like GCP Log services
	log.SetFormatter(&log.JSONFormatter{})
	// Custom logging middleware
	r.Use(Logger())

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

type SessionGetter func()
type UserGetter func()

// In production, this init method should be called with DB-connected code - in test, it should init with mocked calls
type ResourceGetterInit func()

// ResourceGetter is linked to the DB in production and mocked in tests
type ResourceGetter struct {
	init              ResourceGetterInit
	sessions          SessionGetter
	session_feedbacks SessionFeedbackGetter
	users             UserGetter
}

func ping(c *gin.Context) {
	c.String(200, "ping")
	return
}

// ratingIsValid checks if a given value is between 1 and 5 (the accepted rating range)
func ratingIsValid(i int) bool {
	if i >= 1 && i <= 5 {
		return true
	}
	return false
}

func initResourceGetter() *ResourceGetter {
	return &ResourceGetter{}
}

// getAllSessions gets all Session records
func getAllSessions(c *gin.Context, records *[]Session) {
	// SELECT * FROM sessions
	GetDB(c).Find(&records)
}

// getAllUsers gets all User records
func getAllUsers(c *gin.Context, records *[]User) {
	// SELECT * FROM users
	GetDB(c).Find(&records)
}

// TODO: Make all endpoints support filtering (eventually)
// getResources is a convenience method used to contain the logic (at a high level) for all GET endpoints
func getResources(c *gin.Context) {
	var users []User
	var sessions []Session
	switch c.FullPath() {
	case "/sessions":
		getAllSessions(c, &sessions)
		c.JSON(http.StatusOK, gin.H{"sessions": sessions})
		return
	case "/users":
		getAllUsers(c, &users)
		c.JSON(http.StatusOK, gin.H{"users": users})
		return
	case "/sessions/feedback":
		sfg := NewSessionFeedbackGetter(getAllSessionFeedback, getSessionFeedbackByRating, getSessionFeedbackBySessionId, getSessionFeedbackBySessionIdAndRating)
		getSessionFeedback(c, *sfg)
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unrecognized route!"})
		return
	}
}

func createSession(c *gin.Context) {
	var session Session
	session.ID = uuid.NewV4()
	if err := GetDB(c).Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"session": &session})
}

func deleteSession(c *gin.Context) {
	query := c.Request.URL.Query()
	var session Session
	if err := GetDB(c).Where("id = ?", query["id"]).Find(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	// Check if the user even exists - return early if not
	if session.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Session does not exist"})
		return
	}
	// Attempt to delete the user (return an error if something bad happens)
	if err := GetDB(c).Delete(&Session{}, query["id"]).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Session deleted successfully!"})
}

func createSessionFeedback(c *gin.Context) {
	var input CreateSessionFeedbackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// If any required fields are invalid, return before doing any processing
	if input.Rating < 1 || input.Rating > 5 || input.SessionID == uuid.Nil || input.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid values for query parameters - sessionId and userId must be defined and rating must be 1 - 5"})
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
	return
}

func deleteSessionFeedback(c *gin.Context) {
	query := c.Request.URL.Query()
	var sessionFeedback SessionFeedback
	if err := GetDB(c).Where("id = ?", query["id"]).Find(&sessionFeedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	// Check if the user even exists - return early if not
	if sessionFeedback.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "SessionFeedback does not exist"})
		return
	}
	// Attempt to delete the user (return an error if something bad happens)
	if err := GetDB(c).Delete(&SessionFeedback{}, query["id"]).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "SessionFeedback deleted successfully!"})
	return
}

func createUser(c *gin.Context) {
	var user User
	user.ID = uuid.NewV4()
	if err := GetDB(c).Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"user": &user})
	return
}

func deleteUser(c *gin.Context) {
	query := c.Request.URL.Query()
	var user User
	if err := GetDB(c).Where("id = ?", query["id"]).Find(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	// Check if the user even exists - return early if not
	if user.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User does not exist"})
		return
	}
	// Attempt to delete the user (return an error if something bad happens)
	if err := GetDB(c).Delete(&User{}, query["id"]).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User deleted successfully!"})
	return
}

// adds routes to the server
func addRoutes(r *gin.Engine) {

	r.GET("/ping", ping)
	r.GET("/users", getResources)
	r.GET("/sessions", getResources)
	// TODO: Look into how to do wildcards in routes with gin
	r.GET("/sessions/feedback", getResources)
	r.POST("/users/create", createUser)
	r.POST("/sessions/create", createSession)
	r.POST("/sessions/feedback/create", createSessionFeedback)
	r.DELETE("/users", deleteUser)
	r.DELETE("/sessions", deleteSession)
	r.DELETE("/sessions/feedback", deleteSessionFeedback)
}
