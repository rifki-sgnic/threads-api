package repository

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/pagination"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id uint) (*models.Post, error)
	FindReplies(parentID uint, limit int) ([]models.Post, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
	Feed(userID uint, cursor *pagination.Cursor, limit int) ([]models.Post, error)
	CountReplies(postIDs []uint) (map[uint]int64, error)
	FindByUser(userID uint, cursor *pagination.Cursor, limit int) ([]models.Post, error)
	CountByUser(userID uint) (int64, error)
}

type postRepository struct {
	db *gorm.DB
}

type countRow struct {
	PostID uint
	Count  int64
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

func (r *postRepository) Feed(userID uint, cursor *pagination.Cursor, limit int) ([]models.Post, error) {
	following := r.db.Model(&models.Follow{}).
		Select("followee_id").
		Where("follower_id = ?", userID)

	q := r.db.Preload("User").
		Where("parent_id IS NULL").
		Where("user_id = ? OR user_id IN (?)", userID, following)

	q = withCursor(q, cursor)

	var posts []models.Post
	err := q.Order("created_at DESC, id DESC").Limit(limit).Find(&posts).Error
	return posts, err
}

func (r *postRepository) CountReplies(postIDs []uint) (map[uint]int64, error) {
	if len(postIDs) == 0 {
		return map[uint]int64{}, nil
	}

	var rows []countRow
	err := r.db.Model(&models.Post{}).
		Select("parent_id AS post_id, COUNT(*) as count").
		Where("parent_id IN ?", postIDs).
		Group("parent_id").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return toCountMap(rows), nil
}

func (r *postRepository) FindByUser(userID uint, cursor *pagination.Cursor, limit int) ([]models.Post, error) {
	q := r.db.Preload("User").Where("user_id = ? AND parent_id IS NULL", userID)
	q = withCursor(q, cursor)

	var posts []models.Post
	err := q.Order("created_at DESC, id DESC").Limit(limit).Find(&posts).Error
	return posts, err
}

func (r *postRepository) CountByUser(userID uint) (int64, error) {
	var n int64
	err := r.db.Model(&models.Post{}).
		Where("user_id = ? AND parent_id IS NULL", userID).
		Count(&n).Error

	return n, err
}

func toCountMap(rows []countRow) map[uint]int64 {
	m := make(map[uint]int64, len(rows))
	for _, r := range rows {
		m[r.PostID] = r.Count
	}
	return m
}

func withCursor(q *gorm.DB, c *pagination.Cursor) *gorm.DB {
	if c == nil {
		return q
	}
	return q.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
}
