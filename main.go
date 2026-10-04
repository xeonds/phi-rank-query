package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xeonds/phi-plug-go/config"
	"github.com/xeonds/phi-plug-go/lib"
	"github.com/xeonds/phi-plug-go/model"
	"github.com/xeonds/phi-plug-go/service"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	db := lib.NewDB(&cfg.Database, func(db *gorm.DB) error {
		return db.AutoMigrate(&model.User{})
	})
	client, err := service.New(cfg)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(lib.Logger(cfg.Server.LogFile))

	api := router.Group("/api/v1")
	api.GET("/rank_table", func(c *gin.Context) {
		c.JSON(http.StatusOK, client.RankTable())
	})
	api.GET("/leaderboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, service.Leaderboard(db))
	})
	api.POST("/player", handlePlayer(client, db))

	lib.NoRoute(router, "./dist")
	router.Run(cfg.Server.Port)
}

func handlePlayer(client *service.Client, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Session string `json:"session"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.Session == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "session is required"})
			return
		}
		result, err := client.Query(body.Session)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		service.UpdateUser(db, body.Session, result.Player, result.Rks)
		c.JSON(http.StatusOK, result)
	}
}
