package service

import (
	"errors"
	"strings"

	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("you are not allowed to do this")
	ErrEmptyContent = errors.New("content cannot be empty")
)

const maxReplies = 50

type CreatePostInput struct {
	UserID   uint
	Content  string
	ParentID *uint
}

type PostService interface {
	Create(in CreatePostInput) (*models.Post, error)
	GetByID(id uint) (*models.Post, error)
	GetReplies(postID uint) ([]models.Post, error)
	Delete(userID, postID uint) error
}

type postService struct {
	posts repository.PostRepository
}

func NewPostService(posts repository.PostRepository) PostService {
	return &postService{posts: posts}
}

// Create implements [PostService].
func (s *postService) Create(in CreatePostInput) (*models.Post, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return nil, ErrEmptyContent
	}

	if in.ParentID != nil {
		if _, err := s.GetByID(*in.ParentID); err != nil {
			return nil, err
		}
	}

	post := &models.Post{
		UserID:   in.UserID,
		ParentID: in.ParentID,
		Content:  content,
	}

	if err := s.posts.Create(post); err != nil {
		return nil, err
	}

	return s.GetByID(post.ID)
}

// GetByID implements [PostService].
func (s *postService) GetByID(id uint) (*models.Post, error) {
	post, err := s.posts.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return post, nil
}

// GetReplies implements [PostService].
func (s *postService) GetReplies(postID uint) ([]models.Post, error) {
	if _, err := s.GetByID(postID); err != nil {
		return nil, err
	}

	return s.posts.FindReplies(postID, maxReplies)
}

// Delete implements [PostService].
func (s *postService) Delete(userID uint, postID uint) error {
	post, err := s.GetByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return ErrForbidden
	}

	return s.posts.Delete(postID)
}
