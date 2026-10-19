package repository

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FollowRepository interface {
	Create(follow *models.Follow) error
	Delete(followerID, followeeID uint) error
	CountFollowers(userID uint) (int64, error)
	CountFollowing(userID uint) (int64, error)
	IsFollowing(followerID, followeeID uint) (bool, error)
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

func (r *followRepository) CountFollowers(userID uint) (int64, error) {
	var n int64
	err := r.db.Model(&models.Follow{}).
		Where("followee_id = ?", userID).
		Count(&n).Error

	return n, err
}

func (r *followRepository) CountFollowing(userID uint) (int64, error) {
	var n int64
	err := r.db.Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Count(&n).Error

	return n, err
}

func (r *followRepository) IsFollowing(followerID, followeeID uint) (bool, error) {
	var n int64
	err := r.db.Model(&models.Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&n).Error

	return n > 0, err
}
