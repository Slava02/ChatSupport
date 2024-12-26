package messagesrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Slava02/ChatSupport/internal/store"
	"github.com/Slava02/ChatSupport/internal/store/chat"
	"github.com/Slava02/ChatSupport/internal/store/message"
	"time"

	"github.com/Slava02/ChatSupport/internal/types"
)

const (
	maxPageSize = 100
	minPageSize = 10
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
)

type Cursor struct {
	LastCreatedAt time.Time
	PageSize      int
}

func (c *Cursor) Validate() error {
	if c.LastCreatedAt.IsZero() {
		return errors.New("LastCreatedAt field must be specified")
	}

	return validPageSize(c.PageSize)
}

func validPageSize(pageSize int) error {
	if pageSize < minPageSize || pageSize > maxPageSize {
		return fmt.Errorf("PageSize field must be in [%d, %d]", maxPageSize, maxPageSize)
	}

	return nil
}

// GetClientChatMessages returns Nth page of messages in the chat for client side.
func (r *Repo) GetClientChatMessages(
	ctx context.Context,
	clientID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	lastMessage := time.Now().AddDate(100, 0, 0)

	if cursor != nil {
		if err := cursor.Validate(); err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}
		pageSize, lastMessage = cursor.PageSize, cursor.LastCreatedAt
	} else if pageSize != 0 {
		if err := validPageSize(pageSize); err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrInvalidPageSize, err)
		}
	}

	messages, err := r.db.Message(ctx).Query().
		Unique(false).
		Where(
			message.CreatedAtLT(lastMessage),
			message.IsVisibleForClient(true),
			message.HasChatWith(chat.ClientID(clientID)),
		).
		Limit(pageSize + 1).
		Order(store.Desc(message.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("select messages :%v", err)
	}

	res := make([]Message, 0, len(messages))
	for _, m := range messages {
		res = append(res, adaptStoreMessage(m))
	}

	if len(messages) <= pageSize {
		return res, nil, nil
	}

	res = res[:len(res)-1]

	return res, &Cursor{
		LastCreatedAt: res[len(res)-1].CreatedAt,
		PageSize:      pageSize,
	}, nil
}
