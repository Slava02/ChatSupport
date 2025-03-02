package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"github.com/tchap/zapext/v2/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Slava02/ChatSupport/internal/buildinfo"
)

const sentryLvl = zapcore.WarnLevel

var Level zap.AtomicLevel

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options -defaults-from=var
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
	clock          zapcore.Clock
	sentryDSN      string `validate:"omitempty,url"`
	sentryEnv      string `validate:"omitempty,oneof=dev stage prod"`
}

var defaultOptions = Options{
	clock: zapcore.DefaultClock,
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate options: %v", err)
	}

	var err error
	Level, err = zap.ParseAtomicLevel(opts.level)
	if err != nil {
		return fmt.Errorf("parse level: %v", err)
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.NameKey = "component"
	cfg.TimeKey = "T"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if opts.productionMode {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(cfg)
	} else {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(cfg)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, Level),
	}
	if dsn := opts.sentryDSN; dsn != "" {
		sentryClient, err := NewSentryClient(dsn, opts.sentryEnv, buildinfo.Version())
		if err != nil {
			return fmt.Errorf("new sentry client: %v", err)
		}
		cores = append(cores, zapsentry.NewCore(sentryLvl, sentryClient))
	}

	l := zap.New(zapcore.NewTee(cores...), zap.WithClock(opts.clock))
	zap.ReplaceGlobals(l)

	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
