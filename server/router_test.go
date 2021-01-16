package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRouteSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

type RouteTestSuite struct {
	suite.Suite
	router *gin.Engine
	mock   sqlmock.Sqlmock
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

// TestGetSessionsWithData tests that the /sessions route behaves as expected when there are Session records
func (s *RouteTestSuite) TestGetSessionsWithData() {
	// TODO: Figure out how to mock the database records (without actually inserting things in the DB - assume the DB will never be linked to the tests, which is expected in a production environment)
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

	// TODO: Add db setup/teardown to the test logic (find the best place for this)
	// Ensure the Feedback array is empty when the DB is in a fresh state
	s.Assert().Equal(response.Feedback, make([]SessionFeedback, 0))
}
