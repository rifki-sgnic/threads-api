package models

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_posts_user_created,priority:1" json:"user_id"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	Content   string    `gorm:"size:500;not null" json:"content"`
	CreatedAt time.Time `gorm:"index:idx_posts_user_created,priority:2" json:"created_at"`

	User   User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Parent *Post `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"-"`
}
