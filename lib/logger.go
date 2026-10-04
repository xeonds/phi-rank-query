package lib

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger writes one line per request to logFile.
func Logger(logFile string) gin.HandlerFunc {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.LstdFlags)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Printf("%s %s %s %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path, time.Since(start))
	}
}

// NoRoute serves files from dir for unmatched routes (SPA fallback is not
// needed because the frontend uses hash routing).
func NoRoute(router *gin.Engine, dir string) {
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir(dir))))
}
