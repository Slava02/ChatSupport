package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// NewLogging returns a middleware that logs incoming requests with specific details.
func NewLogging(lg *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(eCtx echo.Context) bool {
			return eCtx.Request().Method == echo.OPTIONS
		},
		LogLatency:   true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogRequestID: true,
		LogUserAgent: true,
		LogStatus:    true,
		LogError:     true,
		LogValuesFunc: func(eCtx echo.Context, v middleware.RequestLoggerValues) error {
			lg := lg.With(
				zap.Duration("latency", v.Latency),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("path", v.URIPath),
				zap.String("request_id", v.RequestID),
				zap.String("user_agent", v.UserAgent),
			)

			uid, _ := userID(eCtx)
			lg = lg.With(zap.Stringer("user_id", uid))

			status := v.Status

			if err := v.Error; err != nil {
				lg = lg.With(zap.Error(err))
			}

			lg = lg.With(zap.Int("status", status))

			switch {
			case status >= 1000:
				lg.Error("business logic error")
			case status >= 500:
				lg.Error("server error")
			case status >= 400:
				lg.Error("client error")
			default:
				lg.Info("success")
			}

			return nil
		},
	})
}
