package service

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/repository"
)

type LikeService interface {
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
}

type likeService struct {
	likes repository.LikeRepository
	posts repository.PostRepository
}

func NewLikeService(likes repository.LikeRepository, posts repository.PostRepository) LikeService {
	return &likeService{likes: likes, posts: posts}
}

// Like implements [LikeService].
func (s *likeService) Like(userID uint, postID uint) error {
	exists, err := s.posts.Exists(postID)
	if err != nil {
		return err
	}

	if !exists {
		return ErrPostNotFound
	}

	return s.likes.Create(&models.Like{UserID: userID, PostID: postID})
}

// Unlike implements [LikeService].
func (s *likeService) Unlike(userID uint, postID uint) error {
	return s.likes.Delete(userID, postID)
}
