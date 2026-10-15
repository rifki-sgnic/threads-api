package service

import (
	"github.com/rifki-sgnic/threads-api/internal/pagination"
	"github.com/rifki-sgnic/threads-api/internal/repository"
)

const (
	defaultFeedLimit = 20
	maxFeedLimit     = 50
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
	if limit <= 0 {
		limit = defaultFeedLimit
	}
	if limit > maxFeedLimit {
		limit = maxFeedLimit
	}

	var cursor *pagination.Cursor
	if cursorStr != "" {
		c, err := pagination.Decode(cursorStr)
		if err != nil {
			return nil, err
		}
		cursor = c
	}

	// ambil limit+1: kalau dapet lebih dari limit, berarti masih ada halaman berikutnya
	posts, err := s.posts.Feed(userID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	page := &FeedPage{}
	if len(posts) > limit {
		posts = posts[:limit]
		last := posts[len(posts)-1]
		page.NextCursor = pagination.Cursor{CreatedAt: last.CreatedAt, ID: last.ID}.Encode()
	}

	page.Items, err = s.stats.Attach(userID, posts)
	if err != nil {
		return nil, err
	}
	return page, nil
}
