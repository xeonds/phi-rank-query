package service

import (
	"github.com/xeonds/phi-plug-go/model"
	"gorm.io/gorm"
)

func GetLeaderboard(db *gorm.DB) *[]model.User {
	users := new([]model.User)
	db.Order("rks desc").Find(users).Select("id, username, rks")
	return users
}
