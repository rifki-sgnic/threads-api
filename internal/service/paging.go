package service

import (
	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/pagination"
)

const (
	defaultPageLimit = 20
	maxPageLimit     = 50
)

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultPageLimit
	}

	if limit > maxPageLimit {
		return maxPageLimit
	}

	return limit
}

func decodeOptionalCursor(s string) (*pagination.Cursor, error) {
	if s == "" {
		return nil, nil
	}

	return pagination.Decode(s)
}

func postCursor(p models.Post) pagination.Cursor {
	return pagination.Cursor{
		CreatedAt: p.CreatedAt,
		ID:        p.ID,
	}
}
