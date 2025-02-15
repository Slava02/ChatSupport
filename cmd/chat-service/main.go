package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/Slava02/ChatSupport/internal/clients/keycloak"
	"github.com/Slava02/ChatSupport/internal/config"
	"github.com/Slava02/ChatSupport/internal/logger"
	messagesrepo "github.com/Slava02/ChatSupport/internal/repositories/messages"
	clientv1 "github.com/Slava02/ChatSupport/internal/server-client/v1"
	serverdebug "github.com/Slava02/ChatSupport/internal/server-debug"
	"github.com/Slava02/ChatSupport/internal/store"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

const prod = true

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

	logger.MustInit(
		logger.NewOptions(cfg.Log.Level,
			logger.WithSentryDSN(cfg.Sentry.DSN),
			logger.WithEnv(cfg.Global.Env),
			logger.WithSentryDSN(cfg.Sentry.DSN),
		))
	defer logger.Sync()

	clientv1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %v", err)
	}

	kc, err := keycloakclient.New(keycloakclient.NewOptions(
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
	))
	if err != nil {
		return fmt.Errorf("failed to init keycloak client: %v", err)
	}

	storage, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.Stores.PSQL.Address,
		cfg.Stores.PSQL.Username,
		cfg.Stores.PSQL.Password,
		cfg.Stores.PSQL.Database,
		store.WithDebug(cfg.Stores.PSQL.Debug),
	))
	if err != nil {
		return fmt.Errorf("create store client: %v", err)
	}

	if err = storage.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
	}

	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(store.NewDatabase(storage)))
	if err != nil {
		return fmt.Errorf("failed initializing schema: %v", err)
	}

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	srvClient, err := initServerClient(
		prod,
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		clientv1Swagger,
		kc,
		cfg.Servers.Client.Access.Role,
		cfg.Servers.Client.Access.Resource,
		msgRepo,
	)
	if err != nil {
		return fmt.Errorf("failed to init client server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return srvClient.Run(ctx) })

	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
