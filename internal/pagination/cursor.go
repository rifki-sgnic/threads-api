package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidCursor = errors.New("invalid cursor")

type Cursor struct {
	CreatedAt time.Time
	ID        uint
}

func (c Cursor) Encode() string {
	raw := fmt.Sprintf("%d_%d", c.CreatedAt.UnixMicro(), c.ID)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func Decode(s string) (*Cursor, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	parts := strings.Split(string(b), "_")
	if len(parts) != 2 {
		return nil, ErrInvalidCursor
	}

	micro, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	id, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	return &Cursor{CreatedAt: time.UnixMicro(micro), ID: uint(id)}, nil
}
