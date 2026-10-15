package repository

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FollowRepository interface {
	Create(follow *models.Follow) error
	Delete(followerID, followeeID uint) error
}

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

// Create implements [FollowRepository].
func (r *followRepository) Create(follow *models.Follow) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(follow).Error
}

// Delete implements [FollowRepository].
func (r *followRepository) Delete(followerID uint, followeeID uint) error {
	return r.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&models.Follow{}).Error
}
