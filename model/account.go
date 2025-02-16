package model

type User struct {
	ID           uint32   `gorm:"primary_key" json:"id"`
	SessionToken string   `gorm:"unique" json:"session_token"`
	Username     string   `json:"username"`
	BestN        []Record `gorm:"serializer:json" json:"best_n"`
	Rks          float64  `json:"rks"`
}
