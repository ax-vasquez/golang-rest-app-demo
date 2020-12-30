package server

import (
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

func (s *RouteTestSuite) TestGetSessionsRoute() {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sessions", nil)
	s.NoError(err)
	router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	s.Regexp(`\[\]`, w.Body.String())
}

func (s *RouteTestSuite) TestGetUsersRoute() {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/users", nil)
	s.NoError(err)
	router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	s.Regexp(`\[\]`, w.Body.String())
}

func (s *RouteTestSuite) TestGetSessionFeedbackRoute() {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sessions/feedback", nil)
	s.NoError(err)
	router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	// Test empty response - should always be empty array
	s.Regexp(`\[\]`, w.Body.String())
}
