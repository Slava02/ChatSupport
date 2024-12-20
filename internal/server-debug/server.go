package serverdebug

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/Slava02/ChatSupport/internal/buildinfo"
	"github.com/Slava02/ChatSupport/internal/logger"
	clientv1 "github.com/Slava02/ChatSupport/internal/server-client/v1"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr string `option:"mandatory" validate:"required,hostname_port"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	lg := zap.L().Named("server-debug")

	e := echo.New()
	e.Use(middleware.Recover())

	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
	index := newIndexPage()

	e.GET("/version", s.Version)
	index.addPage("/version", "Get build information")

	e.PUT("/log/level", s.ChangeLogLevel)
	e.GET("/log/level", s.GetLogLevel)

	{
		pprofMux := http.NewServeMux()
		pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
		pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		pprofMux.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
		pprofMux.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
		pprofMux.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
		pprofMux.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
		pprofMux.HandleFunc("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
		pprofMux.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

		e.GET("/debug/pprof/*", echo.WrapHandler(pprofMux))
		index.addPage("/debug/pprof/", "Go std profiler")
		index.addPage("/debug/pprof/profile?seconds=30", "Take half-min profile")
	}

	e.GET("/schema/client", s.GetSchema)
	index.addPage("/schema/client", "Get client OpenAPI specification")

	e.GET("/debug/error", s.SendError)
	index.addPage("/debug/error", "Debug sentry error event")

	e.GET("/", index.handler)
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx) //nolint:contextcheck // graceful shutdown with new context
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}

func (s *Server) Version(ctx echo.Context) error {
	info, err := json.Marshal(buildinfo.BuildInfo)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "couldn't marshal buildinfo")
	}

	_, err = ctx.Response().Write(info)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "couldn't write buildinfo to response")
	}

	return ctx.String(http.StatusOK, "completed")
}

func (s *Server) ChangeLogLevel(ctx echo.Context) error {
	level := ctx.FormValue("level")
	if level == "" {
		return ctx.String(http.StatusBadRequest, "level is required")
	}

	if err := logger.LogLevel.UnmarshalText([]byte(level)); err != nil {
		return ctx.String(http.StatusBadRequest, "parse log level")
	}

	logger.LogLevel.SetLevel(logger.LogLevel.Level())

	return ctx.String(http.StatusOK, "log level updated")
}

func (s *Server) GetLogLevel(ctx echo.Context) error {
	level := logger.LogLevel.String()
	return ctx.JSONPretty(http.StatusOK, map[string]string{"level": level}, " ")
}

func (s *Server) SendError(ctx echo.Context) error {
	s.lg.Error("look for me in sentry")
	return ctx.String(http.StatusOK, "event sent")
}

func (s *Server) GetSchema(ctx echo.Context) error {
	swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("couldn't get swagger")
	}

	return ctx.JSONPretty(http.StatusOK, swagger, " ")
}
