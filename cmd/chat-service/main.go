package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	keycloakclient "github.com/Slava02/ChatSupport/internal/clients/keycloak"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/Slava02/ChatSupport/internal/config"
	"github.com/Slava02/ChatSupport/internal/logger"
	clientv1 "github.com/Slava02/ChatSupport/internal/server-client/v1"
	serverdebug "github.com/Slava02/ChatSupport/internal/server-debug"
)

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

	logger.MustInit(
		logger.NewOptions(cfg.Log.Level,
			logger.WithSentryDSN(cfg.Sentry.DSN),
			logger.WithEnv(cfg.Global.Env),
			logger.WithSentryDSN(cfg.Sentry.DSN),
		))
	defer logger.Sync()

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	clientv1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %v", err)
	}

	// TODO
	keyCloakClient, err := keycloakclient.New(keycloakclient.NewOptions())

	srvClient, err := initServerClient(cfg.Servers.Client.Addr, cfg.Servers.Client.AllowOrigins, clientv1Swagger)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
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
