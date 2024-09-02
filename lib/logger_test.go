package lib_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xeonds/phi-plug-go/lib"
)

func TestLogger(t *testing.T) {
	// test logger middleware
	r := gin.New()
	r.Use(lib.LoggerMiddleware("test.log"))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}
	if w.Body.String() != "{\"message\":\"test\"}" {
		t.Errorf("expected body {\"message\":\"test\"}, got %s", w.Body.String())
	}
}
