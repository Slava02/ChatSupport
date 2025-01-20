package gethistory

import (
	"github.com/Slava02/ChatSupport/internal/types"
	"github.com/Slava02/ChatSupport/internal/validator"
	"time"
)

type Request struct {
	ID       types.RequestID `validate:"required"`
	ClientID types.UserID    `validate:"required"`
	PageSize int             `validate:"omitempty,gte=10,lte=100,required_without=Cursor,excluded_with=Cursor"`
	Cursor   string          `validate:"omitempty,base64url,required_without=PageSize,excluded_with=PageSize"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	Messages   []Message `validate:"required"`
	NextCursor string    `validate:"required"`
}

type Message struct {
	ID         types.MessageID `validate:"required"`
	AuthorID   types.UserID    `validate:"required"`
	Body       string          `validate:"required"`
	CreatedAt  time.Time       `validate:"required"`
	IsReceived bool            `validate:"omitempty"`
	IsBlocked  bool            `validate:"omitempty"`
	IsService  bool            `validate:"omitempty"`
}
