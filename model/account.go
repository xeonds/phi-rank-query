package model

// User is a player stored for the leaderboard.
type User struct {
	ID           uint32  `gorm:"primary_key" json:"id"`
	SessionToken string  `gorm:"unique" json:"-"`
	Username     string  `json:"username"`
	Rks          float64 `json:"rks"`
}
