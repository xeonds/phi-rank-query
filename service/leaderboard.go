package service

import (
	"github.com/xeonds/phi-plug-go/model"
	"gorm.io/gorm"
)

// Leaderboard returns users ordered by rks descending.
func Leaderboard(db *gorm.DB) []model.User {
	var users []model.User
	db.Order("rks desc").Select("id", "username", "rks").Find(&users)
	return users
}

// UpdateUser inserts or updates a user's rks after a query.
func UpdateUser(db *gorm.DB, session, username string, rks float64) {
	var u model.User
	if err := db.Where("session_token = ?", session).FirstOrCreate(&u, model.User{SessionToken: session}).Error; err != nil {
		return
	}
	u.Username = username
	u.Rks = rks
	db.Save(&u)
}
