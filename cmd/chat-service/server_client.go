package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	keycloakclient "github.com/Slava02/ChatSupport/internal/clients/keycloak"
	chatrepo "github.com/Slava02/ChatSupport/internal/repositories/chats"
	messagesrepo "github.com/Slava02/ChatSupport/internal/repositories/messages"
	problemsrepo "github.com/Slava02/ChatSupport/internal/repositories/problems"
	serverclient "github.com/Slava02/ChatSupport/internal/server-client"
	clientv1 "github.com/Slava02/ChatSupport/internal/server-client/v1"
	"github.com/Slava02/ChatSupport/internal/store"
	gethistoryUC "github.com/Slava02/ChatSupport/internal/usecases/client/get-history"
	sendmessageUC "github.com/Slava02/ChatSupport/internal/usecases/client/send-message"
)

const nameServerClient = "server-client"

func initServerClient(
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	chatRepo chatrepo.Repo,
	messageRepo messagesrepo.Repo,
	problemsRepo problemsrepo.Repo,
	keycloak *keycloakclient.Client,
	requiredResource string,
	requiredRole string,
	errHandler echo.HTTPErrorHandler,
	db *store.Database,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistoryUC.New(gethistoryUC.NewOptions(&messageRepo))
	if err != nil {
		return nil, fmt.Errorf("create get-history use case: %v", err)
	}

	sendMessageUseCase, err := sendmessageUC.New(sendmessageUC.NewOptions(&chatRepo, &messageRepo, &problemsRepo, db))
	if err != nil {
		return nil, fmt.Errorf("create send-message use case: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(getHistoryUseCase, sendMessageUseCase))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	srv, err := serverclient.New(serverclient.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		v1Handlers,
		keycloak,
		requiredResource,
		requiredRole,
		errHandler,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
