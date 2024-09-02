package lib

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// 日志中间件，记录客户端IP，请求方法，请求路径，请求耗时，请求头的Aurhorization字段, 将日志保存在path/to/logFile.log中
func LoggerMiddleware(logFile string) gin.HandlerFunc {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.LstdFlags)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Printf("%s %s %s %s %s\n", c.ClientIP(), c.Request.Method, c.Request.URL.Path, time.Since(start), c.GetHeader("Authorization"))
	}
}

// JWT中间件
// 支持添加权限校验（返回error表示校验失败），以及上下文操作
func JWTMiddleware(authToken func(*gin.Context, UserClaim) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		parsed, err := ParseToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if parsed.Valid() != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if authToken != nil && authToken(c, *parsed) != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

// 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// finish this golang cors middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Next()
	}
}
