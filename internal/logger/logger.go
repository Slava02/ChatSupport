package logger

import (
	"errors"
	"fmt"
	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	stdlog "log"
	"os"
	"syscall"

	"github.com/Slava02/ChatSupport/internal/buildinfo"
)

var LogLevel zap.AtomicLevel

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options -defaults-from=var
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
	clock          zapcore.Clock
	sentryDSN      string `validate:"omitempty,url"`
	env            string `validate:"omitempty,oneof=dev stage prod"`
}

var defaultOptions = Options{
	clock:          zapcore.DefaultClock,
	productionMode: false,
	env:            "prod",
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

	if err := LogLevel.UnmarshalText([]byte(opts.level)); err != nil {
		return fmt.Errorf("parse log level: %v", err)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "level",
		NameKey:        "component",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if opts.productionMode {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), LogLevel),
	}

	l := zap.New(zapcore.NewTee(cores...), zap.WithClock(opts.clock))

	if opts.sentryDSN != "" {
		sentryClient, errNewClient := NewSentryClient(opts.sentryDSN, opts.env, buildinfo.Version())
		if errNewClient != nil {
			return fmt.Errorf("couldn't init sentry client: %v", errNewClient)
		}

		core, errNewCore := zapsentry.NewCore(zapsentry.Configuration{
			Level:             zapcore.WarnLevel,
			EnableBreadcrumbs: true,
			BreadcrumbLevel:   zapcore.WarnLevel,
			Tags: map[string]string{
				"component": "system",
			},
		}, zapsentry.NewSentryClientFromClient(sentryClient))

		if errNewCore != nil {
			return fmt.Errorf("couldn't init sentry client: %v", errNewCore)
		}

		l = zapsentry.AttachCoreToLogger(core, l)
	}

	zap.ReplaceGlobals(l)

	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
