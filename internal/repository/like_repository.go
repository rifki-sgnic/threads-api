package repository

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LikeRepository interface {
	Create(like *models.Like) error
	Delete(userID, postID uint) error
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

// Create implements [LikeRepository].
func (r *likeRepository) Create(like *models.Like) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(like).Error
}

// Delete implements [LikeRepository].
func (r *likeRepository) Delete(userID uint, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.Like{}).Error
}
