package sendmessage

import (
	"time"

	"github.com/Slava02/ChatSupport/internal/types"
	"github.com/Slava02/ChatSupport/internal/validator"
)

type Request struct {
	ID          types.RequestID `validate:"required"`
	ClientID    types.UserID    `validate:"required"`
	MessageBody string          `validate:"required,min=1,max=3000"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	MessageID types.MessageID `validate:"required"`
	AuthorID  types.UserID    `validate:"required"`
	CreatedAt time.Time       `validate:"required"`
}
