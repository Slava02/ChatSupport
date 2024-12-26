package messagesrepo

import (
	"github.com/Slava02/ChatSupport/internal/store"
	"github.com/Slava02/ChatSupport/internal/types"
	"time"
)

type Message struct {
	ID                  types.MessageID
	ChatID              types.ChatID
	AuthorID            types.UserID
	Body                string
	CreatedAt           time.Time
	IsVisibleForClient  bool
	IsVisibleForManager bool
	IsBlocked           bool
	IsService           bool
}

func adaptStoreMessage(m *store.Message) Message {
	return Message{
		ID:                  m.ID,
		ChatID:              m.ChatID,
		AuthorID:            m.AuthorID,
		Body:                m.Body,
		CreatedAt:           m.CreatedAt,
		IsVisibleForClient:  m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		IsBlocked:           m.IsBlocked,
		IsService:           m.IsService,
	}
}
