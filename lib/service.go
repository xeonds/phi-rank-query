package lib

import (
	"context"
	"crypto/tls"
	"embed"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jordan-wright/email"
	"gorm.io/gorm"
)

func AddCRUD[T any](router gin.IRouter, path string, db *gorm.DB) *gin.RouterGroup {
	return APIBuilder(router, func(group *gin.RouterGroup) *gin.RouterGroup {
		group.GET("", GetAll[T](db, nil))
		group.GET("/:id", Get[T](db, nil))
		group.POST("", Create[T](db, nil))
		group.PUT("/:id", Update[T](db))
		group.DELETE("/:id", Delete[T](db))
		return group
	})(router, path)
}
func AddCRUDWithAuth[T any](router gin.IRouter, path string, db *gorm.DB, permLo, permHi int) *gin.RouterGroup {
	return APIBuilder(router, func(group *gin.RouterGroup) *gin.RouterGroup {
		group.GET("", GetAll[T](db, nil))
		group.GET("/:id", Get[T](db, nil))
		return group
	}, func(group *gin.RouterGroup) *gin.RouterGroup {
		// use should be in the first line of the function
		group.Use(JWTMiddleware(AuthPermission(permLo, permHi)))
		group.POST("", Create[T](db, nil))
		group.PUT("/:id", Update[T](db))
		group.DELETE("/:id", Delete[T](db))
		return group
	})(router, path)
}
func AddStatic(router *gin.Engine, staticFileDir []string) {
	for _, dir := range staticFileDir {
		router.NoRoute(gin.WrapH(http.FileServer(http.Dir(dir))))
	}
}
func AddStaticFS(router *gin.Engine, fs embed.FS) {
	router.NoRoute(gin.WrapH(http.FileServer(http.FS(fs))))
}
func AddCaptchaAPI(router gin.IRouter, path string, conf1 MailConfig, conf2 CaptchaConfig, rdb *redis.Client) *gin.RouterGroup {
	return APIBuilder(router, func(group *gin.RouterGroup) *gin.RouterGroup {
		group.POST("/gen_captcha", HandleMailSendCaptcha(conf1, conf2, rdb))
		group.POST("/verify_captcha", HandleCaptchaVerify(rdb))
		return group
	})(router, path)
}

// handlers for gorm
func Create[T any](db *gorm.DB, process func(*gorm.DB, *T) *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		d := new(T)
		if err := c.ShouldBindJSON(d); err != nil {
			c.AbortWithStatus(404)
			log.Println("[gorm]parse creation data failed: ", err)
		} else {
			if process != nil {
				if process(db, d).Error != nil {
					c.AbortWithStatus(404)
					log.Println("[gorm] create data process failed: ", err)
				}
			} else if err := db.Create(d).Error; err != nil {
				c.AbortWithStatus(404)
				log.Println("[gorm] create data failed: ", err)
			}
			c.JSON(200, d)
		}
	}
}
func Get[T any](db *gorm.DB, process func(*gorm.DB, *gin.Context) *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, d := c.Param("id"), new(T)
		if process != nil {
			if process(db, c).First(d).Error != nil {
				c.AbortWithStatus(404)
				log.Println("[gorm] db query process failed")
			}
		} else if err := db.Where("id = ?", id).First(d).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		}
		c.JSON(200, d)
	}
}
func GetAll[T any](db *gorm.DB, process func(*gorm.DB, *gin.Context) *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		d := new([]T)
		if process != nil {
			if process(db, c).Find(d).Error != nil {
				c.AbortWithStatus(404)
				log.Println("[gorm] db query all process failed")
			}
		} else if err := db.Find(d).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		}
		c.JSON(200, d)
	}
}
func Update[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var d T
		if err := db.Save(&d).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		} else {
			c.JSON(200, d)
		}
	}
}
func Delete[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		var d T
		if err := db.Where("id = ?", id).Delete(&d).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		} else {
			c.JSON(200, d)
		}
	}
}
func HandleFind[T any](queryProcess func(c *gin.Context) *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		query := new(T)
		err := queryProcess(c).First(query).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Query Content Not Found",
			})
		} else {
			c.JSON(200, query)
		}
	}
}
func HandleFindAll[T any](queryProcess func(c *gin.Context) *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var query []T
		err := queryProcess(c).Find(&query).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Query Content Not Found",
			})
		} else {
			c.JSON(200, query)
		}
	}
}

// 验证码服务，使用redis存储
func HandleMailSendCaptcha(mailConfig MailConfig, captchaConfig CaptchaConfig, rdb *redis.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		mailTo := c.Query("mail") // TODO: 增加对mail的合法性验证
		codeId, code := GenerateCaptcha(captchaConfig.CaptchaLength)
		rdb.Set(context.Background(), codeId, code, time.Minute*time.Duration(captchaConfig.CaptchaAlive))
		e := email.NewEmail()
		e.From = "Get <" + mailConfig.MailServer + ">"
		e.To = []string{mailTo}
		e.Subject = "验证码"
		e.HTML = []byte("你的验证码为：<h1>" + code + "</h1>")
		err := e.SendWithTLS(mailConfig.MailServer+mailConfig.MailServerPort, smtp.PlainAuth("", mailConfig.MailUserName, mailConfig.MailPassword, mailConfig.MailServer),
			&tls.Config{InsecureSkipVerify: true, ServerName: mailConfig.MailServer})
		if err != nil {
			c.Error(err)
		} else {
			c.JSON(200, gin.H{"id": codeId, "code": code})
		}
	}
}
func HandleCaptchaVerify(rdb *redis.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		// TODO: 增加对id和code的合法性验证
		codeId := c.Query("id")
		code := c.Query("code")
		if VerifyCaptcha(codeId, code, rdb) {
			c.JSON(200, gin.H{
				"message": "验证码正确",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "验证码错误",
			})
		}
	}
}
