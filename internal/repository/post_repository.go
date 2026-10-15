package repository

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id uint) (*models.Post, error)
	FindReplies(parentID uint, limit int) ([]models.Post, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// Create implements [PostRepository].
func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// FindByID implements [PostRepository].
func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("User").First(&post, id).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// FindReplies implements [PostRepository].
func (r *postRepository) FindReplies(parentID uint, limit int) ([]models.Post, error) {
	var posts []models.Post

	err := r.db.Preload("User").
		Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Limit(limit).
		Find(&posts).Error

	return posts, err
}

// Delete implements [PostRepository].
func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error

}

func (r *postRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
