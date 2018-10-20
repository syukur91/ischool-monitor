package middleware

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/syukur91/ischool-monitor/pkg/apierror"
	"github.com/syukur91/ischool-monitor/pkg/response"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ErrorHandler ...
func ErrorHandler(logger *zap.Logger) func(error, echo.Context) {
	return func(err error, c echo.Context) {

		tenant := c.Param("tenant")
		if tenant == "" {
			tenant = "none"
		}

		ae, ok := err.(*apierror.APIError)
		if !ok {
			err = apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), errors.Wrap(err, "errorhandler: internal error"))
			ae, _ = err.(*apierror.APIError)
		}

		/*
			// TODO: get only top 2 frames of error stack
			errStack, ok := ae.Err.(stackTracer)
			if !ok {
				panic("oops, err does not implement stackTracer")
			}

			st := errStack.StackTrace()
			fmt.Printf("%+v", st[0:2]) // top two frames
		*/

		req := c.Request()
		path := req.URL.Path
		if path == "" {
			path = "/"
		}

		logger.Error(ae.Message,
			zap.String("path", path),
			zap.String("tenant", tenant),
			zap.String("method", req.Method),
			zap.Error(ae.Err))

		r := new(response.Response)
		es := make([]response.Error, 1)
		es[0] = response.Error{Status: ae.HTTPStatus, Code: ae.Code, Title: http.StatusText(ae.Code), Detail: ae.Message}
		r.Errors = es

		// Send response
		if !c.Response().Committed {
			if c.Request().Method == "HEAD" { // Issue #608
				err = c.NoContent(ae.HTTPStatus)
			} else {
				err = c.JSON(ae.HTTPStatus, r)
			}
			if err != nil {
				logger.Error(err.Error(),
					zap.String("path", path),
					zap.String("tenant", tenant),
					zap.String("method", req.Method),
					zap.Error(err))
			}
		}
	}
}
