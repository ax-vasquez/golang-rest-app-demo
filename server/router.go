package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const ContextKeyDB = "db"

func SetupRouter() *gin.Engine {
	r := gin.Default()
	addMiddleware(r)
	addDatabaseMiddleware(r)
	addRoutes(r)
	return r
}

// retrieves the database from the request context
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
	r.GET("/ping", func(c *gin.Context) {
		var counter Counter
		if err := GetDB(c).FirstOrCreate(&counter).Error; err != nil {
			panic(err)
		}
		counter.Visit++
		if err := GetDB(c).Save(&counter).Error; err != nil {
			panic(err)
		}
		c.JSON(200, &counter)
	})
}
