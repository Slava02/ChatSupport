package clientv1

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	internalerrors "github.com/Slava02/ChatSupport/internal/errors"
	"github.com/Slava02/ChatSupport/internal/middlewares"
	gethistory "github.com/Slava02/ChatSupport/internal/usecases/client/get-history"
	"github.com/Slava02/ChatSupport/pkg/pointer"
)

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	clientID := middlewares.MustUserID(eCtx)

	var req GetHistoryRequest
	if err := eCtx.Bind(&req); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}

	resp, err := h.getHistory.Handle(eCtx.Request().Context(), gethistory.Request{
		ClientID: clientID,
		ID:       params.XRequestID,
		Cursor:   pointer.Indirect(req.Cursor),
		PageSize: pointer.Indirect(req.PageSize),
	})
	if err != nil {
		switch {
		case errors.Is(err, gethistory.ErrInvalidRequest):
			return internalerrors.NewServerError(400, "getHistory InvalidRequest", err)
		case errors.Is(err, gethistory.ErrInvalidCursor):
			return internalerrors.NewServerError(400, "getHistory InvalidCursor", err)
		default:
			return fmt.Errorf("handle `get history`: %v", err)
		}
	}

	messages := make([]Message, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		message := Message{
			m.AuthorID.AsPointer(),
			m.Body,
			m.CreatedAt,
			m.ID,
			m.IsBlocked,
			m.IsReceived,
			m.IsService,
		}

		messages = append(messages, message)
	}

	return eCtx.JSON(200, GetHistoryResponse{
		Data: &MessagesPage{
			Messages: messages,
			Next:     resp.NextCursor,
		},
		Error: nil,
	})
}
