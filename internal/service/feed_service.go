package service

import (
	"github.com/rifki-sgnic/threads-api/internal/pagination"
	"github.com/rifki-sgnic/threads-api/internal/repository"
)

type FeedPage struct {
	Items      []PostView
	NextCursor string
}

type FeedService interface {
	GetFeed(userID uint, cursor string, limit int) (*FeedPage, error)
}

type feedService struct {
	posts repository.PostRepository
	stats *PostStatsLoader
}

func NewFeedService(posts repository.PostRepository, stats *PostStatsLoader) FeedService {
	return &feedService{posts: posts, stats: stats}
}

func (s *feedService) GetFeed(userID uint, cursorStr string, limit int) (*FeedPage, error) {
	limit = normalizeLimit(limit)

	cursor, err := decodeOptionalCursor(cursorStr)
	if err != nil {
		return nil, err
	}

	// ambil limit+1: kalau dapet lebih dari limit, berarti masih ada halaman berikutnya
	posts, err := s.posts.Feed(userID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	posts, next := pagination.Trim(posts, limit, postCursor)

	items, err := s.stats.Attach(userID, posts)
	if err != nil {
		return nil, err
	}

	return &FeedPage{Items: items, NextCursor: next}, nil
}
