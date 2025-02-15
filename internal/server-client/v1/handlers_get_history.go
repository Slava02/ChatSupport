package clientv1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Slava02/ChatSupport/internal/types"
)

var stub = MessagesPage{Messages: []Message{
	{
		AuthorId:  types.NewUserID(),
		Body:      "Здравствуйте! Разберёмся.",
		CreatedAt: time.Now(),
		Id:        types.NewMessageID(),
	},
	{
		AuthorId:  types.MustParse[types.UserID]("7d67b14d-221e-4499-9be2-6707d7df1adc"),
		Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
		CreatedAt: time.Now().Add(-time.Minute),
		Id:        types.NewMessageID(),
	},
}}

func (h Handlers) PostGetHistory(eCtx echo.Context, _ PostGetHistoryParams) error {
	return eCtx.JSON(http.StatusOK, GetHistoryResponse{Data: stub})
}
