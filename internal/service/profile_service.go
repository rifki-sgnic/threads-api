package service

import (
	"errors"
	"strings"

	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/pagination"
	"github.com/rifki-sgnic/threads-api/internal/repository"
	"gorm.io/gorm"
)

type Profile struct {
	User           models.User
	FollowerCount  int64
	FollowingCount int64
	PostCount      int64
	IsFollowing    bool
}

type ProfileService interface {
	GetProfile(viewerID uint, username string) (*Profile, error)
	GetUserPosts(viewerID uint, username, cursor string, limit int) (*FeedPage, error)
	UpdateBio(userID uint, bio string) (*models.User, error)
}

type profileService struct {
	users   repository.UserRepository
	follows repository.FollowRepository
	posts   repository.PostRepository
	stats   *PostStatsLoader
}

func NewProfileService(
	users repository.UserRepository,
	follows repository.FollowRepository,
	posts repository.PostRepository,
	stats *PostStatsLoader,
) ProfileService {
	return &profileService{users: users, follows: follows, posts: posts, stats: stats}
}

func findUserByUsername(users repository.UserRepository, username string) (*models.User, error) {
	user, err := users.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound // dari follow_service.go
		}
		return nil, err
	}
	return user, nil
}

func (s *profileService) GetProfile(viewerID uint, username string) (*Profile, error) {
	user, err := findUserByUsername(s.users, username)
	if err != nil {
		return nil, err
	}

	followers, err := s.follows.CountFollowers(user.ID)
	if err != nil {
		return nil, err
	}
	following, err := s.follows.CountFollowing(user.ID)
	if err != nil {
		return nil, err
	}
	postCount, err := s.posts.CountByUser(user.ID)
	if err != nil {
		return nil, err
	}

	isFollowing := false
	if viewerID != 0 && viewerID != user.ID {
		isFollowing, err = s.follows.IsFollowing(viewerID, user.ID)
		if err != nil {
			return nil, err
		}
	}

	return &Profile{
		User:           *user,
		FollowerCount:  followers,
		FollowingCount: following,
		PostCount:      postCount,
		IsFollowing:    isFollowing,
	}, nil
}

func (s *profileService) GetUserPosts(viewerID uint, username, cursorStr string, limit int) (*FeedPage, error) {
	user, err := findUserByUsername(s.users, username)
	if err != nil {
		return nil, err
	}

	limit = normalizeLimit(limit)

	cursor, err := decodeOptionalCursor(cursorStr)
	if err != nil {
		return nil, err
	}

	posts, err := s.posts.FindByUser(user.ID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	posts, next := pagination.Trim(posts, limit, postCursor)

	items, err := s.stats.Attach(viewerID, posts)
	if err != nil {
		return nil, err
	}
	return &FeedPage{Items: items, NextCursor: next}, nil
}

func (s *profileService) UpdateBio(userID uint, bio string) (*models.User, error) {
	if err := s.users.UpdateBio(userID, strings.TrimSpace(bio)); err != nil {
		return nil, err
	}

	user, err := s.users.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
