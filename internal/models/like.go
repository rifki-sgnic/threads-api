package models

import "time"

type Like struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	PostID    uint      `gorm:"primaryKey" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	Post Post `gorm:"constraint:OnDelete:CASCADE" json:"post"`
}
