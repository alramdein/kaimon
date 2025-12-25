package onrequest

import (
	"log"

	"github.com/alramdein/kaimon/pkg/framework"
	"github.com/alramdein/kaimon/pkg/middleware"
)

func init() {
	// Self-register this middleware
	middleware.RegisterOnRequestFactory("logger", func() middleware.OnRequestMiddleware {
		return NewLoggerMiddleware()
	})
}

// LoggerMiddleware logs incoming requests
type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (m *LoggerMiddleware) Name() string {
	return "logger"
}

func (m *LoggerMiddleware) Handle(next framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx framework.Context) error {
		req := ctx.Request()
		log.Printf("[%s] %s %s", req.Method, req.URL.Path, req.RemoteAddr)
		return next(ctx)
	}
}
