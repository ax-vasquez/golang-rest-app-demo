package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RouteTestSuite struct {
	suite.Suite
	router *gin.Engine
	mock   sqlmock.Sqlmock
}

type KeyValuePair struct {
	Key   string
	Value string
}

// SetupTest performs the shared setup logic for all tests in the RouteTestSuite
func (suite *RouteTestSuite) SetupTest() {
	// Mock this method
	suite.router = SetupMockRouter(suite)
}

// GetUserJSON is used when unmarshalling the response from GET /users
type GetUserJSON struct {
	Users []User `json:"users"`
}

// GetSessionJSON is used when unmarshalling the response from GET /sessions
type GetSessionJSON struct {
	Sessions []Session `json:"sessions"`
}

// GetFeedbackJSON is used when unmarshalling the response from GET /sessions/feedback
type GetFeedbackJSON struct {
	Feedback []SessionFeedback `json:"feedback"`
}

type CreateUserJSON struct {
	User User `json:"user"`
}

type CreateSessionJSON struct {
	Session Session `json:"session"`
}

type CreateSessionFeedbackJSON struct {
	Success         bool            `json:"success"`
	Message         string          `json:"message"`
	SessionFeedback SessionFeedback `json:"sessionFeedback"`
}

func TestRouteSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

func SetupMockRouter(s *RouteTestSuite) *gin.Engine {
	r := gin.Default()
	addMiddleware(r)
	addMockDatabaseMiddleware(r, s)
	addRoutes(r)
	return r
}

func initMockDB() *gorm.DB {
	// Open an in-memory SQLite database (will cease to exist once tests are done)
	// see https://gorm.io/docs/connecting_to_the_database.html#SQLite
	dialector := sqlite.Open("file::memory:")
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = gdb.AutoMigrate(
		&Counter{},
		&User{},
		&Session{},
		&SessionFeedback{},
	)
	if err != nil {
		panic(err)
	}
	return gdb
}

func addMockDatabaseMiddleware(r *gin.Engine, s *RouteTestSuite) {
	gdb := initMockDB()

	// Add database to our context
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyDB, gdb)
	})
}

// TestCreateSession ensures the /sessions/create endpoint creates a Session as expected
func (s *RouteTestSuite) TestCreateSession() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/sessions/create", nil)
	s.NoError(err)
	s.router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	var response CreateSessionJSON
	err = json.Unmarshal([]byte(w.Body.String()), &response)

	// Make sure the Session was assigned a non-nil UUID
	s.NotEqual(response.Session.ID, uuid.Nil)
}

// TestGetSessionsRouteNoData tests that the /sessions route behaves as expected when there are no Session records
func (s *RouteTestSuite) TestGetSessionsRouteNoData() {

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sessions", nil)
	s.NoError(err)
	s.router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	var response GetSessionJSON
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	s.NoError(err)

	// Gorm returns empty slices INSTEAD OF nil slices
	var sessions []Session = make([]Session, 0)
	// Ensure the Sessions array is empty when the DB is in a fresh state
	s.Assert().Equal(response.Sessions, sessions)
}

// TestCreateUser ensures the /users/create endpoint creates a User as expected
func (s *RouteTestSuite) TestCreateUser() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/users/create", nil)
	s.NoError(err)
	s.router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	var response CreateUserJSON
	err = json.Unmarshal([]byte(w.Body.String()), &response)

	// Make sure the User was assigned a non-nil UUID
	s.NotEqual(response.User.ID, uuid.Nil)
}

func (s *RouteTestSuite) TestGetUsersRouteNoData() {

	// See https://golang.org/pkg/net/http/httptest/#ResponseRecorder
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/users", nil)
	s.NoError(err)
	s.router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	var response GetUserJSON
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	s.NoError(err)

	// Gorm returns empty slices INSTEAD OF nil slices
	var users []User = make([]User, 0)
	// Ensure the Users array is empty when the DB is in a fresh state
	s.Assert().Equal(response.Users, users)

}

