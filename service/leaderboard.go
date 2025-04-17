package service

import (
	"github.com/xeonds/phi-plug-go/model"
	"gorm.io/gorm"
)

func GetLeaderboard(db *gorm.DB) *[]model.User {
	users := new([]model.User)
	db.Order("rks desc").Select("id", "username", "rks").Find(users)
	return users
}

func UpdateRank(db *gorm.DB, user *model.User) error {
	u := new(model.User)
	if err := db.FirstOrCreate(u, model.User{SessionToken: user.SessionToken}).Error; err != nil {
		return err
	}
	u.Rks = user.Rks
	u.Username = user.Username
	u.BestN = user.BestN

	if err := db.Save(u).Error; err != nil {
		return err
	}
	return nil
}
