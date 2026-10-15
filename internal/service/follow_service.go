package service

import (
	"errors"

	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrCannotFollowSelf = errors.New("cannot follow yourself")
)

type FollowService interface {
	Follow(followerID uint, username string) error
	Unfollow(followerID uint, username string) error
}

type followService struct {
	follows repository.FollowRepository
	users   repository.UserRepository
}

func NewFollowService(follows repository.FollowRepository, users repository.UserRepository) FollowService {
	return &followService{follows: follows, users: users}
}

// Follow implements [FollowService].
func (s *followService) Follow(followerID uint, username string) error {
	target, err := s.findTarget(username)
	if err != nil {
		return err
	}

	if target.ID == followerID {
		return ErrCannotFollowSelf
	}

	return s.follows.Create(&models.Follow{
		FollowerID: followerID,
		FolloweeID: target.ID,
	})
}

// Unfollow implements [FollowService].
func (s *followService) Unfollow(followerID uint, username string) error {
	target, err := s.findTarget(username)
	if err != nil {
		return err
	}
	return s.follows.Delete(followerID, target.ID)
}

func (s *followService) findTarget(username string) (*models.User, error) {
	user, err := s.users.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
