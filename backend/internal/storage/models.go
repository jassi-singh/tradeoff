package storage

import (
	"time"
)

// PlayerModel represents the GORM model for the players table
type PlayerModel struct {
	ID                 string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username           string    `gorm:"type:varchar(255);not null" json:"username"`
	RefreshToken       string    `gorm:"type:text" json:"-"`
	RefreshTokenExpiry time.Time `gorm:"type:timestamp" json:"-"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (PlayerModel) TableName() string {
	return "players"
}
