package service

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/repository"
)

type PostView struct {
	Post       models.Post
	LikeCount  int64
	ReplyCount int64
	LikedByMe  bool
}

type PostStatsLoader struct {
	posts repository.PostRepository
	likes repository.LikeRepository
}

func NewPostStatsLoader(posts repository.PostRepository, likes repository.LikeRepository) *PostStatsLoader {
	return &PostStatsLoader{posts: posts, likes: likes}
}

func (l *PostStatsLoader) Attach(viewerID uint, posts []models.Post) ([]PostView, error) {
	ids := make([]uint, len(posts))
	for i, p := range posts {
		ids[i] = p.ID
	}

	likeCounts, err := l.likes.CountByPosts(ids)
	if err != nil {
		return nil, err
	}
	replyCounts, err := l.posts.CountReplies(ids)
	if err != nil {
		return nil, err
	}

	liked := map[uint]bool{}
	if viewerID != 0 {
		liked, err = l.likes.LikedByUser(viewerID, ids)
		if err != nil {
			return nil, err
		}
	}

	views := make([]PostView, len(posts))
	for i, p := range posts {
		views[i] = PostView{
			Post:       p,
			LikeCount:  likeCounts[p.ID],
			ReplyCount: replyCounts[p.ID],
			LikedByMe:  liked[p.ID],
		}
	}
	return views, nil
}
