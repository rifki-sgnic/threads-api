package repository

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LikeRepository interface {
	Create(like *models.Like) error
	Delete(userID, postID uint) error
	CountByPosts(postIDs []uint) (map[uint]int64, error)
	LikedByUser(userID uint, postIDs []uint) (map[uint]bool, error)
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

func (r *likeRepository) CountByPosts(postIDs []uint) (map[uint]int64, error) {
	if len(postIDs) == 0 {
		return map[uint]int64{}, nil
	}

	var rows []countRow
	err := r.db.Model(&models.Like{}).
		Select("post_id, COUNT(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return toCountMap(rows), nil
}

func (r *likeRepository) LikedByUser(userID uint, postIDs []uint) (map[uint]bool, error) {
	if len(postIDs) == 0 {
		return map[uint]bool{}, nil
	}

	var ids []uint

	err := r.db.Model(&models.Like{}).
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Pluck("post_id", &ids).Error

	if err != nil {
		return nil, err
	}

	liked := make(map[uint]bool, len(ids))
	for _, id := range ids {
		liked[id] = true
	}

	return liked, nil
}
