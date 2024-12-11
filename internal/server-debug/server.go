package serverdebug

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/Slava02/ChatSupport/internal/buildinfo"
	"github.com/Slava02/ChatSupport/internal/logger"
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

	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))
	index.addPage("/debug/pprof", "Go std profiler")
	index.addPage("/debug/pprof/profile?seconds=30", "Take half-min profile")

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
	return ctx.JSON(http.StatusOK, map[string]string{"level": level})
}

func (s *Server) SendError(ctx echo.Context) error {
	s.lg.Error("look for me in sentry")
	return ctx.String(http.StatusOK, "event sent")
}
