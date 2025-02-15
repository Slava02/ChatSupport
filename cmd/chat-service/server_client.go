package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/Slava02/ChatSupport/internal/clients/keycloak"
	messagesrepo "github.com/Slava02/ChatSupport/internal/repositories/messages"
	serverclient "github.com/Slava02/ChatSupport/internal/server-client"
	"github.com/Slava02/ChatSupport/internal/server-client/errhandler"
	clientv1 "github.com/Slava02/ChatSupport/internal/server-client/v1"
	gethistory "github.com/Slava02/ChatSupport/internal/usecases/client/get-history"
)

const nameServerClient = "server-client"

func initServerClient(
	productionMode bool,
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	keycloak *keycloakclient.Client,
	requiredRole string,
	requiredResource string,
	msgRepo *messagesrepo.Repo,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(msgRepo))
	if err != nil {
		return nil, fmt.Errorf("create get history use case: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(getHistoryUseCase))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	httpErrorHandler, err := errhandler.New(errhandler.NewOptions(lg, productionMode, errhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create http error handler: %v", err)
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
		httpErrorHandler.Handle,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
