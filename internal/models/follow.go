package models

import "time"

type Follow struct {
	FollowerID uint      `gorm:"primaryKey" json:"follower_id"`
	FolloweeID uint      `gorm:"primaryKey" json:"followee_id"` // following
	CreatedAt  time.Time `json:"created_at"`

	Follower User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"-"`
	Followee User `gorm:"foreignKey:FolloweeID;constraint:OnDelete:CASCADE" json:"-"`
}
