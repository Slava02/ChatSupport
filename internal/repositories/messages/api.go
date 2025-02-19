package messagesrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Slava02/ChatSupport/internal/store"
	"github.com/Slava02/ChatSupport/internal/store/message"
	"github.com/Slava02/ChatSupport/internal/types"
)

var ErrMsgNotFound = errors.New("message not found")

func (r *Repo) GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*Message, error) {
	msg, err := r.db.Message(ctx).Query().
		Where(message.InitialRequestID(reqID)).
		Only(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, fmt.Errorf("request id: %v: %w", reqID, ErrMsgNotFound)
		}
		return nil, fmt.Errorf("query message by request id: %v: %v", reqID, err)
	}

	res := adaptStoreMessage(msg)
	return &res, nil
}

// CreateClientVisible creates a message that is visible only to the client.
func (r *Repo) CreateClientVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	msg, err := r.db.Message(ctx).Create().
		SetInitialRequestID(reqID).
		SetProblemID(problemID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetBody(msgBody).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create message visible for client: %v", err)
	}

	res := adaptStoreMessage(msg)
	return &res, nil
}
