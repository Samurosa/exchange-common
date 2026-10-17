package cursor

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

func EncodeCursor(reqCursor Cursor) (string, error) {

	raw := fmt.Sprintf("%d:%s", reqCursor.CreatedAt.UnixNano(), reqCursor.ID)

	return base64.RawURLEncoding.EncodeToString([]byte(raw)), nil
}

func DecodeCursor(reqCursor string) (Cursor, error) {

	data, err := base64.RawURLEncoding.DecodeString(reqCursor)
	if err != nil {
		return Cursor{}, err
	}

	splitCursor := strings.Split(string(data), ":")

	if len(splitCursor) != 2 {
		return Cursor{}, ErrInvalidCursor
	}

	timestamp, err := strconv.ParseInt(splitCursor[0], 10, 64)
	if err != nil {
		return Cursor{}, ErrInvalidCursor
	}

	createdAt := time.Unix(0, timestamp)

	if createdAt.IsZero() {
		return Cursor{}, ErrInvalidCursor
	}

	id := splitCursor[1]

	if _, err := uuid.Parse(id); err != nil {
		return Cursor{}, ErrInvalidCursor
	}

	return Cursor{
		CreatedAt: createdAt,
		ID:        id,
	}, nil
}
