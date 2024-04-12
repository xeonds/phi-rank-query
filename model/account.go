package model

type User struct {
	ID           uint32   `gorm:"primary_key" json:"id"`
	SessionToken string   `gorm:"unique" json:"session_token"`
	Username     string   `json:"username"`
	BestN        []Record `gorm:"foreignKey:UserID" json:"best_n"`
	Rks          float64  `json:"rks"`
}
