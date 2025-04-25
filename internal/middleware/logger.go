package middleware

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	zlog "github.com/rs/zerolog/log"
	"time"
)

func ZerologRequestLogger(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	stop := time.Now()
	latency := stop.Sub(start).Milliseconds()

	statusCode := c.Response().StatusCode()
	if err != nil {
		var e *fiber.Error
		if errors.As(err, &e) {
			statusCode = e.Code
		} else if statusCode == 0 || statusCode == fiber.StatusOK {
			statusCode = fiber.StatusInternalServerError
		}
	}

	logEvent := zlog.Info()

	if statusCode >= 500 || statusCode >= 400 {
		logEvent = zlog.Error().Err(err)

		handlerErr := c.Locals("err")
		if actualErr, ok := handlerErr.(error); ok {
			logEvent = zlog.Error().Err(actualErr)
		}
	}

	logEvent.
		Str("remote_ip", c.IP()).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Int("status", statusCode).
		Str("latency", fmt.Sprintf("%dms", latency)).
		Int("response_size", len(c.Response().Body())).
		Str("user_agent", c.Get(fiber.HeaderUserAgent)).
		Msg("API telah ter hit")

	return err
}
