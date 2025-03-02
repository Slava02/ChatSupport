package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/Slava02/ChatSupport/internal/clients/keycloak"
	"github.com/Slava02/ChatSupport/internal/config"
	"github.com/Slava02/ChatSupport/internal/logger"
	chatrepo "github.com/Slava02/ChatSupport/internal/repositories/chats"
	messagesrepo "github.com/Slava02/ChatSupport/internal/repositories/messages"
	problemrepo "github.com/Slava02/ChatSupport/internal/repositories/problems"
	"github.com/Slava02/ChatSupport/internal/server-client/errhandler"
	clientv1 "github.com/Slava02/ChatSupport/internal/server-client/v1"
	serverdebug "github.com/Slava02/ChatSupport/internal/server-debug"
	"github.com/Slava02/ChatSupport/internal/store"
)

const nameMain = "main"

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	if err := logger.Init(logger.NewOptions(
		cfg.Log.Level,
		logger.WithProductionMode(cfg.Global.IsProduction()),
		logger.WithSentryDSN(cfg.Sentry.DSN),
		logger.WithSentryEnv(cfg.Global.Env),
	)); err != nil {
		return fmt.Errorf("init logger: %v", err)
	}

	defer logger.Sync()

	lg := zap.L().Named(nameMain)

	// Storage.
	if cfg.Global.IsProduction() && cfg.Stores.PSQL.Debug {
		lg.Warn("psql client in the debug mode")
	}

	storage, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.Stores.PSQL.Addr,
		cfg.Stores.PSQL.Username,
		cfg.Stores.PSQL.Password,
		cfg.Stores.PSQL.Database,
		store.WithDebug(cfg.Stores.PSQL.Debug),
	))
	if err != nil {
		return fmt.Errorf("create store client: %v", err)
	}
	defer multierr.AppendInvoke(&errReturned, multierr.Close(storage))

	// Migrations.
	if err = storage.Schema.Create(ctx); err != nil {
		return fmt.Errorf("migrate: %v", err)
	}

	// Repositories.
	db := store.NewDatabase(storage)

	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("messages repo: %v", err)
	}

	chatRepo, err := chatrepo.New(chatrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("chat repo: %v", err)
	}

	problemRepo, err := problemrepo.New(problemrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("problem repo: %v", err)
	}

	// Clients.
	kc, err := keycloakclient.New(keycloakclient.NewOptions(
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
		keycloakclient.WithDebugMode(cfg.Clients.Keycloak.DebugMode),
	))
	if err != nil {
		return fmt.Errorf("create keycloak client: %v", err)
	}
	if cfg.Global.IsProduction() && cfg.Clients.Keycloak.DebugMode {
		zap.L().Warn("keycloak client in the debug mode")
	}

	// Servers.
	clientV1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get client v1 swagger: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(
		lg,
		cfg.Global.IsProduction(),
		errhandler.ResponseBuilder,
	))
	if err != nil {
		return fmt.Errorf("init error handler: %v", err)
	}

	srvClient, err := initServerClient(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		clientV1Swagger,
		*chatRepo,
		*msgRepo,
		*problemRepo,
		kc,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
		errHandler.Handle,
		db,
	)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
	}

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(
		cfg.Servers.Debug.Addr,
		clientV1Swagger,
	))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvClient.Run(ctx) })
	eg.Go(func() error { return srvDebug.Run(ctx) })

	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