// createPostBodyString is a helper method to create a JSON string for a given set of KeyValuePair objects
//
// This method assumes all int values are intended to be processes as actual integers and NOT strings (and doesn't
// wrap the resulting value in parentheses).
func createPostBodyString(pairs ...KeyValuePair) (body string) {
	if len(pairs) > 1 {
		base := ""
		for i := 0; i < len(pairs); i++ {
			if i == 0 {
				base = fmt.Sprintf("{")
			}
			if i == (len(pairs) - 1) {
				valAsInt, err := strconv.Atoi(pairs[i].Value)
				if err == nil { // If the value is an int
					base = fmt.Sprintf("%s\"%s\":%d}", base, pairs[i].Key, valAsInt)
				} else {
					base = fmt.Sprintf("%s\"%s\":\"%s\"}", base, pairs[i].Key, pairs[i].Value)
				}
			} else {
				valAsInt, err := strconv.Atoi(pairs[i].Value)
				if err == nil { // If the value is an int
					base = fmt.Sprintf("%s\"%s\":%d,", base, pairs[i].Key, valAsInt)
				} else {
					base = fmt.Sprintf("%s\"%s\":\"%s\",", base, pairs[i].Key, pairs[i].Value)
				}
			}
		}
		return base
	} else if len(pairs) == 0 {
		return "{}" // empty POST body
	}
	i, err := strconv.Atoi(pairs[0].Value)
	if err == nil { // If the value is an int
		return fmt.Sprintf("{\"%s\":%d}", pairs[0].Key, i)
	} else {
		return fmt.Sprintf("{\"%s\":\"%s\"}", pairs[0].Key, pairs[0].Value)
	}
}

func (s *RouteTestSuite) TestCreateSessionFeedback() {
	var user User
	var session Session
	var sessionFeedback SessionFeedback

	// Create a User to write the feedback
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/users/create", nil)
	s.NoError(err)
	s.router.ServeHTTP(w, req)
	s.Equal(200, w.Code)

	var createUserResponse CreateUserJSON
	err = json.Unmarshal([]byte(w.Body.String()), &createUserResponse)
	s.NoError(err)
	user = createUserResponse.User

	// Create a Session to write feedback for
	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/sessions/create", nil)
	s.NoError(err)
	s.router.ServeHTTP(w, req)
	s.Equal(200, w.Code)

	var createSessionResponse CreateSessionJSON
	err = json.Unmarshal([]byte(w.Body.String()), &createSessionResponse)
	s.NoError(err)
	session = createSessionResponse.Session

	// Create the SessionFeedback
	var sessionIdKeyValuePair KeyValuePair
	sessionIdKeyValuePair.Key = "sessionId"
	sessionIdKeyValuePair.Value = session.ID.String()

	var userIdKeyValuePair KeyValuePair
	userIdKeyValuePair.Key = "userId"
	userIdKeyValuePair.Value = user.ID.String()

	var ratingKeyValuePair KeyValuePair
	ratingKeyValuePair.Key = "rating"
	ratingKeyValuePair.Value = "3"

	w = httptest.NewRecorder()
	postBodyString := createPostBodyString(sessionIdKeyValuePair, userIdKeyValuePair, ratingKeyValuePair)
	postBodyJSON := []byte(postBodyString)
	req, err = http.NewRequest("POST", "/sessions/feedback/create", bytes.NewBuffer(postBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	s.NoError(err)
	s.router.ServeHTTP(w, req)
	s.Equal(200, w.Code)

	var createSessionFeedbackResponse CreateSessionFeedbackJSON
	err = json.Unmarshal([]byte(w.Body.String()), &createSessionFeedbackResponse)
	s.NoError(err)
	sessionFeedback = createSessionFeedbackResponse.SessionFeedback

	// SessionFeedback should not have a nil UUID
	s.NotEqual(sessionFeedback.ID, uuid.Nil)
	// SessionFeedback.UserID should match the given userId
	s.Equal(user.ID, sessionFeedback.UserID)
	// SessionFeedback.SessionID should match the given sessionId
	s.Equal(session.ID, sessionFeedback.SessionID)
}

func (s *RouteTestSuite) TestGetSessionFeedbackRouteNoData() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sessions/feedback", nil)
	s.router.ServeHTTP(w, req)
	s.NoError(err)
	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	var response GetFeedbackJSON
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	s.NoError(err)

	// Ensure the Feedback array is empty when the DB is in a fresh state
	s.Assert().Equal(response.Feedback, make([]SessionFeedback, 0))
}
