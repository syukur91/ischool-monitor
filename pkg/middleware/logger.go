package middleware

import (
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/labstack/echo"
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		Skipper Skipper
		Logger  *zap.Logger
		AppName string
	}
)

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			tenant := c.Param("tenant")
			if tenant == "" {
				tenant = "none"
			}

			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			latency := stop.Sub(start).Nanoseconds() / int64(time.Microsecond)
			latencyHuman := stop.Sub(start).String()
			bytesIn := req.Header.Get(echo.HeaderContentLength)
			if bytesIn == "" {
				bytesIn = "0"
			}

			config.Logger.Info("Request handled",
				zap.String("tenant", tenant),
				zap.String("type", "request"),
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", req.Host),
				zap.String("uri", req.RequestURI),
				zap.String("path", path),
				zap.String("method", req.Method),
				zap.String("referer", req.Referer()),
				zap.String("user_agent", req.UserAgent()),
				zap.Int("status", res.Status),
				zap.Int64("latency", latency),
				zap.String("latency_human", latencyHuman),
				zap.String("bytes_in", bytesIn),
				zap.String("bytes_out", strconv.FormatInt(res.Size, 10)),
			)

			// Todo: add body and header to log

			return nil
		}
	}
}
