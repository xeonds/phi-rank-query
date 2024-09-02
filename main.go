package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/xeonds/phi-plug-go/config"
	"github.com/xeonds/phi-plug-go/lib"
	"github.com/xeonds/phi-plug-go/model"
	"github.com/xeonds/phi-plug-go/service"
	"gorm.io/gorm"
)

func main() {
	config := lib.LoadConfig[config.Config]()
	db := lib.NewDB(&config.DatabaseConfig, func(db *gorm.DB) error {
		return db.AutoMigrate(&model.User{})
	})

	router := gin.Default()
	router.Use(lib.LoggerMiddleware(config.Server.LogFile))
	api := router.Group("/api/v1")
	api.POST("/b19", GetB19(config, db))
	api.POST("/bn", GetBN(config, db))
	api.GET("/leaderboard", GetLeaderboard(db))
	api.GET("/rank_table", GetRankTable(config))
	lib.AddStatic(router, []string{"./dist"})
	router.Run(config.Server.Port)
}

func GetB19(config *config.Config, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		post := new(struct {
			Session string `json:"session"`
		})
		if err := c.BindJSON(post); err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		if post.Session == "" {
			c.JSON(400, gin.H{"msg": "session is empty"})
			return
		}
		accountInfoDump, err := service.GetuserInfo(config, post.Session)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		accountInfo := new(model.GameAccount)
		_ = json.Unmarshal(accountInfoDump, accountInfo)
		userInfoDump, err := service.GetB19Info(config, post.Session)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		userInfo := new(model.GameSave)
		_ = json.Unmarshal(userInfoDump, userInfo)
		saveZip, err := service.GetSaveZip(config, post.Session, userInfo.Results[0].Gamefile.URL)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		game := service.DecryptSaveZip(saveZip)
		bn, rks, phi := service.CalcBNInfo(game, config)
		service.UpdateRank(db, &model.User{
			SessionToken: post.Session,
			Username:     accountInfo.Nickname,
			Rks:          rks,
		})
		c.JSON(200, gin.H{
			"player":    accountInfo.Nickname,
			"b19":       bn[:19],
			"rks":       rks,
			"phi":       phi,
			"date":      userInfo.Results[0].Gamefile.Updatedat,
			"challenge": game.GameProgress.ChallengeModeRank,
		})
	}
}

func GetBN(config *config.Config, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		post := new(struct {
			Session string `json:"session"`
		})
		if err := c.BindJSON(post); err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		if post.Session == "" {
			c.JSON(400, gin.H{"msg": "session is empty"})
			return
		}
		accountInfoDump, err := service.GetuserInfo(config, post.Session)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		accountInfo := new(model.GameAccount)
		_ = json.Unmarshal(accountInfoDump, accountInfo)
		userInfoDump, err := service.GetB19Info(config, post.Session)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		userInfo := new(model.GameSave)
		_ = json.Unmarshal(userInfoDump, userInfo)
		saveZip, err := service.GetSaveZip(config, post.Session, userInfo.Results[0].Gamefile.URL)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		game := service.DecryptSaveZip(saveZip)
		bn, rks, phi := service.CalcBNInfo(game, config)
		service.UpdateRank(db, &model.User{
			SessionToken: post.Session,
			Username:     accountInfo.Nickname,
			Rks:          rks,
		})
		db.Where("session_token = ?", post.Session).FirstOrCreate(&model.User{
			SessionToken: post.Session,
			Username:     accountInfo.Nickname,
			Rks:          rks,
		})
		c.JSON(200, gin.H{
			"player":    accountInfo.Nickname,
			"b19":       bn,
			"rks":       rks,
			"phi":       phi,
			"date":      userInfo.Results[0].Gamefile.Updatedat,
			"challenge": game.GameProgress.ChallengeModeRank,
		})
	}
}

func GetLeaderboard(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		users := service.GetLeaderboard(db)
		c.JSON(200, users)
	}
}

func GetRankTable(config *config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		rankTable, err := lib.LoadCSV(config.Data.Difficulty)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		songInfo, err := lib.LoadCSV(config.Data.Info)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		for i := range rankTable {
			rankTable[i]["title"] = songInfo[i]["song"]
		}
		c.JSON(200, rankTable)
	}
}
