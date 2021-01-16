package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SessionFeedbackGetterAll func(c *gin.Context, records *[]SessionFeedback)
type SessionFeedbackGetterByRating func(c *gin.Context, rating int, records *[]SessionFeedback)
type SessionFeedbackGetterBySessionID func(c *gin.Context, sessionID string, records *[]SessionFeedback)
type SessionFeedbackGetterBySessionIDAndRating func(c *gin.Context, sessionID string, rating int, records *[]SessionFeedback)

// SessionFeedbackGetter facilitates getting SessionFeedback
//
// This object exists so that it can be used in both a production and testing capacity. In production,
// this object is initialized with the methods that actually make calls to the database. In testing,
// this object should be initialized with mock methods that return the data as expected.
//
// Note that the methods used to initialize this object ARE NOT under test; they are merely *used*
// in the actual test (which is only concerned with the logic of getSessionFeedback)
//
// See https://stackoverflow.com/questions/19167970/mock-functions-in-go
type SessionFeedbackGetter struct {
	all                      SessionFeedbackGetterAll
	by_rating                SessionFeedbackGetterByRating
	by_session_id            SessionFeedbackGetterBySessionID
	by_session_id_and_rating SessionFeedbackGetterBySessionIDAndRating
}

func NewSessionFeedbackGetter(getAll SessionFeedbackGetterAll, byRating SessionFeedbackGetterByRating, bySessionID SessionFeedbackGetterBySessionID, bySessionIDAndRating SessionFeedbackGetterBySessionIDAndRating) *SessionFeedbackGetter {
	return &SessionFeedbackGetter{
		all:                      getAll,
		by_rating:                byRating,
		by_session_id:            bySessionID,
		by_session_id_and_rating: bySessionIDAndRating,
	}
}

// getSessionFeedbackBySessionIdAndRating gets SessionFeedback records filtered by SessionID and Rating
func getSessionFeedbackBySessionIdAndRating(c *gin.Context, sessionID string, rating int, records *[]SessionFeedback) {
	// SELECT * FROM session_feedbacks WHERE session_id = ? AND rating = ?
	GetDB(c).Where("session_id = ? AND rating = ?", sessionID, rating).Find(&records)
}

// getSessionFeedbackBySessionId gets SessionFeedback records filtered by the given SessionID
func getSessionFeedbackBySessionId(c *gin.Context, sessionID string, records *[]SessionFeedback) {
	// SELECT * FROM session_feedbacks WHERE session_id = ?
	GetDB(c).Where("session_id = ?", sessionID).Find(&records)
}

// getSessionFeedbackByRating gets SessionFeedback records fultered by the given rating (includes all sessions)
func getSessionFeedbackByRating(c *gin.Context, rating int, records *[]SessionFeedback) {
	// SELECT * FROM session_feedbacks WHERE rating = ?
	GetDB(c).Where("rating = ?", rating).Find(&records)
}

func getAllSessionFeedback(c *gin.Context, records *[]SessionFeedback) {
	// SELECT * FROM session_feedbacks
	GetDB(c).Find(&records)
}

// getSessionFeedback handles the logic for GET requests sent to the /sessions/feedback endpoint - accepts sessionId and/or rating as query parameters
func getSessionFeedback(c *gin.Context, sfg SessionFeedbackGetter) {
	var records []SessionFeedback
	query := c.Request.URL.Query()
	var sessionID = query["sessionId"]
	var rating = query["rating"]
	if sessionID != nil || rating != nil { // Query is filtered by sessionId and/or rating
		if sessionID != nil && rating != nil { // By sessionId and rating
			if ratingInt, err := strconv.Atoi(rating[0]); err == nil {
				if ratingIsValid(ratingInt) {
					sfg.by_session_id_and_rating(c, sessionID[0], ratingInt, &records)
					c.JSON(200, gin.H{"feedback": &records})
					return
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be an integer from 1 through 5"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else if sessionID == nil && rating != nil { // By rating only
			if ratingInt, err := strconv.Atoi(rating[0]); err == nil {
				if ratingIsValid(ratingInt) {
					sfg.by_rating(c, ratingInt, &records)
					c.JSON(200, gin.H{"feedback": &records})
					return
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be an integer from 1 through 5"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else if sessionID != nil && rating == nil { // By sessionId only
			sfg.by_session_id(c, sessionID[0], &records)
			c.JSON(200, gin.H{"feedback": &records})
			return
		}
	} else { // Unfiltered query
		sfg.all(c, &records)
		c.JSON(200, gin.H{"feedback": &records})
		return
	}
}
