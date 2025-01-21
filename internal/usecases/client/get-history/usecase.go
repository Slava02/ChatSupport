package gethistory

import (
	"context"
	"errors"
	"fmt"
	"github.com/Slava02/ChatSupport/internal/cursor"

	messagesrepo "github.com/Slava02/ChatSupport/internal/repositories/messages"
	"github.com/Slava02/ChatSupport/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=gethistorymocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type messagesRepository interface {
	GetClientChatMessages(
		ctx context.Context,
		clientID types.UserID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]messagesrepo.Message, *messagesrepo.Cursor, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("UseCase validation: %w", err)
	}

	return UseCase{opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("request validation: %w: %v", ErrInvalidRequest, err)
	}

	var reqCursor *messagesrepo.Cursor
	if req.Cursor != "" {
		if err := cursor.Decode(req.Cursor, &reqCursor); err != nil {
			return Response{}, fmt.Errorf("decode cursor: %w: %v", ErrInvalidCursor, err)
		}
	}

	messages, next, err := u.msgRepo.GetClientChatMessages(ctx, req.ClientID, req.PageSize, reqCursor)
	if err != nil {
		switch {
		case errors.Is(err, messagesrepo.ErrInvalidCursor):
			return Response{}, fmt.Errorf("get client chat messages: %w: %v", ErrInvalidCursor, err)
		default:
			return Response{}, fmt.Errorf("get client chat messages: %v", err)
		}
	}

	var nextCoursor string
	if next != nil {
		data, err := cursor.Encode(next)
		if err != nil {
			return Response{}, fmt.Errorf("encode cursor: %v", err)
		}
		nextCoursor = data
	}

	messagesArr := make([]Message, 0, len(messages))
	for _, m := range messages {
		messagesArr = append(messagesArr, Message{
			m.ID,
			m.AuthorID,
			m.Body,
			m.CreatedAt,
			m.IsVisibleForManager && !m.IsBlocked,
			m.IsBlocked,
			m.IsService,
		})
	}

	return Response{
		Messages:   messagesArr,
		NextCursor: nextCoursor,
	}, nil
}
