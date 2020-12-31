package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRouteSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

type RouteTestSuite struct {
	suite.Suite
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

func (s *RouteTestSuite) TestGetSessionsRouteNoData() {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sessions", nil)
	s.NoError(err)
	router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	var response GetSessionJSON
	err2 := json.Unmarshal([]byte(w.Body.String()), &response)
	s.NoError(err2)

	// Ensure the Sessions array is empty when the DB is in a fresh state
	s.Assert().Equal(response.Sessions, make([]Session, 0))

}

func (s *RouteTestSuite) TestGetUsersRouteNoData() {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/users", nil)
	s.NoError(err)
	router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	var response GetUserJSON
	err2 := json.Unmarshal([]byte(w.Body.String()), &response)
	s.NoError(err2)

	// Ensure the Users array is empty when the DB is in a fresh state
	s.Assert().Equal(response.Users, make([]User, 0))

}

func (s *RouteTestSuite) TestGetSessionFeedbackRouteNoData() {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req1, err1 := http.NewRequest("GET", "/sessions/feedback", nil)
	router.ServeHTTP(w, req1)
	s.NoError(err1)
	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	var response GetFeedbackJSON
	err2 := json.Unmarshal([]byte(w.Body.String()), &response)
	s.NoError(err2)

	// TODO: Add db setup/teardown to the test logic (find the best place for this)
	// Ensure the Feedback array is empty when the DB is in a fresh state
	s.Assert().Equal(response.Feedback, make([]SessionFeedback, 0))
}
