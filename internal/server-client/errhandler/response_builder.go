package errhandler

import (
	clientv1 "github.com/Slava02/ChatSupport/internal/server-client/v1"
	"github.com/Slava02/ChatSupport/pkg/pointer"
)

type Response struct {
	Error clientv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	return Response{
		clientv1.Error{
			Code:    clientv1.ErrorCode(code), //nolint:unconvert // unified types
			Message: msg,
			Details: pointer.PtrWithZeroAsNil(details),
		},
	}
}
