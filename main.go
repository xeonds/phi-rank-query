package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/xeonds/phi-plug-go/config"
	"github.com/xeonds/phi-plug-go/lib"
	"github.com/xeonds/phi-plug-go/service"
	"gorm.io/gorm"
)

func main() {
	config := lib.LoadConfig[config.Config]()
	db := lib.NewDB(&config.DatabaseConfig, func(db *gorm.DB) error {
		return db.AutoMigrate(&lib.User{})
	})

	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/b19", GetB19(config))
	// TODO: 登陆后查询
	lib.AddLoginAPI(api, "/login", db)
	lib.AddStatic(router, []string{"./dist"})
	router.Run(config.Server.Port)
}

func GetB19(config *config.Config) func(c *gin.Context) {
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
		accountInfo := new(service.GameAccount)
		_ = json.Unmarshal(accountInfoDump, accountInfo)
		userInfoDump, err := service.GetB19Info(config, post.Session)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		userInfo := new(service.GameSave)
		_ = json.Unmarshal(userInfoDump, userInfo)
		saveZip, err := service.GetSaveZip(config, post.Session, userInfo.Results[0].Gamefile.URL)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		game := service.DecryptSaveZip(saveZip)
		b19, rks, phi := service.CalcBNInfo(game, config, 19)
		c.JSON(200, gin.H{
			"player":    accountInfo.Nickname,
			"b19":       b19,
			"rks":       rks,
			"phi":       phi,
			"date":      userInfo.Results[0].Gamefile.Updatedat,
			"challenge": game.GameProgress.ChallengeModeRank,
		})
	}
}
